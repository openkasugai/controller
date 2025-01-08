/*****************************************************************
 * Copyright 2024 NTT Corporation, FUJITSU LIMITED
 *****************************************************************/
#ifndef RECV_HEADER_LOG_HPP
#define RECV_HEADER_LOG_HPP

// header verification result
#ifdef MODULE_FPGA
typedef struct header_verify_info {
	int8_t result_marker;
	int8_t result_payload_len;
	int8_t result_data_id;
} header_verify_info_t;
#else
typedef struct header_verify_info {
	int8_t result_marker;
	int8_t result_payload_len;
	int8_t result_payload_type;
	int8_t result_channel_id;
	int8_t result_color_space;
	int8_t result_data_type;
	int8_t result_num_ch;
	int8_t result_width;
	int8_t result_height;
} header_verify_info_t;
#endif //MODULE_FPGA

namespace RecvHeaderLog {
	extern int32_t set(const FrameHeader_t &fh, size_t task);
	extern void close();
}

#endif //RECV_HEADER_LOG_HPP
