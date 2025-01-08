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

uint32_t read_buffer(uint32_t cmd_idx, void **data, int *size);
int clear_buffer(uint32_t cmd_idx);
int init_mem(char* file_prefix, bool shmem_secondary);
int init_worker(char* fpga_dev, char* file_prefix, char* connector_id);
int finish_worker();
uint64_t get_current_time();
