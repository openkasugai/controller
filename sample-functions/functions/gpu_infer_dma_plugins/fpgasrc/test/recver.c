/*****************************************************************
 * Copyright 2024 NTT Corporation, FUJITSU LIMITED
 *****************************************************************/

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#include "libdma.h"
#include "liblogging.h"
#include "libshmem.h"

#include "util.h"

#define NUM_FRAMES 8

static char* fpga_dev;
int width;
int height;
char* connector_id;

int main() {

  int ret = 0;
  char *strenv;

  int log_level;
  strenv = getenv("LOG_LEVEL");
  if (strenv == NULL) {
    log_level = LIBFPGA_LOG_INFO;
  } else {
    log_level = atoi(strenv);
  }
  libfpga_log_set_level(log_level);
  
  strenv = getenv("FPGA_DEV");
  if (strenv == NULL) {
    printf("environment variable FPGA_DEV must be specified\n");
    return -1;
  }
  fpga_dev = strenv;

  strenv = getenv("OUTPUT_HEIGHT");
  if (strenv == NULL) {
    printf("environment variable INPUT_HEIGHT must be specified\n");
    return -1;
  }
  height = atoi(strenv);

  strenv = getenv("OUTPUT_WIDTH");
  if (strenv == NULL) {
    printf("environment variable INPUT_WIDTH must be specified\n");
    return -1;
  }
  width = atoi(strenv);

  strenv = getenv("CONNECTOR_ID");
  if (strenv == NULL) {
    log_libfpga(LIBFPGA_LOG_ERROR,"environment variable CONNECTOR_ID must be specified in tester-mode\n");
    return -1;
  }
  connector_id = strenv;

  ret = fpga_shmem_init_sys("recver", NULL, NULL, NULL, 0);
  
  void*   addr;
  size_t  header_size;
  size_t  buffer_size;
  size_t  allocate_size;
  size_t  tx_size;

  int num_channels = 3;
  size_t payload_size = sizeof(uint8_t) * width * height * num_channels;

  size_t array_size = NUM_FRAMES;

  addr = alloc_array_buffer_host(payload_size, array_size, &tx_size, &buffer_size);
  if(addr == NULL) {
    printf("failed to allocate CPU memory\n");
    return -1;
  }

  array_buffer_info_t abi;
  set_array_buffer_info(&abi, buffer_size, array_size, addr);

  int dev_id;
  ret = fpga_dev_init(fpga_dev, &dev_id);
  if (ret < 0) {
    printf("fpga_dev_init failed: ret=%d\n", ret);
    return -1;
  }

  dma_info_t dmainfo;
  memset(&dmainfo, 0, sizeof(dma_info_t));
  ret = fpga_lldma_queue_setup(connector_id, &dmainfo);
  if (ret < 0) {
    printf("fpga_lldma_queue_setup failed: ret=%d\n", ret);
    return -1;
  }
  
  dmacmd_info_t dmacmdinfo[NUM_FRAMES];
  for (int queue_id = 0; queue_id < array_size; queue_id++) {
    printf("Recving data[%d]\n", queue_id);

    ret = set_dma_cmd(&dmacmdinfo[queue_id], queue_id, abi.buffer_addrs[queue_id], tx_size);
    if (ret < 0) {
      printf("set_dma_cmd failed: %d\n", ret);
      return -1;
    }

    ret = wait_fpga_enqueue(&dmainfo, &dmacmdinfo[queue_id], ENQUEUE_WAIT_MSEC);
    if (ret < 0) {
      printf("wait_fpga_enqueue failed ret=%d\n", ret);
      return -1;
    }

    ret = wait_fpga_dequeue(&dmainfo, &dmacmdinfo[queue_id], DEQUEUE_WAIT_MSEC);
    if (ret < 0) {
      printf("wait_fpga_dequeue failed ret=%d\n", ret);
      return -1;
    }
    frameheader_t fh = *(frameheader_t*)((uint64_t)dmacmdinfo[queue_id].data_addr);
    if (fh.payload_len == payload_size) {
      printf("OK! (%d/%d)\n", queue_id + 1, NUM_FRAMES);
    } else {
      printf("NG! (%d/%d)\n", queue_id + 1, NUM_FRAMES);
    }

    sleep(1);
  }

  printf("Recving done\n");
  
  return 0;
}
