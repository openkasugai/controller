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

int write_buffer(uint32_t cmd_idx);
int read_buffer_physical(uint32_t cmd_idx);
int free_buffer(uint32_t cmd_idx);

int init_tester(int argc, char **argv, char* connector_id, char* fpga_dev, char* file_prefix);
int init_tester_physical(int argc, char **argv, char* connector_id, unsigned long paddr);
int finish_tester();

unsigned long get_current_time();
