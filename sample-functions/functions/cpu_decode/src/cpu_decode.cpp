/*****************************************************************
 * Copyright 2024 NTT Corporation, FUJITSU LIMITED
 *****************************************************************/
#include <iostream>
#include <thread>
#include <atomic>
#include <iomanip>
#include <string>
#include <cerrno>
#include <cstring>
#include <unistd.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <opencv2/opencv.hpp>
#include "cpu_decode.hpp"
#include "tp.hpp"
#include "connect_thread_queue.hpp"
#include "liblogging.h"
#include "recv_header_log.hpp"
#include "libshmem.h"
#include "libfpgactl.h"
#include "libchain.h"
#include "libdmacommon.h"

//#define ALLOCATE_SRC_SHMEM_MODE
//#define EXEC_FPGA_ENQUEUE_MODE
//#define EXEC_FPGA_DEQUEUE_MODE
//#define DEBUG_PRINT
//#define DEBUG_FPGA_OUT_IMAGE_TO_MP4

// Initial parameter value
const char *FPGA_DEV_NAME = "/dev/xpcie0";
uint16_t FPGA_CH_ID = 0;
double FRAME_FPS = 5.0;
uint32_t FRAME_WIDTH = 3840;
uint32_t FRAME_HEIGHT = 2160;
uint32_t FRAME_CHANNEL = 3;
uint32_t FPGA_OUT_FRAME_WIDTH = 1280;
uint32_t FPGA_OUT_FRAME_HEIGHT = 1280;
void *FPGA_SRC_SHMEM_ADDR = nullptr;
const char *FPGA_DMA_HOST_TO_DEV_CONNECTOR_ID = "enq_connector_id";
const char *DPDK_FILE_PREFIX = "0";
uint32_t VIDEO_CONNECT_LIMIT = 0;
videosrc_protocol VIDEOSRC_PROTOCOL = videosrc_protocol::RTP;
protocol OUTDST_PROTOCOL = protocol::DMA;
std::string OUTDST_IPA("127.0.0.1");
uint16_t OUTDST_PORT = 0;

// TCP connect retry limit
static constexpr uint32_t TCP_CONNECT_RETRY_LIMIT = 60;

// decode_thread termination notification
std::atomic<bool> is_decode_thread_fin(false);
// recv_fpga_deq_thread termination notification
std::atomic<bool> is_recv_fpga_deq_thread_fin(false);

// queue info
dma_info_t enqdmainfo_channel;
dma_info_t enqdmainfo;
dma_info_t deqdmainfo_channel;
dma_info_t deqdmainfo;

// Chain Control Settings
#ifdef MODULE_FPGA
const chain_ctrl_cid_t g_param_chain_ctrl_tbls[CH_NUM_MAX] = {
	/* { FUNCTION CHID, LLDMA CID } */
	{ 0,  0 },
	{ 1,  1 },
	{ 2,  2 },
	{ 3,  3 },
	{ 4,  4 },
	{ 5,  5 },
	{ 6,  6 },
	{ 7,  7 },
	{ 0,  8 },
	{ 1,  9 },
	{ 2, 10 },
	{ 3, 11 },
	{ 4, 12 },
	{ 5, 13 },
	{ 6, 14 },
	{ 7, 15 },
};
#else
const chain_ctrl_cid_t g_param_chain_ctrl_tbls[CH_NUM_MAX] = {
	/* { FUNCTION CHID, LLDMA CID } */
	{ 0, 512 },
	{ 1, 513 },
	{ 2, 514 },
	{ 3, 515 },
	{ 4, 516 },
	{ 5, 517 },
	{ 6, 518 },
	{ 7, 519 },
	{ 0, 520 },
	{ 1, 521 },
	{ 2, 522 },
	{ 3, 523 },
	{ 4, 524 },
	{ 5, 525 },
	{ 6, 526 },
	{ 7, 527 },
};
#endif

// Get Date String
std::string datetime2str() {
	time_t t = time(nullptr);
	const tm* lt = localtime(&t);
	std::stringstream s;
	s << lt->tm_year + 1900;
	s << std::setw(2) << std::setfill('0') << lt->tm_mon + 1;
	s << std::setw(2) << std::setfill('0') << lt->tm_mday;
	s << std::setw(2) << std::setfill('0') << lt->tm_hour;
	s << std::setw(2) << std::setfill('0') << lt->tm_min;
	s << std::setw(2) << std::setfill('0') << lt->tm_sec;

	return s.str();
}

// Get string from environment variable
int32_t env2str(const char *env, std::string &str, bool isRequired) {
	const char* e = std::getenv(env);
	if (e == nullptr) {
		if (isRequired) {
			std::string err = "getenv(\"" + std::string(env) + "\") failed!";
			throw std::logic_error(err);
		} else {
			return -1;
		}
	}
	str = std::string(e);

	return 0;
}

// Get numeric value from environment variable
template <typename T>
int32_t env2num(const char *env, T &num, bool isRequired) {
	std::string s;
	int32_t ret = env2str(env, s, isRequired);
	if (ret < 0) {
		return -1;
	}
	num = stold(s);

	return 0;
}

template <typename T>
bool recvthreadque_polling(ConnectThreadQueue<T> &cq)
{
	bool is_fin = false;
	while (true) {
		// Checking whether there is an input data queue
		if (cq.size() == 0) {
			// Checking the Status of recv_fpga_deq_thread
			if (is_recv_fpga_deq_thread_fin.load(std::memory_order_acquire)) {
				// Flag exit if no input data queue and recv_fpga_deq_thread completes
				is_fin = true;
				break;
			}
		} else {
			break;
		}
		std::this_thread::sleep_for(std::chrono::microseconds(5000));
	}
	return is_fin;
}

template <typename T>
bool decodethreadque_polling(ConnectThreadQueue<T> &cq)
{
	bool is_fin = false;
	while (true) {
		// Checking whether there is an input data queue
		if (cq.size() == 0) {
			// Check the status of decode_thread
			if (is_decode_thread_fin.load(std::memory_order_acquire)) {
				// End flag if no input data queue and decode thread is complete
				is_fin = true;
				break;
			}
		} else {
			break;
		}
		std::this_thread::sleep_for(std::chrono::microseconds(5000));
	}
	return is_fin;
}

int32_t tcp_client_open(int32_t &sockfd)
{
	struct sockaddr_in addr;
	memset(&addr, 0, sizeof(struct sockaddr_in));

	// Setting the destination IP address and port number
	addr.sin_family = AF_INET;
	addr.sin_port = htons(OUTDST_PORT);
	addr.sin_addr.s_addr = inet_addr(OUTDST_IPA.c_str());

	uint32_t connect_cnt = 0;
	bool is_socket_open = false;
	while(true) {
		if (!is_socket_open) {
			// socket creation
			sockfd = socket(AF_INET, SOCK_STREAM, 0);
			if (sockfd < 0) {
				Logging::set(LogLevel::ERROR, "TCP socket error. errno(" + std::to_string(errno) + ": " + strerror(errno) + ")");
				return -1;
			}
			is_socket_open = true;
		}

		// Request to connect to server
		connect_cnt++;
		int32_t ret = connect(sockfd, (struct sockaddr*)&addr, sizeof(struct sockaddr_in));
		if (ret < 0) {
			Logging::set(LogLevel::INFO, "TCP connect try " + std::to_string(connect_cnt) + ". (" + std::string(strerror(errno)) + ")");
			shutdown(sockfd, SHUT_RDWR);
			close(sockfd);
			is_socket_open = false;
			if (connect_cnt > TCP_CONNECT_RETRY_LIMIT - 1) {
				Logging::set(LogLevel::ERROR, "TCP connect error. connect retry count reached limit (" + std::to_string(TCP_CONNECT_RETRY_LIMIT) + ")");
				return -1;
			} else {
				std::this_thread::sleep_for(std::chrono::seconds(1));
				// connection request retry
				continue;
			}
		}
		Logging::set(LogLevel::INFO, "TCP connect try " + std::to_string(connect_cnt) + ". (Success)");
		break;
	}

	return 0;
}

int32_t tcp_client_close(int32_t &sockfd)
{
	int32_t ret = shutdown(sockfd, SHUT_RDWR);
	if (ret < 0) {
		Logging::set(LogLevel::ERROR, "TCP socket shutdown error. errno(" + std::to_string(errno) + ": " + strerror(errno) + ")");
		return -1;
	}

	ret = close(sockfd);
	if (ret < 0) {
		Logging::set(LogLevel::ERROR, "TCP socket close error. errno(" + std::to_string(errno) + ": " + strerror(errno) + ")");
		return -1;
	}
	Logging::set(LogLevel::INFO, "TCP socket close Success");

	return 0;
}

void recv_data_thread(ConnectThreadQueue<recv_fpga_data_t> &cq)
{
	Logging::set(LogLevel::INFO, "--- recv_data_thread start ---");

	try {
#if defined(DEBUG_FPGA_OUT_IMAGE_TO_MP4)
		double out_fps = FRAME_FPS;
		int32_t fourcc = cv::VideoWriter::fourcc('a', 'v', 'c', '1');
		cv::VideoWriter writer;
		std::string outfile = "fpga_out_frame_" + datetime2str() + ".mp4";
		std::string videosink = "appsrc ! videoconvert ! queue ! openh264enc complexity=low rate-control=off ! video/x-h264,stream-format=byte-stream,alignment=au,profile=high ! queue ! h264parse ! qtmux ! filesink location=./" + outfile;
		writer.open(videosink, cv::CAP_GSTREAMER, fourcc, out_fps, cv::Size(FPGA_OUT_FRAME_WIDTH, FPGA_OUT_FRAME_HEIGHT));
		if (writer.isOpened()) {
			Logging::set(LogLevel::INFO, "VideoWriter opened.");
		} else {
			Logging::set(LogLevel::ERROR, "VideoWriter open failed.");
		}
#endif //DEBUG_FPGA_OUT_IMAGE_TO_MP4

		size_t frame_cnt = 0;

		while (true) {
			if (recvthreadque_polling<recv_fpga_data_t>(cq)) {
				// Completion if there is no input data queue and recv_fpga_deq_thread completes
#if defined(DEBUG_FPGA_OUT_IMAGE_TO_MP4)
				writer.release();
#endif //DEBUG_FPGA_OUT_IMAGE_TO_MP4
				break;
			}

			// FPGA output data acquisition
			recv_fpga_data_t recvdata = cq.pop();
			frame_cnt++;
#if defined(DEBUG_PRINT)
			std::cout << "recv data: " << frame_cnt << std::endl;
#endif //DEBUG_PRINT

			//-----------------------------------------------------
			// Log output of received frame header
			//-----------------------------------------------------
			RecvHeaderLog::set(recvdata.frame_header, frame_cnt);

			//-----------------------------------------------------
			// Output of received image data
			//-----------------------------------------------------
#if defined(DEBUG_FPGA_OUT_IMAGE_TO_MP4)
			writer << recvdata.mat;
#endif
		}

	} catch (const std::exception& error) {
		Logging::set(LogLevel::ERROR, error.what());
	} catch (...) {
		Logging::set(LogLevel::ERROR, "Unknown/internal exception.");
	}

	Logging::set(LogLevel::INFO, "--- recv_data_thread finish ---");
#if defined(DEBUG_PRINT)
	std::cout << "recv_data_thread finish" << std::endl;
#endif //DEBUG_PRINT
}

void recv_fpga_deq_thread(ConnectThreadQueue<dmacmd_info_t> &cq, ConnectThreadQueue<recv_fpga_data_t> &cq_recvdata)
{
	Logging::set(LogLevel::INFO, "--- recv_fpga_deq_thread start ---");

	try {
		std::ostringstream debugstr;
		debugstr << "CH(" << FPGA_CH_ID << ") DMA TX dma_info: dir(" << deqdmainfo.dir << ") chid(" << deqdmainfo.chid << ") ";
		debugstr << "queue_addr(" << std::hex << deqdmainfo.queue_addr <<  ") queue_size(" << std::dec << deqdmainfo.queue_size << ")";
		Logging::set(LogLevel::INFO, debugstr.str());
		debugstr.str("");

		size_t frame_cnt = 0;

		while (true) {
			if (decodethreadque_polling<dmacmd_info_t>(cq)) {
				// Processing complete if there is no input data queue and the decode thread completes
				break;
			}
			// Data to transfer to recv_data_thread
			recv_fpga_data_t recvdata;

			// Get dmacmd to transfer from FPGA to host
			dmacmd_info_t deqdmacmdinfo = cq.pop();
			frame_cnt++;
#if defined(DEBUG_PRINT)
			std::cout << "recv fpga deq: " << frame_cnt << std::endl;
#endif //DEBUG_PRINT

			//-----------------------------------------------------
			// fpga_dequeue
			//-----------------------------------------------------
			uint32_t enq_id = frame_cnt - 1;
			debugstr << std::dec << "CH(" << FPGA_CH_ID << ") DMA TX dmacmd_info: enq(" << enq_id << ") ";
			debugstr << "task_id(" << deqdmacmdinfo.task_id << ") ";
			debugstr << "dst_len(" << deqdmacmdinfo.data_len << ") ";
			debugstr << "dst_addr(" << std::hex << deqdmacmdinfo.data_addr << ")";
			Logging::set(LogLevel::INFO, debugstr.str());
			debugstr.str("");
			int32_t ret = wait_dma_fpga_dequeue(deqdmainfo, deqdmacmdinfo, enq_id, WAIT_TIME_DMA_TX_DEQUEUE);
			if (ret < 0) {
				Logging::set(LogLevel::ERROR, "DMA TX deqerror CH(" + std::to_string(FPGA_CH_ID) + ") enq(" + std::to_string(enq_id) + ")");
			}
			debugstr << std::dec << "CH(" << FPGA_CH_ID << ") DMA TX dmacmd_info: enq(" << enq_id << ") ";
		       	debugstr << "result_task_id(" << deqdmacmdinfo.result_task_id << ") ";
		       	debugstr << "result_status(" << deqdmacmdinfo.result_status << ") ";
		       	debugstr << "result_data_len(" << deqdmacmdinfo.result_data_len << ")";
			Logging::set(LogLevel::INFO, debugstr.str());
			debugstr.str("");

			//-----------------------------------------------------
			// received frame header
			//-----------------------------------------------------
			// Frame header area of DST shared memory
			size_t head_len = sizeof(FrameHeader_t);
			void *head_addr = deqdmacmdinfo.data_addr;
			FrameHeader_t *fh = static_cast<FrameHeader_t*>(head_addr);
			// Get frame header from DST shared memory
			memcpy(&recvdata.frame_header, fh, head_len);	

			//-----------------------------------------------------
			// incoming payload data
			//-----------------------------------------------------
			// DST Shared Memory Payload Area
			void *payload_addr = (void*)((uint64_t)deqdmacmdinfo.data_addr + head_len);
			// Retrieving Image Data from DST Shared Memory
			cv::Mat recv_mat(FPGA_OUT_FRAME_HEIGHT, FPGA_OUT_FRAME_WIDTH, CV_MAKETYPE(CV_8U, FRAME_CHANNEL), payload_addr);
			recv_mat.copyTo(recvdata.mat);

			//-----------------------------------------------------
			// Clear the DST shared memory header area 0xFF after data acquisition
			//-----------------------------------------------------
			memset(fh, 0xFF, sizeof(FrameHeader_t));

			//-----------------------------------------------------
			// interthread data transfer
			//-----------------------------------------------------
			// Forward to recv_data_thread
			cq_recvdata.push(recvdata);

		}

	} catch (const std::exception& error) {
		Logging::set(LogLevel::ERROR, error.what());
	} catch (...) {
		Logging::set(LogLevel::ERROR, "Unknown/internal exception.");
	}

	// thread termination notification
	is_recv_fpga_deq_thread_fin.store(true, std::memory_order_release);

	Logging::set(LogLevel::INFO, "--- recv_fpga_deq_thread finish ---");
#if defined(DEBUG_PRINT)
	std::cout << "recv_fpga_deq_thread finish" << std::endl;
#endif //DEBUG_PRINT
}

void recv_fpga_enq_thread(ConnectThreadQueue<dmacmd_info_t> &cq)
{
	Logging::set(LogLevel::INFO, "--- recv_fpga_enq_thread start ---");

	try {
		size_t frame_cnt = 0;

		while (true) {
			if (decodethreadque_polling<dmacmd_info_t>(cq)) {
				// Processing complete if there is no input data queue and the decode thread completes
				break;
			}

			// Get dmacmd to transfer from FPGA to host
			dmacmd_info_t deqdmacmdinfo = cq.pop();
			frame_cnt++;
#if defined(DEBUG_PRINT)
			std::cout << "recv fpga enq: " << frame_cnt << std::endl;
#endif //DEBUG_PRINT

			//-----------------------------------------------------
			// fpga_enqueue
			//-----------------------------------------------------
			uint32_t enq_id = frame_cnt - 1;
			int32_t ret = wait_dma_fpga_enqueue(deqdmainfo, deqdmacmdinfo, enq_id, WAIT_TIME_DMA_TX_ENQUEUE);
			if (ret < 0) {
				Logging::set(LogLevel::ERROR, "DMA TX enqerror CH(" + std::to_string(FPGA_CH_ID) + ") enq(" + std::to_string(enq_id) + ")");
			}

		}

	} catch (const std::exception& error) {
		Logging::set(LogLevel::ERROR, error.what());
	} catch (...) {
		Logging::set(LogLevel::ERROR, "Unknown/internal exception.");
	}

	Logging::set(LogLevel::INFO, "--- recv_fpga_enq_thread finish ---");
#if defined(DEBUG_PRINT)
	std::cout << "recv_fpga_enq_thread finish" << std::endl;
#endif //DEBUG_PRINT
}

void send_fpga_deq_thread(ConnectThreadQueue<dmacmd_info_t> &cq)
{
	Logging::set(LogLevel::INFO, "--- send_fpga_deq_thread start ---");

	try {
		std::ostringstream debugstr;
		debugstr << "CH(" << FPGA_CH_ID << ") DMA RX dma_info: dir(" << enqdmainfo.dir << ") chid(" << enqdmainfo.chid << ") ";
		debugstr << "queue_addr(" << std::hex << enqdmainfo.queue_addr <<  ") queue_size(" << std::dec << enqdmainfo.queue_size << ")";
		Logging::set(LogLevel::INFO, debugstr.str());
		debugstr.str("");

		size_t frame_cnt = 0;

		while (true) {
			if (decodethreadque_polling<dmacmd_info_t>(cq)) {
				// Processing complete if there is no input data queue and the decode thread completes
				break;
			}

			// dmacmd get to transfer decoded frames from host to FPGA
			dmacmd_info_t enqdmacmdinfo = cq.pop();
			frame_cnt++;
#if defined(DEBUG_PRINT)
			std::cout << "send fpga deq: " << frame_cnt << std::endl;
#endif //DEBUG_PRINT

			//-----------------------------------------------------
			// fpga_dequeue
			//-----------------------------------------------------
			uint32_t enq_id = frame_cnt - 1;
			debugstr << std::dec << "CH(" << FPGA_CH_ID << ") DMA RX dmacmd_info: enq(" << enq_id << ") ";
			debugstr << "task_id(" << enqdmacmdinfo.task_id << ") ";
			debugstr << "src_len(" << enqdmacmdinfo.data_len << ") ";
			debugstr << "src_addr(" << std::hex << enqdmacmdinfo.data_addr << ")";
			Logging::set(LogLevel::INFO, debugstr.str());
			debugstr.str("");
			int32_t ret = wait_dma_fpga_dequeue(enqdmainfo, enqdmacmdinfo, enq_id, WAIT_TIME_DMA_RX_DEQUEUE);
			if (ret < 0) {
				Logging::set(LogLevel::ERROR, "DMA RX deqerror CH(" + std::to_string(FPGA_CH_ID) + ") enq(" + std::to_string(enq_id) + ")");
			}
			debugstr << std::dec << "CH(" << FPGA_CH_ID << ") DMA RX dmacmd_info: enq(" << enq_id << ") ";
		       	debugstr << "result_task_id(" << enqdmacmdinfo.result_task_id << ") ";
		       	debugstr << "result_status(" << enqdmacmdinfo.result_status << ") ";
		       	debugstr << "result_data_len(" << enqdmacmdinfo.result_data_len << ")";
			Logging::set(LogLevel::INFO, debugstr.str());
			debugstr.str("");

		}

	} catch (const std::exception& error) {
		Logging::set(LogLevel::ERROR, error.what());
	} catch (...) {
		Logging::set(LogLevel::ERROR, "Unknown/internal exception.");
	}

	Logging::set(LogLevel::INFO, "--- send_fpga_deq_thread finish ---");
#if defined(DEBUG_PRINT)
	std::cout << "send_fpga_deq_thread finish" << std::endl;
#endif //DEBUG_PRINT
}

void send_fpga_enq_thread(ConnectThreadQueue<send_fpga_info_t> &cq)
{
	Logging::set(LogLevel::INFO, "--- send_fpga_enq_thread start ---");

	try {
		// FPS period
		int64_t fps_nsec = (static_cast<double>(1000000000)/FRAME_FPS); // nsec

		size_t frame_cnt = 0;

		while (true) {
			auto t1 = std::chrono::high_resolution_clock::now();

			if (decodethreadque_polling<send_fpga_info_t>(cq)) {
				// Processing complete if there is no input data queue and the decode thread completes
				break;
			}

			// Get decoded frames and dmacmd transfer decoded frames from host to FPGA
			send_fpga_info_t sendfpgainfo = cq.pop();
			frame_cnt++;
#if defined(DEBUG_PRINT)
			std::cout << "send fpga enq frame: " << frame_cnt << std::endl;
#endif //DEBUG_PRINT

			//-----------------------------------------------------
			// transmit frame header
			//-----------------------------------------------------
			size_t head_len = sizeof(FrameHeader_t);
			FrameHeader_t fh;
#ifdef MODULE_FPGA
			fh.marker          = 0xE0FF10AD;
			fh.payload_len     = FRAME_WIDTH * FRAME_HEIGHT * FRAME_CHANNEL;
			fh.sequence_num    = frame_cnt;
			fh.timestamp       = 0.0;
			fh.data_id         = FPGA_CH_ID;
			fh.header_checksum = 0x0000;
			memset(&fh.reserved1, 0x00, sizeof(fh.reserved1)/sizeof(uint8_t));
			memset(&fh.reserved2, 0x00, sizeof(fh.reserved2)/sizeof(uint8_t));
			memset(&fh.reserved3, 0x00, sizeof(fh.reserved3)/sizeof(uint8_t));
			memset(&fh.reserved4, 0x00, sizeof(fh.reserved4)/sizeof(uint8_t));
#else
			fh.marker       = 0xE0FF10AD;
			fh.payload_len  = FRAME_WIDTH * FRAME_HEIGHT * FRAME_CHANNEL;
			fh.payload_type = 0x01;
			fh.reserved1    = 0x00;
			fh.channel_id   = FPGA_CH_ID % 8;
			fh.frame_index  = frame_cnt;
			fh.color_space  = 0x01;
			fh.data_type    = 0x00;
			fh.num_ch       = 0x03;
			fh.width        = FRAME_WIDTH;
			fh.height       = FRAME_HEIGHT;
			fh.local_ts     = 0.0;
			memset(&fh.reserved2, 0x00, sizeof(fh.reserved2)/sizeof(uint8_t));
#endif //MODULE_FPGA
			// SRC shared memory frame header area
#if defined(ALLOCATE_SRC_SHMEM_MODE)
			void *head_addr = sendfpgainfo.dmacmdinfo.data_addr;
#else
			void *head_addr = FPGA_SRC_SHMEM_ADDR;
#endif
			// Set frame header in SRC shared memory
			memcpy(head_addr, &fh, head_len);

			//-----------------------------------------------------
			// outgoing payload data
			//-----------------------------------------------------
			// SRC Shared Memory Payload Area
#if defined(ALLOCATE_SRC_SHMEM_MODE)
			void *payload_addr = (void*)((uint64_t)sendfpgainfo.dmacmdinfo.data_addr + head_len);
#else
			void *payload_addr = (void*)((uint64_t)FPGA_SRC_SHMEM_ADDR + head_len);
#endif
			// Set decoded image data in SRC shared memory
			memcpy(payload_addr, sendfpgainfo.mat.ptr(), fh.payload_len);

			//-----------------------------------------------------
			// fpga_enqueue
			//-----------------------------------------------------
			uint32_t enq_id = frame_cnt - 1;
#if defined(EXEC_FPGA_ENQUEUE_MODE)
			int32_t ret = wait_dma_fpga_enqueue(enqdmainfo, sendfpgainfo.dmacmdinfo, enq_id, WAIT_TIME_DMA_RX_ENQUEUE);
			if (ret < 0) {
				Logging::set(LogLevel::ERROR, "DMA RX enqerror CH(" + std::to_string(FPGA_CH_ID) + ") enq(" + std::to_string(enq_id) + ")");
			}
#endif

			auto t2 = std::chrono::high_resolution_clock::now();
			auto duration = (std::chrono::duration_cast<std::chrono::nanoseconds>(t2 - t1)).count();
			if (fps_nsec > duration) {
				// FPS time lapse
				auto sleep_time = fps_nsec - duration;
				Logging::set(LogLevel::DEBUG, "CH(" + std::to_string(FPGA_CH_ID) + ") DMA RX enq(" + std::to_string(enq_id) + ") duration time: " + std::to_string(duration) + " fps_nsec (<" + std::to_string(fps_nsec) + "), sleep_time: " + std::to_string(sleep_time));
				std::this_thread::sleep_for(std::chrono::nanoseconds(sleep_time));
			}
		}

	} catch (const std::exception& error) {
		Logging::set(LogLevel::ERROR, error.what());
	} catch (...) {
		Logging::set(LogLevel::ERROR, "Unknown/internal exception.");
	}

	Logging::set(LogLevel::INFO, "--- send_fpga_enq_thread finish ---");
#if defined(DEBUG_PRINT)
	std::cout << "send_fpga_enq_thread finish" << std::endl;
#endif //DEBUG_PRINT
}

void send_tcp_thread(int32_t &sockfd, ConnectThreadQueue<cv::Mat> &cq)
{
	Logging::set(LogLevel::INFO, "--- send_tcp_thread start ---");

	try {
		// FPS period
		int64_t fps_nsec = (static_cast<double>(1000000000)/FRAME_FPS); // nsec

		int32_t data_size = sizeof(FrameHeader_t) + FRAME_WIDTH * FRAME_HEIGHT * FRAME_CHANNEL;
		uint8_t *data_addr = (uint8_t*)malloc(data_size);

		size_t frame_cnt = 0;
		while (true) {
			auto t1 = std::chrono::high_resolution_clock::now();

			if (decodethreadque_polling<cv::Mat>(cq)) {
				// Processing complete if there is no input data queue and the decode thread completes
				break;
			}

			// decode frame acquisition
			cv::Mat sendmat = cq.pop();
			frame_cnt++;
#if defined(DEBUG_PRINT)
			std::cout << "send tcp frame: " << frame_cnt << std::endl;
#endif //DEBUG_PRINT

			//-----------------------------------------------------
			// transmit frame header
			//-----------------------------------------------------
			size_t head_len = sizeof(FrameHeader_t);
			FrameHeader_t fh;
#ifdef MODULE_FPGA
			fh.marker          = 0xE0FF10AD;
			fh.payload_len     = FRAME_WIDTH * FRAME_HEIGHT * FRAME_CHANNEL;
			fh.sequence_num    = frame_cnt;
			fh.timestamp       = 0.0;
			fh.data_id         = FPGA_CH_ID;
			fh.header_checksum = 0x0000;
			memset(&fh.reserved1, 0x00, sizeof(fh.reserved1)/sizeof(uint8_t));
			memset(&fh.reserved2, 0x00, sizeof(fh.reserved2)/sizeof(uint8_t));
			memset(&fh.reserved3, 0x00, sizeof(fh.reserved3)/sizeof(uint8_t));
			memset(&fh.reserved4, 0x00, sizeof(fh.reserved4)/sizeof(uint8_t));
#else
			fh.marker       = 0xE0FF10AD;
			fh.payload_len  = FRAME_WIDTH * FRAME_HEIGHT * FRAME_CHANNEL;
			fh.payload_type = 0x01;
			fh.reserved1    = 0x00;
			fh.channel_id   = FPGA_CH_ID % 8;
			fh.frame_index  = frame_cnt;
			fh.color_space  = 0x01;
			fh.data_type    = 0x00;
			fh.num_ch       = 0x03;
			fh.width        = FRAME_WIDTH;
			fh.height       = FRAME_HEIGHT;
			fh.local_ts     = 0.0;
			memset(&fh.reserved2, 0x00, sizeof(fh.reserved2)/sizeof(uint8_t));
#endif //MODULE_FPGA

			// Frame header area of the transmit memory
			void *head_addr = data_addr;

			// Set frame header in transmit memory
			memcpy(head_addr, &fh, head_len);

			//-----------------------------------------------------
			// outgoing payload data
			//-----------------------------------------------------
			// payload area of the transmit memory
			void *payload_addr = data_addr + head_len;

			// Set the decoded image data in the transmission memory
			memcpy(payload_addr, sendmat.ptr(), fh.payload_len);

			//-----------------------------------------------------
			// TCP Send to Server
			//-----------------------------------------------------
			uint32_t enq_id = frame_cnt - 1;
			Logging::set(LogLevel::DEBUG, "CH(" + std::to_string(FPGA_CH_ID) + ") TCP send enq(" + std::to_string(enq_id) + ")");
			int32_t send_ret = send(sockfd, data_addr, data_size, 0);
			if (send_ret < 0) {
				Logging::set(LogLevel::ERROR, "TCP send error. errno(" + std::to_string(errno) + ": " + strerror(errno) + ")");
			} else if (send_ret < data_size) {
				// If the specified size cannot be sent, send the rest
				int32_t send_ret2 = send(sockfd, (data_addr + send_ret), (data_size - send_ret), 0);
				if (send_ret2 < 0) {
					Logging::set(LogLevel::ERROR, "TCP send error. errno(" + std::to_string(errno) + ": " + strerror(errno) + ")");
				}
			}

			sendmat.release();

			auto t2 = std::chrono::high_resolution_clock::now();
			auto duration = (std::chrono::duration_cast<std::chrono::nanoseconds>(t2 - t1)).count();
			if (fps_nsec > duration) {
				// FPS time lapse
				auto sleep_time = fps_nsec - duration;
				Logging::set(LogLevel::DEBUG, "CH(" + std::to_string(FPGA_CH_ID) + ") TCP send enq(" + std::to_string(enq_id) + ") duration time: " + std::to_string(duration) + " fps_nsec (<" + std::to_string(fps_nsec) + "), sleep_time: " + std::to_string(sleep_time));
				std::this_thread::sleep_for(std::chrono::nanoseconds(sleep_time));
			}
		}

		free(data_addr);

	} catch (const std::exception& error) {
		Logging::set(LogLevel::ERROR, error.what());
	} catch (...) {
		Logging::set(LogLevel::ERROR, "Unknown/internal exception.");
	}

	Logging::set(LogLevel::INFO, "--- send_tcp_thread finish ---");
#if defined(DEBUG_PRINT)
	std::cout << "send_tcp_thread finish" << std::endl;
#endif //DEBUG_PRINT
}

void decode_to_tcp_thread(cv::VideoCapture &cap, ConnectThreadQueue<cv::Mat> &cq_send)
{
	Logging::set(LogLevel::INFO, "--- decode_to_tcp_thread start ---");

	try {
		size_t frame_cnt = 0;

		while (true) {
			// Data to transfer to send_tcp_thread
			cv::Mat sendmat;

			//-----------------------------------------------------
			// Read (Decode) frames from video distribution source
			//-----------------------------------------------------
			if (!cap.read(sendmat)) {
				if (sendmat.empty()) {
					break;
				} else {
					throw std::runtime_error("Failed to get frame from cv::VideoCapture");
				}
			}
			frame_cnt++;
#if defined(DEBUG_PRINT)
			std::cout << "decode frame: " << frame_cnt << std::endl;
#endif //DEBUG_PRINT

			//-----------------------------------------------------
			// interthread data transfer
			//-----------------------------------------------------
			// Forward to send_tcp_thread
			cq_send.push(sendmat);

			sendmat.release();
		}

	} catch (const std::exception& error) {
		Logging::set(LogLevel::ERROR, error.what());
	} catch (...) {
		Logging::set(LogLevel::ERROR, "Unknown/internal exception.");
	}

	// decode thread termination notification
	is_decode_thread_fin.store(true, std::memory_order_release);

	Logging::set(LogLevel::INFO, "--- decode_to_tcp_thread finish ---");
#if defined(DEBUG_PRINT)
	std::cout << "decode_to_tcp_thread finish" << std::endl;
#endif //DEBUG_PRINT
}

void decode_to_dma_thread(cv::VideoCapture &cap, const mngque_t &mngq, ConnectThreadQueue<send_fpga_info_t> &cq_sendenq, ConnectThreadQueue<dmacmd_info_t> &cq_senddeq, ConnectThreadQueue<dmacmd_info_t> &cq_recvenq, ConnectThreadQueue<dmacmd_info_t> &cq_recvdeq)
{
	Logging::set(LogLevel::INFO, "--- decode_to_dma_thread start ---");

	try {
		size_t frame_cnt = 0;

		while (true) {
			// data to send_fpga_enq_thread
			send_fpga_info_t sendfpgainfo;

			//-----------------------------------------------------
			// Read (Decode) frames from video distribution source
			//-----------------------------------------------------
			if (!cap.read(sendfpgainfo.mat)) {
				if (sendfpgainfo.mat.empty()) {
					break;
				} else {
					throw std::runtime_error("Failed to get frame from cv::VideoCapture");
				}
			}
			frame_cnt++;
#if defined(DEBUG_PRINT)
			std::cout << "decode frame: " << frame_cnt << std::endl;
#endif //DEBUG_PRINT

			//-----------------------------------------------------
			// set dmacmd info
			//-----------------------------------------------------
			uint32_t enq_id = frame_cnt - 1;
			// dmacmd configuration to transfer from host to FPGA
			dmacmd_info_t enqdmacmdinfo;
#if defined(EXEC_FPGA_ENQUEUE_MODE)
			if (tp_enqueue_set_dma_cmd(enq_id, mngq, enqdmacmdinfo) < 0) {
				Logging::set(LogLevel::ERROR, "tp_enqueue_set_dma_cmd error!!");
			}
			memcpy(&sendfpgainfo.dmacmdinfo, &enqdmacmdinfo, sizeof(dmacmd_info_t));
#endif

			// FPGA-to-host dmacmd configuration
			dmacmd_info_t deqdmacmdinfo;
#if defined(EXEC_FPGA_DEQUEUE_MODE)
			if (tp_dequeue_set_dma_cmd(enq_id, mngq, deqdmacmdinfo) < 0) {
				Logging::set(LogLevel::ERROR, "tp_dequeue_set_dma_cmd error!!");
			}
#endif

			//-----------------------------------------------------
			// interthread data transfer
			//-----------------------------------------------------
#if defined(EXEC_FPGA_DEQUEUE_MODE)
			// Forward to recv_fpga_deq_thread
			cq_recvdeq.push(deqdmacmdinfo);
			// Forward to recv_fpga_enq_thread
			cq_recvenq.push(deqdmacmdinfo);
#endif
#if defined(EXEC_FPGA_ENQUEUE_MODE)
			// send_fpga_deq_thread
			cq_senddeq.push(enqdmacmdinfo);
#endif
			// send_fpga_enq_thread
			cq_sendenq.push(sendfpgainfo);
		}

	} catch (const std::exception& error) {
		Logging::set(LogLevel::ERROR, error.what());
	} catch (...) {
		Logging::set(LogLevel::ERROR, "Unknown/internal exception.");
	}

	// decode thread termination notification
	is_decode_thread_fin.store(true, std::memory_order_release);

	Logging::set(LogLevel::INFO, "--- decode_to_dma_thread finish ---");
#if defined(DEBUG_PRINT)
	std::cout << "decode_to_dma_thread finish" << std::endl;
#endif //DEBUG_PRINT
}

void prlog_mngque(const mngque_t &m)
{
	Logging::set(LogLevel::DEBUG, "prlog_mngque...");

	Logging::set(LogLevel::DEBUG, "  mngque.srcdsize(" + std::to_string(m.srcdsize) + ")");
	Logging::set(LogLevel::DEBUG, "  mngque.dstdsize(" + std::to_string(m.dstdsize) + ")");
	Logging::set(LogLevel::DEBUG, "  mngque.srcbuflen(" + std::to_string(m.srcbuflen) + ")");
	Logging::set(LogLevel::DEBUG, "  mngque.dstbuflen(" + std::to_string(m.dstbuflen) + ")");
	std::ostringstream bufpval;
	bufpval << std::hex << m.srcbufp;
	Logging::set(LogLevel::DEBUG, "  mngque.srcbufp(" + bufpval.str() + ")");
	bufpval.str("");
	bufpval << std::hex << m.dstbufp;
	Logging::set(LogLevel::DEBUG, "  mngque.dstbufp(" + bufpval.str() + ")");
	bufpval.str("");
}

void prlog_dma_info(const dma_info_t &i)
{
	Logging::set(LogLevel::DEBUG, "prlog_dma_info...");

	Logging::set(LogLevel::DEBUG, "  dev_id(" + std::to_string(i.dev_id) + ")");
	Logging::set(LogLevel::DEBUG, "  dir(" + std::to_string(i.dir) + ")");
	Logging::set(LogLevel::DEBUG, "  chid(" + std::to_string(i.chid) + ")");
	std::ostringstream bufpval;
	bufpval << std::hex << i.queue_addr;
	Logging::set(LogLevel::DEBUG, "  queue_addr(" + bufpval.str() + ")");
	bufpval.str("");
	Logging::set(LogLevel::DEBUG, "  queue_size(" + std::to_string(i.queue_size) + ")");
	if (i.connector_id == nullptr) {
		Logging::set(LogLevel::DEBUG, "  connector_id()");
	} else {
		Logging::set(LogLevel::DEBUG, "  connector_id(" + std::string(i.connector_id) + ")");
	}
}

void prlog_dmacmd_info(const dmacmd_info_t &i, uint32_t enq_id)
{
	Logging::set(LogLevel::DEBUG, "prlog_dmacmd_info...");

	std::string enq_id_str = std::to_string(enq_id);
	Logging::set(LogLevel::DEBUG, "  ENQ(" + enq_id_str + ") task_id(" + std::to_string(i.task_id) + ")");
	Logging::set(LogLevel::DEBUG, "  ENQ(" + enq_id_str + ") data_len(" + std::to_string(i.data_len) + ")");
	std::ostringstream bufpval;
	bufpval << std::hex << i.data_addr;
	Logging::set(LogLevel::DEBUG, "  ENQ(" + enq_id_str + ") data_addr(" + bufpval.str() + ")");
	bufpval.str("");
	bufpval << std::hex << i.desc_addr;
	Logging::set(LogLevel::DEBUG, "  ENQ(" + enq_id_str + ") desc_addr(" + bufpval.str() + ")");
	bufpval.str("");
	Logging::set(LogLevel::DEBUG, "  ENQ(" + enq_id_str + ") result_status(" + std::to_string(i.result_status) + ")");
	Logging::set(LogLevel::DEBUG, "  ENQ(" + enq_id_str + ") result_task_id(" + std::to_string(i.result_task_id) + ")");
	Logging::set(LogLevel::DEBUG, "  ENQ(" + enq_id_str + ") result_data_len(" + std::to_string(i.result_data_len) + ")");
	bufpval << std::hex << i.result_data_addr;
	Logging::set(LogLevel::DEBUG, "  ENQ(" + enq_id_str + ") result_data_addr(" + bufpval.str() + ")");
	bufpval.str("");
}

void prlog_fpga_chain_ddr_info(const fpga_chain_ddr_t &i, uint32_t dev_id, uint32_t krnl_id, uint32_t extif_id)
{
	Logging::set(LogLevel::DEBUG, "prlog_fpga_chain_ddr_info dev(" + std::to_string(dev_id) + ") kernel(" + std::to_string(krnl_id) + ") extif(" + std::to_string(extif_id) + ")");
	std::ostringstream bufpval;
	bufpval << std::hex << i.base;
	Logging::set(LogLevel::DEBUG, "  base      0x" + bufpval.str());
	bufpval.str("");
	bufpval << std::hex << i.rx_offset;
	Logging::set(LogLevel::DEBUG, "  rx_offset 0x" + bufpval.str());
	bufpval.str("");
	bufpval << std::hex << i.rx_stride;
	Logging::set(LogLevel::DEBUG, "  rx_stride 0x" + bufpval.str());
	bufpval.str("");
	bufpval << std::hex << static_cast<uint32_t>(i.rx_size);
	Logging::set(LogLevel::DEBUG, "  rx_size   0x" + bufpval.str());
	bufpval.str("");
	bufpval << std::hex << i.tx_offset;
	Logging::set(LogLevel::DEBUG, "  tx_offset 0x" + bufpval.str());
	bufpval.str("");
	bufpval << std::hex << i.tx_stride;
	Logging::set(LogLevel::DEBUG, "  tx_stride 0x" + bufpval.str());
	bufpval.str("");
	bufpval << std::hex << static_cast<uint32_t>(i.tx_size);
	Logging::set(LogLevel::DEBUG, "  tx_size   0x" + bufpval.str());
	bufpval.str("");
}

int main(void)
{
	try {
		int32_t ret = 0;

		libfpga_log_set_level(LIBFPGA_LOG_ALL);

#if defined(APPLOG_PRINT)
		// Log standard output setting
		Logging::set_stdoutmode(true);
#endif //APPLOG_PRINT

		// Log Level Settings
		uint8_t param_applog_level;
		if (env2num("DECENV_APPLOG_LEVEL", param_applog_level, 0) != -1) {
			Logging::setlevel(static_cast<LogLevel>(param_applog_level));
		}

		// Show Version
		Logging::set(LogLevel::PRINT, "Version: " + std::string(VERSION));

		// interthread data transfer
		ConnectThreadQueue<send_fpga_info_t> cq_send_enq;
		ConnectThreadQueue<dmacmd_info_t> cq_send_deq;
		ConnectThreadQueue<dmacmd_info_t> cq_recv_enq;
		ConnectThreadQueue<dmacmd_info_t> cq_recv_deq;
		ConnectThreadQueue<recv_fpga_data_t> cq_recv_data;
		ConnectThreadQueue<cv::Mat> cq_send;

		// Thread
		std::thread decode_th;
		std::thread send_fpga_enq_th;
		std::thread send_fpga_deq_th;
		std::thread recv_fpga_enq_th;
		std::thread recv_fpga_deq_th;
		std::thread recv_data_th;
		std::thread send_tcp_th;

		// video capture count
		size_t cap_cnt = 0;

		// Variable Initialization for FPGA
		uint32_t dev_id = 0;
		mngque_t mngq;
		mngq.srcdsize = 0;
		mngq.srcbuflen = 0;
		mngq.dstdsize = 0;
		mngq.dstbuflen = 0;
		mngq.srcbufp = nullptr;
		mngq.dstbufp = nullptr;

		// TCP variable initialization
		int32_t sockfd;

		//--------------------------------------------------------------------
		// Parameter Settings
		//--------------------------------------------------------------------
		// input video delivery protocol: required environment variables
		std::string param_videosrc_protocol;
		env2str("DECENV_VIDEOSRC_PROTOCOL", param_videosrc_protocol, 1);
		Logging::set(LogLevel::INFO, "parameter VIDEOSRC_PROTOCOL: \"" + param_videosrc_protocol + "\"");
		VIDEOSRC_PROTOCOL =
			param_videosrc_protocol == "RTP" ? videosrc_protocol::RTP :
			param_videosrc_protocol == "RTSP" ? videosrc_protocol::RTSP :
			throw std::runtime_error("VIDEOSRC_PROTOCOL: \"" + param_videosrc_protocol + "\" is invalid.");

		// ingress video distribution IPv4 ports: required environment variables
		std::string param_videosrc_port;
		env2str("DECENV_VIDEOSRC_PORT", param_videosrc_port, 1);

		// video distribution source
		std::string videosrc;
		if (VIDEOSRC_PROTOCOL == videosrc_protocol::RTP) {
			// RTP
			videosrc = "udpsrc port=" + param_videosrc_port + " buffer-size=512000000 caps=application/x-rtp ! rtpjitterbuffer ! rtph264depay ! h264parse ! openh264dec ! queue ! videoconvert ! appsink";
		} else if (VIDEOSRC_PROTOCOL == videosrc_protocol::RTSP) {
			// video distribution RTSP server IP address: required environment variable
			std::string param_videosrc_ipa;
			env2str("DECENV_VIDEOSRC_IPA", param_videosrc_ipa, 1);

			// RTSP
			std::string url = "rtsp://" + param_videosrc_ipa + ":" + param_videosrc_port + "/test";
			videosrc = "rtspsrc location=" + url + " ! application/x-rtp ! rtpjitterbuffer ! rtph264depay ! h264parse ! openh264dec ! queue ! videoconvert ! appsink";
		}
		Logging::set(LogLevel::INFO, "parameter videosrc: \"" + videosrc + "\"");

		// frame FPS: required environment variables
		env2num("DECENV_FRAME_FPS", FRAME_FPS, 1);
		Logging::set(LogLevel::INFO, "parameter FRAME_FPS: " + std::to_string(FRAME_FPS));

		// frame size: required environment variable
		env2num("DECENV_FRAME_WIDTH", FRAME_WIDTH, 1);
		Logging::set(LogLevel::INFO, "parameter FRAME_WIDTH: " + std::to_string(FRAME_WIDTH));
		env2num("DECENV_FRAME_HEIGHT", FRAME_HEIGHT, 1);
		Logging::set(LogLevel::INFO, "parameter FRAME_HEIGHT: " + std::to_string(FRAME_HEIGHT));

		// delivery video connection limit: any environment variable
		uint32_t param_video_connect_limit;
		ret = env2num("DECENV_VIDEO_CONNECT_LIMIT", param_video_connect_limit, 0);
		if (ret != -1) {
			VIDEO_CONNECT_LIMIT = param_video_connect_limit;
		}
		Logging::set(LogLevel::INFO, "parameter VIDEO_CONNECT_LIMIT: " + std::to_string(VIDEO_CONNECT_LIMIT));

		// output protocols: required environment variables
		std::string param_outdst_protocol;
		env2str("DECENV_OUTDST_PROTOCOL", param_outdst_protocol, 1);
		Logging::set(LogLevel::INFO, "parameter OUTDST_PROTOCOL: \"" + param_outdst_protocol + "\"");
		OUTDST_PROTOCOL =
			param_outdst_protocol == "DMA" ? protocol::DMA :
			param_outdst_protocol == "TCP" ? protocol::TCP :
			throw std::runtime_error("OUTDST_PROTOCOL: \"" + param_outdst_protocol + "\" is invalid.");

		std::string param_dpdk_file_prefix;
		std::string param_fpga_dma_host_to_dev_connector_id;
		std::string param_fpga_dev_name;
		if (OUTDST_PROTOCOL == protocol::DMA) {
			// DPDK file prefixes: arbitrary environment variables
#if defined(DPDK_SECONDARY_PROC_MODE)
			env2str("DECENV_DPDK_FILE_PREFIX", param_dpdk_file_prefix, 1);
			DPDK_FILE_PREFIX = param_dpdk_file_prefix.c_str();
			Logging::set(LogLevel::INFO, "parameter DPDK_FILE_PREFIX: \"" + std::string(DPDK_FILE_PREFIX) + "\"");
#endif

			// connector ID: required environment variable
			env2str("DECENV_FPGA_DMA_HOST_TO_DEV_CONNECTOR_ID", param_fpga_dma_host_to_dev_connector_id, 1);
			FPGA_DMA_HOST_TO_DEV_CONNECTOR_ID = param_fpga_dma_host_to_dev_connector_id.c_str();
			Logging::set(LogLevel::INFO, "parameter FPGA_DMA_HOST_TO_DEV_CONNECTOR_ID: \"" + std::string(FPGA_DMA_HOST_TO_DEV_CONNECTOR_ID) + "\"");

			// Destination shared memory address (SRC shared memory on the FPGA): any environment variable
#if !(defined(ALLOCATE_SRC_SHMEM_MODE))
			uint64_t param_fpga_src_shmem_addr;	
			env2num("DECENV_FPGA_SRC_SHMEM_ADDR", param_fpga_src_shmem_addr, 1);
			FPGA_SRC_SHMEM_ADDR = reinterpret_cast<void*>(param_fpga_src_shmem_addr);
			std::ostringstream oss;
			oss.setf(std::ios::hex, std::ios::basefield);
			oss << "parameter FPGA_SRC_SHMEM_ADDR: " << std::hex << FPGA_SRC_SHMEM_ADDR;
			Logging::set(LogLevel::INFO, oss.str());
			oss.str("");
			oss.setf(std::ios::dec, std::ios::basefield);
#endif

			// FPGA device files: arbitrary environment variables
#if defined(CONTROL_FPGA_DEV_INIT) || defined(EXEC_FPGA_DEQUEUE_MODE)
			ret = env2str("DECENV_FPGA_DEV_NAME", param_fpga_dev_name, 0);
			if (ret != -1) {
				FPGA_DEV_NAME = param_fpga_dev_name.c_str();
			}
			Logging::set(LogLevel::INFO, "parameter FPGA_DEV_NAME: \"" + std::string(FPGA_DEV_NAME) + "\"");
#endif

			// FPGA channel ID: any environment variable
#if defined(CONTROL_FPGA_ENQUEUE_LLDMA_INIT) || defined(EXEC_FPGA_DEQUEUE_MODE)
			uint16_t param_fpga_ch_id;
			ret = env2num("DECENV_FPGA_CH_ID", param_fpga_ch_id, 0);
			if (ret != -1) {
				FPGA_CH_ID = param_fpga_ch_id;
			}
			Logging::set(LogLevel::INFO, "parameter FPGA_CH_ID: " + std::to_string(FPGA_CH_ID));
#endif

			// FPGA output frame size: arbitrary environment variable
#if defined(CONTROL_FPGA_FUNC_INIT) || defined(EXEC_FPGA_DEQUEUE_MODE)
			uint32_t param_fpga_out_frame_width;
			ret = env2num("DECENV_FPGA_OUT_FRAME_WIDTH", param_fpga_out_frame_width, 0);
			if (ret != -1) {
				FPGA_OUT_FRAME_WIDTH = param_fpga_out_frame_width;
			}
			Logging::set(LogLevel::INFO, "parameter FPGA_OUT_FRAME_WIDTH: " + std::to_string(FPGA_OUT_FRAME_WIDTH));
			uint32_t param_fpga_out_frame_height;
			ret = env2num("DECENV_FPGA_OUT_FRAME_HEIGHT", param_fpga_out_frame_height, 0);
			if (ret != -1) {
				FPGA_OUT_FRAME_HEIGHT = param_fpga_out_frame_height;
			}
			Logging::set(LogLevel::INFO, "parameter FPGA_OUT_FRAME_HEIGHT: " + std::to_string(FPGA_OUT_FRAME_HEIGHT));
#endif
		} else if (OUTDST_PROTOCOL == protocol::TCP) {
			// TCP output destination IP address: required environment variable
			env2str("DECENV_OUTDST_IPA", OUTDST_IPA, 1);
			Logging::set(LogLevel::INFO, "parameter OUTDST_IPA: " + OUTDST_IPA);

			// TCP output destination IP port: required environment variables
			env2num("DECENV_OUTDST_PORT", OUTDST_PORT, 1);
			Logging::set(LogLevel::INFO, "parameter OUTDST_PORT: " + std::to_string(OUTDST_PORT));

			// FPGA channel ID: any environment variable (for frame headers)
			uint16_t param_fpga_ch_id;
			ret = env2num("DECENV_FPGA_CH_ID", param_fpga_ch_id, 0);
			if (ret != -1) {
				FPGA_CH_ID = param_fpga_ch_id;
			}
			Logging::set(LogLevel::INFO, "parameter FPGA_CH_ID: " + std::to_string(FPGA_CH_ID));
		}


		//--------------------------------------------------------------------
		// Output Settings
		//--------------------------------------------------------------------
		if (OUTDST_PROTOCOL == protocol::DMA) {
			//----------------------------------------------
			// initialize DPDK
			//----------------------------------------------
#if defined(DPDK_SECONDARY_PROC_MODE)
			ret = fpga_shmem_init(DPDK_FILE_PREFIX, 0, 0);
			if (ret < 0) {
				throw std::runtime_error("Initialize DPDK failed!! fpga_shmem_init: ret(" + std::to_string(ret) + ")");
			}
			Logging::set(LogLevel::DEBUG, "fpga_shmem_init: ret(" + std::to_string(ret) + ")");
#else
			ret = fpga_shmem_init_sys(0, 0, 0, 0, 0);
			if (ret < 0) {
				throw std::runtime_error("Initialize DPDK failed!! fpga_shmem_init_sys: ret(" + std::to_string(ret) + ")");
			}
			Logging::set(LogLevel::DEBUG, "fpga_shmem_init_sys: ret(" + std::to_string(ret) + ")");
#endif

			//----------------------------------------------
			// set CPU affinity
			//----------------------------------------------
			cpu_set_t mask;
			size_t len = sizeof(mask);
			pid_t pid = getpid();
			CPU_ZERO(&mask);
			int64_t ncores = sysconf(_SC_NPROCESSORS_CONF);
			for (int64_t i=0; i<ncores; i++) {
				CPU_SET(i, &mask);
			}
			ret = sched_setaffinity(pid, len, &mask);
			if (ret < 0) {
				throw std::runtime_error("Set CPU affinity failed!! sched_setaffinity: ret(" + std::to_string(ret) + ")");
			}
			Logging::set(LogLevel::DEBUG, "sched_setaffinity: ret(" + std::to_string(ret) + ")");
			std::ostringstream oss_bits;
			oss_bits.setf(std::ios::hex, std::ios::basefield);
			oss_bits << std::hex << *mask.__bits;
			Logging::set(LogLevel::INFO, "PID(" + std::to_string(pid) + ") CPU affinity(0x" + oss_bits.str() + ")");
			oss_bits.str("");
			oss_bits.setf(std::ios::dec, std::ios::basefield);

			//----------------------------------------------
			// initialize FPGA
			//----------------------------------------------
#if defined(CONTROL_FPGA_DEV_INIT)
			ret = fpga_dev_init(FPGA_DEV_NAME, &dev_id);
			if (ret < 0) {
				throw std::runtime_error("fpga_dev_init error!!: ret(" + std::to_string(ret) + ")" );
			}
			Logging::set(LogLevel::DEBUG, "fpga_dev_init: ret(" + std::to_string(ret) + ")");

			int32_t fpga_num = fpga_get_num();
			if (fpga_num != 1) {
				throw std::runtime_error("Num of FPGA error(" +  std::to_string(fpga_num) + ")");
			}
#endif

			//----------------------------------------------
			// fpga enable regrw
			//----------------------------------------------
			ret = fpga_enable_regrw(dev_id);
			if (ret < 0) {
				Logging::set(LogLevel::ERROR, "fpga_enable_regrw error!!: ret(" + std::to_string(ret) + ")");
				goto _END1;
			}
			Logging::set(LogLevel::DEBUG, "fpga_enable_regrw: ret(" + std::to_string(ret) + ")");

			//----------------------------------------------
			// Shared Memory Settings
			//----------------------------------------------
			mngq.srcdsize = sizeof(FrameHeader_t) + FRAME_WIDTH * FRAME_HEIGHT * FRAME_CHANNEL;
			mngq.srcbuflen = (mngq.srcdsize + (ALIGN_BUF_LEN - 1)) & ~(ALIGN_BUF_LEN - 1);
#if defined(ALLOCATE_SRC_SHMEM_MODE)
			Logging::set(LogLevel::INFO, "--- fpga_shmem_alloc for srcbuf ---");
			mngq.srcbufp = fpga_shmem_aligned_alloc(mngq.srcbuflen);
			if (mngq.srcbufp == nullptr) {
				Logging::set(LogLevel::ERROR, "fpga_shmem_alloc error!!: ret(" + std::to_string(ret) + ")");
				goto _END1;
			}
#else
			mngq.srcbufp = FPGA_SRC_SHMEM_ADDR;
#endif

#if defined(EXEC_FPGA_DEQUEUE_MODE)
			mngq.dstdsize = sizeof(FrameHeader_t) + FPGA_OUT_FRAME_WIDTH * FPGA_OUT_FRAME_HEIGHT * FRAME_CHANNEL;
			mngq.dstbuflen = (mngq.dstdsize + (ALIGN_BUF_LEN - 1)) & ~(ALIGN_BUF_LEN - 1);
			Logging::set(LogLevel::INFO, "--- fpga_shmem_alloc for dstbuf ---");
			mngq.dstbufp = fpga_shmem_aligned_alloc(mngq.dstbuflen);
			if (mngq.dstbufp == nullptr) {
				Logging::set(LogLevel::ERROR, "fpga_shmem_alloc error!!: ret(" + std::to_string(ret) + ")");
				goto _END1;
			}
#endif

			prlog_mngque(mngq);

			//----------------------------------------------
			// fpga lldma init
			//----------------------------------------------
#if defined(CONTROL_FPGA_ENQUEUE_LLDMA_INIT)
			// Execute ENQUEUE fpga_lldma_init() in the local process
			if (tp_enqueue_lldma_init(dev_id) < 0) {
				// error
				goto _END1;
			}
#else
			// If another process is running ENQUEUE fpga_lldma_init(),
			// Connector ID setting for environment variable DECENV_FPGA_DMA_HOST_TO_DEV_CONNECTOR_ID
			enqdmainfo_channel.connector_id = const_cast<char*>(FPGA_DMA_HOST_TO_DEV_CONNECTOR_ID);
#endif
#if defined(EXEC_FPGA_DEQUEUE_MODE)
			if (tp_dequeue_lldma_init(dev_id) < 0) {
				// error
				goto _END1;
			}
#endif

			//----------------------------------------------
			// fpga lldma queue setup (set dmainfo)
			//----------------------------------------------
			if (tp_enqueue_lldma_queue_setup() < 0) {
				// error
				goto _END2;
			}
#if !(defined(CONTROL_FPGA_ENQUEUE_LLDMA_INIT))
			// If another process is running ENQUEUE fpga_lldma_init(),
			// Acquires the FPGA CH ID associated with the connector ID of the environment variable DECENV_FPGA_DMA_HOST_TO_DEV_CONNECTOR_ID.
			FPGA_CH_ID = enqdmainfo.chid;
			Logging::set(LogLevel::INFO, "parameter FPGA_CH_ID: " + std::to_string(FPGA_CH_ID));
#else
			// The local process is executing ENQUEUE fpga_lldma_init().
#endif

#if defined(EXEC_FPGA_DEQUEUE_MODE)
			if (tp_dequeue_lldma_queue_setup() < 0) {
				// error
				goto _END2;
			}
#endif

			//----------------------------------------------
			// FPGA kernel init
			//----------------------------------------------
#if defined(CONTROL_FPGA_FUNC_INIT)
			// function filter_resize init
			ret = tp_function_filter_resize_init(dev_id);
			if (ret < 0) {
				// error
				switch (ret) {
					case -2:
						goto _END4;
						break;
					default:
						goto _END3;
				}
			}
#endif

			//----------------------------------------------
			// function chain control
			//----------------------------------------------
#if defined(CONTROL_FPGA_CHAIN_CONNECT)
			if (tp_chain_connect(dev_id) < 0) {
				// error
				goto _END4;
			}
#endif
		}


		//--------------------------------------------------------------------
		// start of video input processing
		//--------------------------------------------------------------------
		while (true) {
			//----------------------------------------------
			// Thread Start
			//----------------------------------------------
			if (OUTDST_PROTOCOL == protocol::DMA) {
#if defined(EXEC_FPGA_DEQUEUE_MODE)
				recv_data_th = std::thread(recv_data_thread, std::ref(cq_recv_data));
				recv_fpga_deq_th = std::thread(recv_fpga_deq_thread, std::ref(cq_recv_deq), std::ref(cq_recv_data));
				recv_fpga_enq_th = std::thread(recv_fpga_enq_thread, std::ref(cq_recv_enq));
#endif
#if defined(EXEC_FPGA_ENQUEUE_MODE)
				send_fpga_deq_th = std::thread(send_fpga_deq_thread, std::ref(cq_send_deq));
#endif
				send_fpga_enq_th = std::thread(send_fpga_enq_thread, std::ref(cq_send_enq));
			} else if (OUTDST_PROTOCOL == protocol::TCP) {
				if (tcp_client_open(sockfd) < 0) {
					throw std::runtime_error("TCP client open error!!");
				}
				send_tcp_th = std::thread(send_tcp_thread, std::ref(sockfd), std::ref(cq_send));
			}

			//----------------------------------------------
			// Connect with Video Sources
			//----------------------------------------------
			cv::VideoCapture cap;
			if (VIDEOSRC_PROTOCOL == videosrc_protocol::RTP) {
				if (cap.open(videosrc, cv::CAP_GSTREAMER)) {
					Logging::set(LogLevel::INFO, "VideoCapture opened.");
				} else {
					throw std::runtime_error("VideoCapture open failed.");
				}
			} else if (VIDEOSRC_PROTOCOL == videosrc_protocol::RTSP) {
				while (true) {
					if (cap.open(videosrc), cv::CAP_GSTREAMER) {
						Logging::set(LogLevel::INFO, "VideoCapture opened.");
						break;
					} else {
						std::this_thread::sleep_for(std::chrono::seconds(1));
					}
				}
			}
			cap_cnt++;

			//----------------------------------------------
			// decode thread start
			//----------------------------------------------
			if (OUTDST_PROTOCOL == protocol::DMA) {
				decode_th = std::thread(
						decode_to_dma_thread,
						std::ref(cap),
						std::ref(mngq),
						std::ref(cq_send_enq),
						std::ref(cq_send_deq),
						std::ref(cq_recv_enq),
						std::ref(cq_recv_deq)
				);
			} else if (OUTDST_PROTOCOL == protocol::TCP) {
				decode_th = std::thread(
						decode_to_tcp_thread,
						std::ref(cap),
						std::ref(cq_send)
				);
			}

			//----------------------------------------------
			// Waiting for thread termination
			//----------------------------------------------
			decode_th.join();

			if (OUTDST_PROTOCOL == protocol::DMA) {
				send_fpga_enq_th.join();
#if defined(EXEC_FPGA_ENQUEUE_MODE)
				send_fpga_deq_th.join();
#endif
#if defined(EXEC_FPGA_DEQUEUE_MODE)
				recv_fpga_enq_th.join();
				recv_fpga_deq_th.join();
				recv_data_th.join();
#endif
			} else if (OUTDST_PROTOCOL == protocol::TCP) {
				send_tcp_th.join();
				if (tcp_client_close(sockfd) < 0) {
					throw std::runtime_error("TCP client close error!!");
				}
				std::this_thread::sleep_for(std::chrono::seconds(1));
			}

#if defined(DEBUG_PRINT)
			std::cout << "all thread finish" << std::endl;
#endif //DEBUG_PRINT

			//----------------------------------------------
			// Video Input Closed
			//----------------------------------------------
			cap.release();

			//----------------------------------------------
			// Video connection limit
			// If VIDEO_CONNECT_LIMIT is 0: Unlimited
			// If VIDEO_CONNECT_LIMIT is greater than 0: Connect the specified number of times before terminating.
			//----------------------------------------------
			if (VIDEO_CONNECT_LIMIT != 0) {
				if (cap_cnt >= VIDEO_CONNECT_LIMIT)
					break;
			}

			// end notification initialization
			is_decode_thread_fin.store(false, std::memory_order_release);
			is_recv_fpga_deq_thread_fin.store(false, std::memory_order_release);
		}


		//--------------------------------------------------------------------
		// termination processing
		//--------------------------------------------------------------------
_END4:
		if (OUTDST_PROTOCOL == protocol::DMA) {
			//----------------------------------------------
			// FPGA kernel finish
			//----------------------------------------------
#if defined(CONTROL_FPGA_FUNC_INIT) && defined(CONTROL_FPGA_FUNC_FINISH)
			tp_function_finish(dev_id);
#endif
		}

_END3:
		if (OUTDST_PROTOCOL == protocol::DMA) {
			//----------------------------------------------
			// lldma queue finish
			//----------------------------------------------
			tp_enqueue_lldma_queue_finish();
#if defined(EXEC_FPGA_DEQUEUE_MODE)
			tp_dequeue_lldma_queue_finish();
#endif
		}

_END2:
		if (OUTDST_PROTOCOL == protocol::DMA) {
			//----------------------------------------------
			// lldma finish
			//----------------------------------------------
#if defined(CONTROL_FPGA_ENQUEUE_LLDMA_INIT)
			tp_enqueue_lldma_finish();
#endif
#if defined(EXEC_FPGA_DEQUEUE_MODE)
			tp_dequeue_lldma_finish();
#endif
		}

_END1:
		if (OUTDST_PROTOCOL == protocol::DMA) {
			//----------------------------------------------
			// shared memory release
			//----------------------------------------------
#if defined(ALLOCATE_SRC_SHMEM_MODE)
			if (mngq.srcbufp != nullptr) {
				Logging::set(LogLevel::INFO, "--- fpga_shmem_free for srcbuf ---");
				fpga_shmem_free(mngq.srcbufp);
			}
#endif

#if defined(EXEC_FPGA_DEQUEUE_MODE)
			if (mngq.dstbufp != nullptr) {
				Logging::set(LogLevel::INFO, "--- fpga_shmem_free for dstbuf ---");
				fpga_shmem_free(mngq.dstbufp);
			}
#endif

			//----------------------------------------------
			// fpga disable regrw
			//----------------------------------------------
			ret = fpga_disable_regrw(dev_id);
			if (ret < 0) {
				Logging::set(LogLevel::ERROR, "fpga_disable_regrw error!!: ret(" + std::to_string(ret) + ")");
			}
			Logging::set(LogLevel::DEBUG, "fpga_disable_regrw: ret(" + std::to_string(ret) + ")");

			//----------------------------------------------
			// finish FPGA
			//----------------------------------------------
#if defined(CONTROL_FPGA_DEV_INIT) && defined(CONTROL_FPGA_FINISH)
			ret = fpga_finish();
			if (ret < 0) {
				Logging::set(LogLevel::ERROR, "fpga_finish error!!: ret(" + std::to_string(ret) + ")");
			}
			Logging::set(LogLevel::DEBUG, "fpga_finish: ret(" + std::to_string(ret) + ")");
#endif

			//----------------------------------------------
			// finish DPDK shmem 
			//----------------------------------------------
#if defined(CONTROL_FPGA_SHMEM_FINISH)
			ret = fpga_shmem_finish();
			if (ret < 0) {
				Logging::set(LogLevel::ERROR, "fpga_shmem_finish error!!: ret("+ std::to_string(ret) + ")");
			}
			Logging::set(LogLevel::DEBUG, "fpga_shmem_finish: ret(" + std::to_string(ret) + ")");
#endif
		}

		Logging::set(LogLevel::INFO, "--- main finish ---");

	} catch (const std::exception& error) {
		std::cerr << "[error] " << error.what() << std::endl;
		Logging::set(LogLevel::ERROR, error.what());
		return EXIT_FAILURE;
	} catch (...) {
		std::cerr << "[error] Unknown/internal exception." << std::endl;
		Logging::set(LogLevel::ERROR, "Unknown/internal exception.");
		return EXIT_FAILURE;
	}

	std::cout << "//--- finish ---//" << std::endl;
	Logging::set(LogLevel::PRINT, "//--- finish ---//");
	Logging::close();

	return EXIT_SUCCESS;
}
