/*****************************************************************
 * Copyright 2024 NTT Corporation, FUJITSU LIMITED
 *****************************************************************/
#include <iostream>
#include <fstream>
#include <sstream>
#include <string>
#include <chrono>
#include <ctime>
#include <iomanip>
#include "cpu_decode.hpp"
#include "recv_header_log.hpp"

//#define DEBUG_PRINT_FPGA_OUT_HEADER

static std::ofstream ofs;

namespace RecvHeaderLog {
	// Log start time prefix
	static std::ostringstream time_prefix() {
		const std::chrono::system_clock::time_point now = std::chrono::system_clock::now();
		const std::time_t t = std::chrono::system_clock::to_time_t(now);
		std::tm *lt = std::localtime(&t);
		std::ostringstream oss;
		oss << lt->tm_year + 1900;
		oss << std::setfill('0') << std::right << std::setw(2) << lt->tm_mon + 1;
		oss << std::setfill('0') << std::right << std::setw(2) << lt->tm_mday;
		oss << "-";
		oss << std::setfill('0') << std::right << std::setw(2) << lt->tm_hour;
		oss << std::setfill('0') << std::right << std::setw(2) << lt->tm_min;
		oss << std::setfill('0') << std::right << std::setw(2) << lt->tm_sec;
		return oss;
	}

	// Log File Start
	static int32_t open() {
		std::ostringstream tp = time_prefix();
		std::string logfile = "recv_header_" + tp.str() + ".log";
		ofs.open(logfile);
		if (!ofs) {
			return -1;
		}
		ofs << "recv header log start..." << tp.str() << std::endl;

		return 0;
	}

	// Display comparison results
	static const char* comp2rslt(int32_t r) {
		return r == 0 ? "OK" : "NG" ;
	}

	// Receive Header Set
	int32_t set(const FrameHeader_t &header, size_t task) {
		if (!ofs.is_open()) {
			if (open() < 0) {
				std::cerr << "RecvHeaderLog.set(): logfile open failed." << std::endl;
				return -1;
			}
		}

		//-----------------------------------------------------
		// expected header value
		//-----------------------------------------------------
		FrameHeader_t exp_fh;
#ifdef MODULE_FPGA
		exp_fh.marker      = 0xE0FF10AD;
		exp_fh.payload_len = FPGA_OUT_FRAME_WIDTH * FPGA_OUT_FRAME_HEIGHT * FRAME_CHANNEL;
		exp_fh.data_id     = FPGA_CH_ID;
#else
		exp_fh.marker       = 0xE0FF10AD;
		exp_fh.payload_len  = FPGA_OUT_FRAME_WIDTH * FPGA_OUT_FRAME_HEIGHT * FRAME_CHANNEL;
		exp_fh.payload_type = 0x01;
		exp_fh.reserved1    = 0x00;
		exp_fh.channel_id   = FPGA_CH_ID % 8;
		//exp_fh.color_space  = 0x01;
		exp_fh.color_space  = 0x00;
		exp_fh.data_type    = 0x00;
		//exp_fh.num_ch       = 0x03;
		exp_fh.num_ch       = 0x00;
		exp_fh.width        = FPGA_OUT_FRAME_WIDTH;
		exp_fh.height       = FPGA_OUT_FRAME_HEIGHT;
		exp_fh.local_ts     = 0.0;
#endif //MODULE_FPGA

		//-----------------------------------------------------
		// header validation
		//-----------------------------------------------------
		header_verify_info_t header_result;
#ifdef MODULE_FPGA
		header_result.result_marker = header.marker == exp_fh.marker ? 0 : 1;
		header_result.result_payload_len = header.payload_len == exp_fh.payload_len ? 0 : 1;
		header_result.result_data_id = header.data_id == exp_fh.data_id ? 0 : 1;
#else
		header_result.result_marker = header.marker == exp_fh.marker ? 0 : 1;
		header_result.result_payload_len = header.payload_len == exp_fh.payload_len ? 0 : 1;
		header_result.result_payload_type = header.payload_type == exp_fh.payload_type ? 0 : 1;
		header_result.result_channel_id = header.channel_id == exp_fh.channel_id ? 0 : 1;
		header_result.result_color_space = header.color_space == exp_fh.color_space ? 0 : 1;
		header_result.result_data_type = header.data_type == exp_fh.data_type ? 0 : 1;
		header_result.result_num_ch = header.num_ch == exp_fh.num_ch ? 0 : 1;
		header_result.result_width = header.width == exp_fh.width ? 0 : 1;
		header_result.result_height = header.height == exp_fh.height ? 0 : 1;
#endif //MODULE_FPGA

		//-----------------------------------------------------
		// Output of header verification results
		//-----------------------------------------------------
		std::ostringstream oss;
#ifdef MODULE_FPGA
		oss << "----------------------------------------------------------------" << std::endl;
#else
		oss << "--------------------------------------------------------" << std::endl;
#endif
		oss << "FrameHeader CH(" << FPGA_CH_ID << ") TASK(" << task << ")" << std::endl;
#ifdef MODULE_FPGA
		oss << "                        Receive value | Expected value | compare" << std::endl;
#else
		oss << "                Receive value | Expected value | compare" << std::endl;
#endif //MODULE_FPGA

		oss.setf(std::ios::hex, std::ios::basefield);
		char originalfill = oss.fill('0');
#ifdef MODULE_FPGA
		oss << "  marker         :         0x" << std::setw(8) << header.marker
			<< " |     0x" << std::setw(8) << exp_fh.marker
			<< " |   " << comp2rslt(header_result.result_marker) << std::endl;
		oss << "  payload_len    :         0x" << std::setw(8) << header.payload_len
			<< " |     0x" << std::setw(8) << exp_fh.payload_len
			<< " |   " << comp2rslt(header_result.result_payload_len) << std::endl;
		oss << "  sequence_num   :         0x" << std::setw(8) << header.sequence_num
			<< " |              -"
			<< " |   -" << std::endl;
		oss << "  timestamp      : 0x" << std::setw(16) << static_cast<int64_t>(header.timestamp)
			<< " |              -"
			<< " |   -" << std::endl;
		oss << "  data_id        :         0x" << std::setw(8) << header.data_id
			<< " |     0x" << std::setw(8) << exp_fh.data_id
			<< " |   " << comp2rslt(header_result.result_data_id) << std::endl;
		oss << "  header_checksum:             0x" << std::setw(4) << header.header_checksum
			<< " |              -"
			<< " |   -" << std::endl;
#else
		oss << "  marker      :    0x" << std::setw(8) << header.marker
			<< " |     0x" << std::setw(8) << exp_fh.marker
			<< " |   " << comp2rslt(header_result.result_marker) << std::endl;
		oss << "  payload_len :    0x" << std::setw(8) << header.payload_len
			<< " |     0x" << std::setw(8) << exp_fh.payload_len
			<< " |   " << comp2rslt(header_result.result_payload_len) << std::endl;
		oss << "  payload_type:          0x" << std::setw(2) << static_cast<uint16_t>(header.payload_type)
			<< " |           0x" << std::setw(2) << static_cast<uint16_t>(exp_fh.payload_type)
			<< " |   " << comp2rslt(header_result.result_payload_type) << std::endl;
		oss << "  channel_id  :        0x" << std::setw(4) << header.channel_id
			<< " |         0x" << std::setw(4) << exp_fh.channel_id
			<< " |   " << comp2rslt(header_result.result_channel_id) << std::endl;
		oss << "  frame_index :    0x" << std::setw(8) << header.frame_index
			<< " |              -"
			<< " |   -" << std::endl;
		oss << "  color_space :          0x" << std::setw(2) << static_cast<uint16_t>(header.color_space)
			<< " |           0x" << std::setw(2) << static_cast<uint16_t>(exp_fh.color_space)
			<< " |   " << comp2rslt(header_result.result_color_space) << std::endl;
		oss << "  data_type   :          0x" << std::setw(2) << static_cast<uint16_t>(header.data_type)
			<< " |           0x" << std::setw(2) << static_cast<uint16_t>(exp_fh.data_type)
			<< " |   " << comp2rslt(header_result.result_data_type) << std::endl;
		oss << "  num_ch      :        0x" << std::setw(4) << header.num_ch
			<< " |         0x" << std::setw(4) << exp_fh.num_ch
			<< " |   " << comp2rslt(header_result.result_num_ch) << std::endl;
		oss << "  width       :        0x" << std::setw(4) << header.width
			<< " |         0x" << std::setw(4) << exp_fh.width
			<< " |   " << comp2rslt(header_result.result_width) << std::endl;
		oss << "  height      :        0x" << std::setw(4) << header.height
			<< " |         0x" << std::setw(4) << exp_fh.height
			<< " |   " << comp2rslt(header_result.result_height) << std::endl;
#endif //MODULE_FPGA
		oss.setf(std::ios::dec, std::ios::basefield);
		oss.fill(originalfill);

		ofs << oss.str() << std::endl;
#if defined(DEBUG_PRINT_FPGA_OUT_HEADER)
		std::cout << oss.str() << std::endl;
#endif //DEBUG_PRINT_FPGA_OUT_HEADER
		oss.str("");

		return 0;
	}

	// End Log File
	void close() {
		if (ofs.is_open()) {
			ofs.close();
		}
	}
}
