/*****************************************************************
 * Copyright 2024 NTT Corporation, FUJITSU LIMITED
 *****************************************************************/
#ifndef CPU_DECODE_HPP
#define CPU_DECODE_HPP

#include <opencv2/opencv.hpp>
#include "logging.hpp"
#include "libdmacommon.h"
#include "libchain.h"

//-----------------------------------------------------
// define
//-----------------------------------------------------
#define VERSION "0.6.00"

#define CH_NUM_MAX      16
#define LANE_NUM_MAX                    2
#define PTU_KRNL_NUM_MAX                2
#define FRAMEWORK_KRNL_NUM_MAX          2
#define FRAMEWORK_SUB_KRNL_NUM_MAX      2
#define FUNCTION_KRNL_NUM_MAX           2
#define EXTIF_NUM_MAX                   2

#ifdef MODULE_FPGA
#define MODULE_ID_FPGA_CHAIN		0x0000F0C0 // Chain control
#define MODULE_ID_FPGA_DIRECT		0x0000F3C0 // Direct transfer
#define MODULE_ID_FPGA_FILTER_RESIZE	0x0000F2C2 // Function (filter/resize)
#define MODULE_ID_FPGA_CONV		0x0000F1C2 // Conversion(filter/resize)
#endif // MODULE_FPGA

#define FW_IDX          FRAMEWORK_KRNL_NUM_MAX / LANE_NUM_MAX
#define FWSUB_IDX       FRAMEWORK_SUB_KRNL_NUM_MAX / LANE_NUM_MAX
#define FUNC_IDX        FUNCTION_KRNL_NUM_MAX / LANE_NUM_MAX

#define ALIGN_BUF_LEN   64

#define WAIT_TIME_DMA_RX_ENQUEUE	2000	//msec
#define WAIT_TIME_DMA_RX_DEQUEUE	2000	//msec
#define WAIT_TIME_DMA_TX_ENQUEUE	2000	//msec
#define WAIT_TIME_DMA_TX_DEQUEUE	2000	//msec

enum class videosrc_protocol : uint8_t {
	RTP = 0,
	RTSP
};

enum class protocol : uint8_t {
	DMA = 0,
	TCP
};

//-----------------------------------------------------
// struct
//-----------------------------------------------------
#ifdef MODULE_FPGA
typedef struct FrameHeader {
	uint32_t marker;
	uint32_t payload_len;
	uint8_t reserved1[4];
	uint32_t sequence_num; // oldframe_index
	uint8_t reserved2[8];
	double timestamp;
	uint32_t data_id; // oldchannel_id
	uint8_t reserved3[8];
	uint16_t header_checksum;
	uint8_t reserved4[2];
} FrameHeader_t;
#else
typedef struct FrameHeader {
	uint32_t marker;
	uint32_t payload_len;
	uint8_t payload_type;
	uint8_t reserved1;
	uint16_t channel_id;
	uint32_t frame_index;
	uint8_t color_space;
	uint8_t data_type;
	uint16_t num_ch;
	uint16_t width;
	uint16_t height;
	double local_ts;
	uint8_t reserved2[16];
} FrameHeader_t;
#endif //MODULE_FPGA

typedef struct mngque {
	uint32_t srcdsize;
	uint32_t dstdsize;
	uint32_t srcbuflen;
	uint32_t dstbuflen;
	void *srcbufp;
	void *dstbufp;
} mngque_t;

typedef struct chain_ctrl_cid {
	uint32_t function_chid;
	uint32_t lldma_cid;
} chain_ctrl_cid_t;

typedef struct send_fpga_info {
	dmacmd_info_t dmacmdinfo;
	cv::Mat mat;
} send_fpga_info_t;

typedef struct recv_fpga_data {
	FrameHeader_t frame_header;
	cv::Mat mat;
} recv_fpga_data_t;

//-----------------------------------------------------
// global variables
//-----------------------------------------------------
extern const char *FPGA_DEV_NAME;
extern uint16_t FPGA_CH_ID;
extern uint32_t FRAME_WIDTH;
extern uint32_t FRAME_HEIGHT;
extern uint32_t FRAME_CHANNEL;
extern uint32_t FPGA_OUT_FRAME_WIDTH;
extern uint32_t FPGA_OUT_FRAME_HEIGHT;
extern void *FPGA_SRC_SHMEM_ADDR;
extern const char *FPGA_DMA_HOST_TO_DEV_CONNECTOR_ID;
extern const char *DPDK_FILE_PREFIX;
extern protocol OUTDST_PROTOCOL;
extern std::string OUTDST_IPA;
extern uint16_t OUTDST_PORT;

extern dma_info_t enqdmainfo_channel;
extern dma_info_t enqdmainfo;
extern dma_info_t deqdmainfo_channel;
extern dma_info_t deqdmainfo;
extern dmacmd_info_t enqdmacmdinfo;
extern dmacmd_info_t deqdmacmdinfo;

extern const chain_ctrl_cid_t g_param_chain_ctrl_tbls[CH_NUM_MAX];

//-----------------------------------------------------
// function
//-----------------------------------------------------
extern void prlog_mngque(const mngque_t &m);
extern void prlog_dma_info(const dma_info_t &i);
extern void prlog_dmacmd_info(const dmacmd_info_t &i, uint32_t enq_id);
extern void prlog_fpga_chain_ddr_info(const fpga_chain_ddr_t &i, uint32_t dev_id, uint32_t krnl_id, uint32_t extif_id);

#endif //CPU_DECODE_HPP
