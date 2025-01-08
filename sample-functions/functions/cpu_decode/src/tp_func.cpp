/*****************************************************************
 * Copyright 2024 NTT Corporation, FUJITSU LIMITED
 *****************************************************************/
#include <iostream>
#include <sstream>
#include <chrono>
#include <thread>
#include "cpu_decode.hpp"
#include "libshmem.h"
#include "libfpgactl.h"
#include "libchain.h"
#include "liblldma.h"
#include "libdmacommon.h"
#include "libdma.h"
#include "liblogging.h"
#include "libfunction.h"
#include "libdirecttrans.h"
#include "libfunction_filter_resize.h"
#include "libfunction_conv.h"
#include "libchain.h"

#define JSON_FORMAT "{ \
  \"i_width\"   :%u, \
  \"i_height\"  :%u, \
  \"o_width\"   :%u, \
  \"o_height\"  :%u \
}"

int32_t tp_function_filter_resize_init(uint32_t dev_id)
{
	int32_t ret = 0;

	uint32_t lane_id = 0;
	if (FPGA_CH_ID > (CH_NUM_MAX / LANE_NUM_MAX - 1)) {
		lane_id = 1;
	}
	uint32_t krnl_id = lane_id * FUNC_IDX;

	std::ostringstream debugstr;

#ifdef MODULE_FPGA
	bool module_id_err = false;
	uint32_t module_id = 0;

	Logging::set(LogLevel::INFO, "--- module_id check ---");

	fpga_chain_get_module_id(dev_id, krnl_id, &module_id);
	debugstr << "dev(" << dev_id << ") kernel_id(" << krnl_id << ") chain_control module_id(0x" << std::hex << module_id << ")" ;
	if (module_id == MODULE_ID_FPGA_CHAIN) {
		Logging::set(LogLevel::DEBUG, debugstr.str());
	} else {	
		debugstr << " error! expected module_id(0x" << MODULE_ID_FPGA_CHAIN << ")";
		Logging::set(LogLevel::ERROR, debugstr.str());
		module_id_err = true;
	}
	debugstr.str("");

	module_id = 0;
	fpga_direct_get_module_id(dev_id, krnl_id, &module_id);
	debugstr << "dev(" << dev_id << ") kernel_id(" << krnl_id << ") direct_trans_adapter module_id(0x" << std::hex << module_id << ")" ;
	if (module_id == MODULE_ID_FPGA_DIRECT) {
		Logging::set(LogLevel::DEBUG, debugstr.str());
	} else {	
		debugstr << " error! expected module_id(0x" << MODULE_ID_FPGA_DIRECT << ")";
		Logging::set(LogLevel::ERROR, debugstr.str());
		module_id_err = true;
	}
	debugstr.str("");

	module_id = 0;
	fpga_filter_resize_get_module_id(dev_id, krnl_id, &module_id);
	debugstr << "dev(" << dev_id << ") kernel_id(" << krnl_id << ") function module_id(0x" << std::hex << module_id << ")" ;
	if (module_id == MODULE_ID_FPGA_FILTER_RESIZE) {
		Logging::set(LogLevel::DEBUG, debugstr.str());
	} else {	
		debugstr << " error! expected module_id(0x" << MODULE_ID_FPGA_FILTER_RESIZE << ")";
		Logging::set(LogLevel::ERROR, debugstr.str());
		module_id_err = true;
	}
	debugstr.str("");

	module_id = 0;
	fpga_conv_get_module_id(dev_id, krnl_id, &module_id);
	debugstr << "dev(" << dev_id << ") kernel_id(" << krnl_id << ") conversion_adapter module_id(0x" << std::hex << module_id << ")" ;
	if (module_id == MODULE_ID_FPGA_CONV) {
		Logging::set(LogLevel::DEBUG, debugstr.str());
	} else {	
		debugstr << " error! expected module_id(0x" << MODULE_ID_FPGA_CONV << ")";
		Logging::set(LogLevel::ERROR, debugstr.str());
		module_id_err = true;
	}
	debugstr.str("");

	if (module_id_err) {
		//error
		return -1;
	}
#endif //MODULE_FPGA

	Logging::set(LogLevel::INFO, "--- fpga_function filter_resize init ---");

	char json_txt[256];
	snprintf(json_txt, 256, JSON_FORMAT, FRAME_WIDTH, FRAME_HEIGHT, FPGA_OUT_FRAME_WIDTH, FPGA_OUT_FRAME_HEIGHT);

	debugstr << "dev(" << dev_id << ") func_kernel(" << krnl_id << ") json_txt: " << std::string(json_txt);
	Logging::set(LogLevel::DEBUG, debugstr.str());
	debugstr.str("");

	debugstr << "dev(" << dev_id << ") func_kernel(" << krnl_id << ") fpga_function_config";
	Logging::set(LogLevel::DEBUG, debugstr.str());
	debugstr.str("");
	ret = fpga_function_config(dev_id, krnl_id, "filter_resize");
	if (ret < 0) {
		//error
		Logging::set(LogLevel::ERROR, "fpga_function_config error!!: ret(" + std::to_string(ret) + ")");
		return -1;
	}

	debugstr << "dev(" << dev_id << ") func_kernel(" << krnl_id << ") fpga_function_init";
	Logging::set(LogLevel::DEBUG, debugstr.str());
	debugstr.str("");
	ret = fpga_function_init(dev_id, krnl_id, NULL);
	if (ret < 0) {
		//error
		Logging::set(LogLevel::ERROR, "fpga_function_init error!!: ret(" + std::to_string(ret) + ")");
		return -1;
	}

	debugstr << "dev(" << dev_id << ") func_kernel(" << krnl_id << ") fpga_function_set";
	Logging::set(LogLevel::DEBUG, debugstr.str());
	debugstr.str("");
	ret = fpga_function_set(dev_id, krnl_id, json_txt);
	if (ret < 0) {
		//error
		Logging::set(LogLevel::ERROR, "fpga_function_set error!!: ret(" + std::to_string(ret) + ")");
		return -2;
	}

#ifdef MODULE_FPGA
	// Chain Control DDR Offset Settings
	for (size_t i=0; i<EXTIF_NUM_MAX; i++) {
		uint32_t extif_id = i;
		debugstr << "dev(" << dev_id << ") kernel(" << krnl_id << ") extif(" << extif_id << ") fpga_chain_set_ddr";
		Logging::set(LogLevel::DEBUG, debugstr.str());
		debugstr.str("");
		ret = fpga_chain_set_ddr(dev_id, krnl_id, extif_id);
		if (ret < 0) {
			//error
			Logging::set(LogLevel::ERROR, "fpga_chain_set_ddr error!!: ret(" + std::to_string(ret) + ")");
			return -2;
		}

		// Verifying the DDR Configuration
		fpga_chain_ddr_t chain_ddr;
		ret = fpga_chain_get_ddr(dev_id, krnl_id, extif_id, &chain_ddr);
		if (ret < 0) {
			//error
			Logging::set(LogLevel::ERROR, "fpga_chain_get_ddr error!!: ret(" + std::to_string(ret) + ")");
		} else {
			prlog_fpga_chain_ddr_info(chain_ddr, dev_id, krnl_id, extif_id);
		}
	}

	// direct transfer adapter kernel module boot
	debugstr << "dev(" << dev_id << ") kernel(" << krnl_id << ") fpga_direct_start";
	Logging::set(LogLevel::DEBUG, debugstr.str());
	debugstr.str("");
	ret = fpga_direct_start(dev_id, krnl_id);
	if (ret < 0) {
		//error
		Logging::set(LogLevel::ERROR, "fpga_direct_start error!!: ret(" + std::to_string(ret) + ")");
		return -2;
	}

	// chain control kernel module startup
	debugstr << "dev(" << dev_id << ") kernel(" << krnl_id << ") fpga_chain_start";
	Logging::set(LogLevel::DEBUG, debugstr.str());
	debugstr.str("");
	ret = fpga_chain_start(dev_id, krnl_id);
	if (ret < 0) {
		//error
		Logging::set(LogLevel::ERROR, "fpga_chain_start error!!: ret(" + std::to_string(ret) + ")");
		return -2;
	}
#endif //MODULE_FPGA

	return 0;
}

void tp_function_finish(uint32_t dev_id)
{
	int32_t ret = 0;

	Logging::set(LogLevel::INFO, "--- fpga_function_finish ---");

	uint32_t lane_id = 0;
	if (FPGA_CH_ID > (CH_NUM_MAX / LANE_NUM_MAX - 1)) {
		lane_id = 1;
	}
	uint32_t krnl_id = lane_id * FUNC_IDX;

	std::ostringstream debugstr;

#ifdef MODULE_FPGA
	// direct transfer adapter kernel module stop
	debugstr << "dev(" << dev_id << ") kernel(" << krnl_id << ") fpga_direct_stop";
	Logging::set(LogLevel::DEBUG, debugstr.str());
	debugstr.str("");
	ret = fpga_direct_stop(dev_id, krnl_id);
	if (ret < 0) {
		//error
		Logging::set(LogLevel::ERROR, "fpga_direct_stop error!!: ret(" + std::to_string(ret) + ")");
	}

	// chain control kernel module stop
	debugstr << "dev(" << dev_id << ") kernel(" << krnl_id << ") fpga_chain_stop";
	Logging::set(LogLevel::DEBUG, debugstr.str());
	debugstr.str("");
	ret = fpga_chain_stop(dev_id, krnl_id);
	if (ret < 0) {
		//error
		Logging::set(LogLevel::ERROR, "fpga_chain_stop error!!: ret(" + std::to_string(ret) + ")");
	}
#endif //MODULE_FPGA

	debugstr << "dev(" << dev_id << ") func_kernel(" << krnl_id << ") fpga_function_finish";
	Logging::set(LogLevel::DEBUG, debugstr.str());
	debugstr.str("");
	ret = fpga_function_finish(dev_id, krnl_id, NULL);
	if (ret < 0) {
		//error
		Logging::set(LogLevel::ERROR, "fpga_function_finish error!!: ret(" + std::to_string(ret) + ")");
	}
}

int32_t tp_enqueue_lldma_init(uint32_t dev_id)
{
	int32_t ret = 0;

	Logging::set(LogLevel::INFO, "--- enqueue fpga_lldma_init  ---");

	memset(&enqdmainfo_channel, 0, sizeof(dma_info_t));
	const char *connector_id = FPGA_DMA_HOST_TO_DEV_CONNECTOR_ID;
	std::ostringstream debugstr;
	debugstr << "dev(" << dev_id << ") CH(" << FPGA_CH_ID << ") enqueue fpga_lldma_init";
	Logging::set(LogLevel::DEBUG, debugstr.str());
	debugstr.str("");
	ret = fpga_lldma_init(dev_id, DMA_HOST_TO_DEV, FPGA_CH_ID, const_cast<char*>(connector_id), &enqdmainfo_channel);
	if (ret < 0) {
		// error
		Logging::set(LogLevel::ERROR, "enqueue fpga_lldma_init error!!: ret(" + std::to_string(ret) + ")");
		return -1;
	}
	prlog_dma_info(enqdmainfo_channel);

        return 0;
}

int32_t tp_dequeue_lldma_init(uint32_t dev_id)
{
	int32_t ret = 0;

	Logging::set(LogLevel::INFO, "--- dequeue fpga_lldma_init  ---");

	memset(&deqdmainfo_channel, 0, sizeof(dma_info_t));
	std::string cs = "deq_" + std::string(FPGA_DMA_HOST_TO_DEV_CONNECTOR_ID);
	const char *connector_id = cs.c_str();
	std::ostringstream debugstr;
	debugstr << "dev(" << dev_id << ") CH(" << FPGA_CH_ID << ") dequeue fpga_lldma_init";
	Logging::set(LogLevel::DEBUG, debugstr.str());
	debugstr.str("");
	ret = fpga_lldma_init(dev_id, DMA_DEV_TO_HOST, FPGA_CH_ID, const_cast<char*>(connector_id), &deqdmainfo_channel);
	if (ret < 0) {
		// error
		Logging::set(LogLevel::ERROR, "dequeue fpga_lldma_init error!!: ret(" + std::to_string(ret) + ")");
		return -1;
	}
	prlog_dma_info(deqdmainfo_channel);

        return 0;
}

void tp_enqueue_lldma_finish(void)
{
	int32_t ret = 0;

	Logging::set(LogLevel::INFO, "--- enqueue fpga_lldma_finish ---");

	ret = fpga_lldma_finish(&enqdmainfo_channel);
	if (ret < 0) {
		// error
		Logging::set(LogLevel::ERROR, "enqueue fpga_lldma_finish error!!: ret(" + std::to_string(ret) + ")");
	}
	enqdmainfo_channel.connector_id = nullptr;
	prlog_dma_info(enqdmainfo_channel);
}

void tp_dequeue_lldma_finish(void)
{
	int32_t ret = 0;

	Logging::set(LogLevel::INFO, "--- dequeue fpga_lldma_finish ---");

	ret = fpga_lldma_finish(&deqdmainfo_channel);
	if (ret < 0) {
		// error
		Logging::set(LogLevel::ERROR, "dequeue fpga_lldma_finish error!!: ret(" + std::to_string(ret) + ")");
	}
	deqdmainfo_channel.connector_id = nullptr;
	prlog_dma_info(deqdmainfo_channel);
}

int32_t tp_enqueue_lldma_queue_setup(void)
{
	int32_t ret = 0;

	Logging::set(LogLevel::INFO, "--- enqueue fpga_lldma_queue_setup ---");

	memset(&enqdmainfo, 0, sizeof(dma_info_t));
	const char *connector_id = enqdmainfo_channel.connector_id;
	Logging::set(LogLevel::DEBUG, "CH(" + std::to_string(FPGA_CH_ID) + ") enqueue fpga_lldma_queue_setup");
	ret = fpga_lldma_queue_setup(const_cast<char*>(connector_id), &enqdmainfo);
	if (ret < 0) {
		// error
		Logging::set(LogLevel::ERROR, "enqueue fpga_lldma_queue_setup error!!: ret(" + std::to_string(ret) + ")");
		return -1;
	}
	prlog_dma_info(enqdmainfo);

	return 0;
}

int32_t tp_dequeue_lldma_queue_setup(void)
{
	int32_t ret = 0;

	Logging::set(LogLevel::INFO, "--- dequeue fpga_lldma_queue_setup ---");

	memset(&deqdmainfo, 0, sizeof(dma_info_t));
	const char *connector_id = deqdmainfo_channel.connector_id;
	Logging::set(LogLevel::DEBUG, "CH(" + std::to_string(FPGA_CH_ID) + ") dequeue fpga_lldma_queue_setup");
	ret = fpga_lldma_queue_setup(const_cast<char*>(connector_id), &deqdmainfo);
	if (ret < 0) {
		// error
		Logging::set(LogLevel::ERROR, "dequeue fpga_lldma_queue_setup error!!: ret(" + std::to_string(ret) + ")");
		return -1;
	}
	prlog_dma_info(deqdmainfo);

	return 0;
}

void tp_enqueue_lldma_queue_finish(void)
{
	int32_t ret = 0;

	Logging::set(LogLevel::INFO, "--- enqueue fpga_lldma_queue_finish ---");

	Logging::set(LogLevel::DEBUG, "CH(" + std::to_string(FPGA_CH_ID) + ") enqueue fpga_lldma_queue_finish");
	ret = fpga_lldma_queue_finish(&enqdmainfo);
	if (ret < 0) {
		Logging::set(LogLevel::ERROR, "enqueue fpga_lldma_queue_finish error!!: ret(" + std::to_string(ret) + ")");
	}
	enqdmainfo.connector_id = nullptr;
	prlog_dma_info(enqdmainfo);
}

void tp_dequeue_lldma_queue_finish(void)
{
	int32_t ret = 0;

	Logging::set(LogLevel::INFO, "--- dequeue fpga_lldma_queue_finish ---");

	Logging::set(LogLevel::DEBUG, "CH(" + std::to_string(FPGA_CH_ID) + ") dequeue fpga_lldma_queue_finish");
	ret = fpga_lldma_queue_finish(&deqdmainfo);
	if (ret < 0) {
		Logging::set(LogLevel::ERROR, "dequeue fpga_lldma_queue_finish error!!: ret(" + std::to_string(ret) + ")");
	}
	deqdmainfo.connector_id = nullptr;
	prlog_dma_info(deqdmainfo);
}

int32_t tp_chain_connect(uint32_t dev_id)
{
	int32_t ret = 0;

	Logging::set(LogLevel::INFO, "--- fpga_chain_connect ---");

	uint32_t lane_id = 0;
	if (FPGA_CH_ID > (CH_NUM_MAX / LANE_NUM_MAX - 1)) {
		lane_id = 1;
	}
	uint32_t func_krnl_pid = lane_id * FUNC_IDX;
	uint32_t fchid = g_param_chain_ctrl_tbls[FPGA_CH_ID].function_chid;
	uint32_t ingress_cid = g_param_chain_ctrl_tbls[FPGA_CH_ID].lldma_cid;
	uint32_t egress_cid = g_param_chain_ctrl_tbls[FPGA_CH_ID].lldma_cid;
	uint32_t ingress_extif_id = 0;
	uint32_t egress_extif_id = 0;
	uint8_t direct_flag = 0;
	uint8_t ig_active_flag = 1;
	uint8_t eg_active_flag = 1;
	uint8_t virtual_flag = 0;
	uint8_t blocking_flag = 1;
	std::ostringstream debugstr;

#ifdef MODULE_FPGA
	// ingress Confirm Connection Establishment
	uint32_t con_status = 0;
	bool con_err = false;
	ret = fpga_chain_get_con_status(dev_id, func_krnl_pid, ingress_extif_id, ingress_cid, &con_status);
	if (ret < 0) {
		debugstr << "dev(" << dev_id << ") CH(" << FPGA_CH_ID << ") func_kernel_id(" << func_krnl_pid  << ") ingress fpga_chain_get_con_status error!!: ret(" << ret << ")";
		Logging::set(LogLevel::ERROR, debugstr.str());	
		debugstr.str("");
		con_err = true;
	}
	if (con_status == 0) {
		debugstr << "dev(" << dev_id << ") CH(" << FPGA_CH_ID << ") func_kernel_id(" << func_krnl_pid  << ") fpga_chain_get_con_status() chain connection error. ";
		debugstr << "ingress_extif_id(" << ingress_extif_id << ") ingress_cid(" << ingress_cid << ") status(0x" << con_status << ")";
		Logging::set(LogLevel::ERROR, debugstr.str());
		debugstr.str("");
		con_err = true;
	} else {
		debugstr << "dev(" << dev_id << ") CH(" << FPGA_CH_ID << ") func_kernel_id(" << func_krnl_pid  << ") fpga_chain_get_con_status() chain connection established. ";
		debugstr << "ingress_extif_id(" << ingress_extif_id << ") ingress_cid(" << ingress_cid << ") status(0x" << con_status << ")";
		Logging::set(LogLevel::DEBUG, debugstr.str());
		debugstr.str("");
	}
	// Verify egress connection establishment
	con_status = 0;
	ret = fpga_chain_get_con_status(dev_id, func_krnl_pid, egress_extif_id, egress_cid, &con_status);
	if (ret < 0) {
		debugstr << "dev(" << dev_id << ") CH(" << FPGA_CH_ID << ") func_kernel_id(" << func_krnl_pid  << ") egress fpga_chain_get_con_status error!!: ret(" << ret << ")";
		Logging::set(LogLevel::ERROR, debugstr.str());	
		debugstr.str("");
		con_err = true;
	}
	if (con_status == 0) {
		debugstr << "dev(" << dev_id << ") CH(" << FPGA_CH_ID << ") func_kernel_id(" << func_krnl_pid  << ") fpga_chain_get_con_status() chain connection error. ";
		debugstr << "egress_extif_id(" << egress_extif_id << ") egress_cid(" << egress_cid << ") status(0x" << con_status << ")";
		Logging::set(LogLevel::ERROR, debugstr.str());
		debugstr.str("");
		con_err = true;
	} else {
		debugstr << "dev(" << dev_id << ") CH(" << FPGA_CH_ID << ") func_kernel_id(" << func_krnl_pid  << ") fpga_chain_get_con_status() chain connection established. ";
		debugstr << "egress_extif_id(" << egress_extif_id << ") egress_cid(" << egress_cid << ") status(0x" << con_status << ")";
		Logging::set(LogLevel::DEBUG, debugstr.str());
		debugstr.str("");
	}
	if (con_err) {
		//error
		return -1;
	}
#endif //MODULE_FPGA

	debugstr << "dev(" << dev_id << ") CH(" << FPGA_CH_ID << ") fpga_chain_connect";
	Logging::set(LogLevel::DEBUG, debugstr.str());
	debugstr.str("");
	debugstr << "  func_kernel_id(" << func_krnl_pid << ") fchid(" << fchid << ") ingress_cid(" << ingress_cid << ") egress_cid(" << egress_cid << ")";
	Logging::set(LogLevel::DEBUG, debugstr.str());
	debugstr.str("");
#ifdef MODULE_FPGA
	debugstr << "  ingress_extif_id(" << ingress_extif_id << ") egress_extif_id(" << egress_extif_id << ") ";
	debugstr << "ingress_active_flag(" << std::to_string(ig_active_flag) << ") egress_active_flag(" << std::to_string(eg_active_flag) << ") ";
	debugstr << "direct_flag(" << std::to_string(direct_flag) << ") ";
	debugstr << "virtual_flag(" << std::to_string(virtual_flag) << ") ";
	debugstr << "blocking_flag(" << std::to_string(blocking_flag) << ")";
	Logging::set(LogLevel::DEBUG, debugstr.str());
	debugstr.str("");
#endif //MODULE_FPGA
	ret = fpga_chain_connect(
			dev_id,
			func_krnl_pid,
			fchid,
			ingress_extif_id,
			ingress_cid,
			egress_extif_id,
			egress_cid,
			ig_active_flag,
			eg_active_flag,
			direct_flag,
			virtual_flag,
			blocking_flag
	);
	if (ret < 0) {
		// error
		Logging::set(LogLevel::ERROR, "fpga_chain_connect error!!: ret(" + std::to_string(ret) + ")");
		return -1;
	}

	return 0;
}

int32_t tp_enqueue_set_dma_cmd(uint32_t enq_id, const mngque_t &mngq, dmacmd_info_t &enqdmacmdinfo)
{
	int32_t ret = 0;

	Logging::set(LogLevel::DEBUG, "--- enqueue set_dma_cmd ---");

	static uint16_t taskidx = 1;
	uint16_t task_id = taskidx;
	uint32_t data_len = mngq.srcbuflen;
	uint32_t dsize = mngq.srcdsize;
	void *data_addr = mngq.srcbufp;
	Logging::set(LogLevel::DEBUG, "CH(" + std::to_string(FPGA_CH_ID) + ") dma rx data size=" + std::to_string(dsize) + " Byte");
	memset(&enqdmacmdinfo, 0, sizeof(dmacmd_info_t));
	Logging::set(LogLevel::DEBUG, "CH(" + std::to_string(FPGA_CH_ID) + ") ENQ(" + std::to_string(enq_id) + ") set_dma_cmd");
	ret = set_dma_cmd(&enqdmacmdinfo, task_id, data_addr, data_len);
	if (ret < 0) {
		// error
		Logging::set(LogLevel::ERROR, "enqueue set_dma_cmd error!!: ret(" + std::to_string(ret) + ")");
		return -1;
	}
	prlog_dmacmd_info(enqdmacmdinfo, enq_id);

	if (taskidx == 0xFFFF) {
		taskidx = 1;
	} else {
		taskidx++;
	}

	return 0;
}

int32_t tp_dequeue_set_dma_cmd(uint32_t enq_id, const mngque_t &mngq, dmacmd_info_t &deqdmacmdinfo)
{
	int32_t ret = 0;

	Logging::set(LogLevel::DEBUG, "--- dequeue set_dma_cmd ---");

	static uint16_t taskidx = 1;
	uint32_t task_id = taskidx;
	uint32_t data_len = mngq.dstbuflen;
	uint32_t dsize = mngq.dstdsize;
	void *data_addr = mngq.dstbufp;
	Logging::set(LogLevel::DEBUG, "CH(" + std::to_string(FPGA_CH_ID) + ") dma rx data size=" + std::to_string(dsize) + " Byte");
	memset(&deqdmacmdinfo, 0, sizeof(dmacmd_info_t));
	Logging::set(LogLevel::DEBUG, "CH(" + std::to_string(FPGA_CH_ID) + ") ENQ(" + std::to_string(enq_id) + ") set_dma_cmd");
	ret = set_dma_cmd(&deqdmacmdinfo, task_id, data_addr, data_len);
	if (ret < 0) {
		// error
		Logging::set(LogLevel::ERROR, "dequeue set_dma_cmd error!!: ret(" + std::to_string(ret) + ")");
		return -1;
	}
	prlog_dmacmd_info(deqdmacmdinfo, enq_id);

	if (taskidx == 0xFFFF) {
		taskidx = 1;
	} else {
		taskidx++;
	}

	return 0;
}

int32_t wait_dma_fpga_enqueue(dma_info_t &dmainfo, dmacmd_info_t &dmacmdinfo, uint32_t enq_id, const uint32_t msec)
{
	uint32_t wait_msec = 100; //fpga_enqueue timeout 100msec
	int32_t ret = 0;

	uint16_t ch_id = dmainfo.chid;
	uint16_t task_id = dmacmdinfo.task_id;
	std::string dir = "RX";
	if (dmainfo.dir == DMA_DEV_TO_HOST) {
		dir = "TX";
	}

	uint32_t cnt = msec/wait_msec;
	for (size_t i=0; i < cnt; i++) {
		ret = fpga_enqueue(&dmainfo, &dmacmdinfo);
		if (ret == 0) {
			return 0;
		} else if (ret == -ENQUEUE_QUEFULL) {
			Logging::set(LogLevel::DEBUG, "  CH(" + std::to_string(ch_id) + ") enq(" + std::to_string(enq_id) + ") task_id(" + std::to_string(task_id) + ") DMA " + dir + " fpga_enqueue que full(" + std::to_string(ret) + ")");
                        std::this_thread::sleep_for(std::chrono::microseconds(wait_msec * 1000));
		} else {
			Logging::set(LogLevel::ERROR, "  CH(" + std::to_string(ch_id) + ") enq(" + std::to_string(enq_id) + ") task_id(" + std::to_string(task_id) + ") DMA " + dir + " fpga_enqueue error!!!(" + std::to_string(ret) + ")");
			return ret;
		}
	}

	Logging::set(LogLevel::ERROR, "  CH(" + std::to_string(ch_id) + ") enq(" + std::to_string(enq_id) + ") task_id(" + std::to_string(task_id) + ") DMA " + dir + " enqueue timeout!!!");

	return -1;
}

int32_t wait_dma_fpga_dequeue(dma_info_t &dmainfo, dmacmd_info_t &dmacmdinfo, uint32_t enq_id, const uint32_t msec)
{
	uint32_t timeout = 100; //fpga_deuque timeout 100msec

	uint16_t ch_id = dmainfo.chid;
	uint16_t task_id = dmacmdinfo.task_id;
	std::string dir = "RX";
	if (dmainfo.dir == DMA_DEV_TO_HOST) {
		dir = "TX";
	}

	uint32_t cnt = msec/timeout;
	for (size_t i=0; i < cnt; i++) {
		int32_t ret = fpga_dequeue(&dmainfo, &dmacmdinfo);
		if (ret == 0) {
			return 0;
		}
	}

	Logging::set(LogLevel::ERROR, "  CH(" + std::to_string(ch_id) + ") enq(" + std::to_string(enq_id) + ") task_id(" + std::to_string(task_id) + ") DMA " + dir + " dequeue timeout!!!");

	return -1;
}

