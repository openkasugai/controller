/*****************************************************************
 * Copyright 2024 NTT Corporation, FUJITSU LIMITED
 *****************************************************************/

#include <stdint.h>

//-----------------------------------------------------
// frame header
//-----------------------------------------------------
#ifdef MODULE_FPGA
typedef struct frameheader {
	uint32_t marker;
	uint32_t payload_len;
	uint8_t reserved1[4];
	uint32_t sequence_num; // old frame_index
	uint8_t reserved2[8];
	double timestamp;
	uint32_t data_id; // old num_ch
	uint8_t reserved3[8];
	uint16_t header_checksum;
	uint8_t reserved4[2];
} frameheader_t;
#else
typedef struct frameheader {
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
} frameheader_t;
#endif // MODULE_FPGA