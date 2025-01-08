/*****************************************************************
 * Copyright 2024 NTT Corporation, FUJITSU LIMITED
 *****************************************************************/

#include <stdbool.h>
#include <stdint.h>

struct data_size {
  uint32_t input_height;
  uint32_t input_width;
  uint32_t output_height;
  uint32_t output_width;
};

int init_DPDK(char* file_prefix);
int fini_DPDK(char* file_prefix);
int init_FPGA(char* fpga_dev, int ch_id);
int fini_FPGA(int ch_id);
