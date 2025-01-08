/*****************************************************************
 * Copyright 2024 NTT Corporation, FUJITSU LIMITED
 *****************************************************************/

// FPGA frame header
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
  uint64_t local_ts;
  uint8_t reserved2[16];
} frameheader_t;
