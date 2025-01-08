/*****************************************************************
 * Copyright 2024 NTT Corporation, FUJITSU LIMITED
 *****************************************************************/

#include <assert.h>
#include <stdio.h>
#include <string.h>
#include <sys/time.h>
#include <unistd.h>

#include "libdma.h"
#include "libdmacommon.h"
#include "liblogging.h"
#include "libshmem.h"

#include "lib_worker.h"
#include "util.h"

#define GST_PLUGIN

#define BUFFER_SIZE 64
#define WAIT_TIME_DEQUEUE       5000    //msec

uint32_t dev_id = 0;

#define MARKER 0xE0FF10AD

struct data_size ds;
static uint32_t lane_num = 1;
static uint32_t lane_id = 0;
static dma_info_t dmainfo;
static dmacmd_info_t dmacmdinfo[BUFFER_SIZE];

uint64_t get_current_time() {
  struct timeval tv;
  int rc = gettimeofday(&tv, NULL);
  assert(rc == 0);
  unsigned long current_time = 1UL*1000*1000*1000*tv.tv_sec + 1UL*1000*tv.tv_usec;
  return current_time;
}

uint32_t read_buffer(uint32_t cmd_idx, void **data, int *size) {
  int ret;
  dma_info_t *pdmainfo = &dmainfo;
  dmacmd_info_t *pdmacmdinfo = &dmacmdinfo[cmd_idx];
  ret = wait_fpga_dequeue(pdmainfo, pdmacmdinfo, WAIT_TIME_DEQUEUE);
  if (ret < 0) {
    log_libfpga(LIBFPGA_LOG_WARN, "wait_fpga_dequeue failed: ret=%d\n", ret);
    return -1;
  }
  frameheader_t *fh = (frameheader_t*)((uint64_t)pdmacmdinfo->data_addr);
  if (fh->marker == MARKER) {
    uint64_t image_size = fh->payload_len;
    uint64_t exepected_size = ds.output_height * ds.output_width * 3;
    if (image_size == exepected_size) {
      log_libfpga(LIBFPGA_LOG_INFO, "OK!: expected size=%lu, received size=%lu, frame_index=%u\n", exepected_size, image_size, fh->frame_index);
      *size = image_size;
      *data = (void*)(pdmacmdinfo->data_addr + sizeof(frameheader_t));
    } else {
      log_libfpga(LIBFPGA_LOG_WARN, "NG!: expected size=%lu, received size=%lu, frame_index=%u\n", exepected_size, image_size, fh->frame_index);
      *size = -1;
      *data = NULL;
    }
  } else {
    printf("NG!: invalid data format\n");
    *size = -1;
    *data = NULL;
  }
  uint32_t next_cmd_idx = cmd_idx + 1;
  uint32_t deq_num = BUFFER_SIZE;
  return next_cmd_idx % deq_num;
}

int clear_buffer(uint32_t cmd_idx) {
  int ret;
  dma_info_t *pdmainfo = &dmainfo;
  dmacmd_info_t *pdmacmdinfo = &dmacmdinfo[cmd_idx];
  frameheader_t *fh = (frameheader_t*)((uint64_t)pdmacmdinfo->data_addr);
  fh->payload_len = 0;
  ret = fpga_enqueue(pdmainfo, pdmacmdinfo);
  if (ret < 0) {
    log_libfpga(LIBFPGA_LOG_ERROR, "fpga_enqueue failed: ret=%d\n", ret);
    return -1;
  }
  return 0;
}

int init_mem(char* file_prefix, bool shmem_secondary) {
  // initialize DPDK
  int ret;
  if (shmem_secondary) {
    ret = fpga_shmem_init(file_prefix, NULL, 0);
    if (ret < 0) {
      log_libfpga(LIBFPGA_LOG_ERROR, "fpga_shmem_init error!! ret=%d\n", ret);
      return -1;
    }
  } else {
    ret = fpga_shmem_init_sys(file_prefix, NULL, NULL, NULL, 0);
    if (ret < 0) {
      log_libfpga(LIBFPGA_LOG_ERROR, "fpga_shmem_init_sys error!! ret=%d\n", ret);
      return -1;
    }
  }
  log_libfpga(LIBFPGA_LOG_INFO, "fpga_shmem_init:ret(%d)\n",ret);
  return 0;
}

int init_worker(char* fpga_device, char* file_prefix, char* connector_id) {
  int ret;
  ret = fpga_dev_init(fpga_device, &dev_id);
  if (ret < 0) {
    printf("fpga_dev_init failed: ret=%d\n", ret);
    return -1;
  }
  ret = fpga_enable_regrw(dev_id);
  if (ret < 0) {
    printf("fpga_enable_regrw failed: ret=%d\n", ret);
    return -1;
  }
  int32_t fpga_num = fpga_get_num();
  if (fpga_num != 1) {
    printf(" Num of FPGA error(%d)\n", fpga_num);
    return -1;
  }
  void*   addr;
  size_t  header_size;
  size_t  buffer_size;
  size_t  allocate_size;
  size_t  tx_size_cpu;
  size_t  allocate_size_cpu;
  int num_channels = 3;
  size_t payload_size = sizeof(uint8_t) * ds.output_width * ds.output_height * num_channels;
  size_t array_size = BUFFER_SIZE;
  addr = alloc_array_buffer_host(payload_size, array_size, &tx_size_cpu, &buffer_size);
  if(addr == NULL) {
    printf("failed to allocate CPU memory\n");
    return -1;
  }
  array_buffer_info_t abi;
  set_array_buffer_info(&abi, buffer_size, array_size, addr);
  for (int queue_id = 0; queue_id < array_size; queue_id++) {
    ret = set_dma_cmd(&dmacmdinfo[queue_id], queue_id, abi.buffer_addrs[queue_id], tx_size_cpu);
  }
  memset(&dmainfo, 0, sizeof(dma_info_t));
  ret = fpga_lldma_queue_setup(connector_id, &dmainfo);
  if (ret < 0) {
    printf("fpga_lldma_queue_setup failed: ret=%d\n",ret);
    return -1;
  }
  // for debugging purpose
  fpga_queue_t *queue = dmainfo.queue_addr;
  fpga_desc_t *desc;
  int invalid_cmds, ready_cmds, done_cmds;
  invalid_cmds = ready_cmds = done_cmds = 0;
  for (int i = 0; i < queue->size; i++) {
    desc = &queue->ring[i];
    if (desc->op == CMD_INVALID) {
      invalid_cmds++;
    } else if (desc->op == CMD_READY) {
      ready_cmds++;
    } else {
      done_cmds++;
    }
  }
  printf("queue info: readhead=%d, writehead=%d, (invalid, ready, done) = (%d, %d, %d)\n",
	 queue->readhead, queue->writehead, invalid_cmds, ready_cmds, done_cmds);

  for (size_t j =0; j < BUFFER_SIZE; j++) {
    dmacmd_info_t *pdmacmdinfo = &dmacmdinfo[j];
    frameheader_t *fh = (frameheader_t*)((uint64_t)pdmacmdinfo->data_addr);
    fh->payload_len = 0;
    ret = fpga_enqueue(&dmainfo, pdmacmdinfo);
    if (ret < 0) {
      printf("fpga_enqueue failed: ret=%d\n",ret);
      return -1;
    }
  }
  return 0;
}

int finish_worker() {
  int ret;
  dma_info_t *pdmainfo = &dmainfo;
  // clear queue
  int wait_time = 30;
  int count = 0;
  fpga_queue_t *queue = pdmainfo->queue_addr;
  fpga_desc_t *desc;
  int invalid_cmds, ready_cmds, done_cmds;
  while (true) {
    invalid_cmds = ready_cmds = done_cmds = 0;
    for (int i = 0; i < queue->size; i++) {
      desc = &queue->ring[i];
      if (desc->op == CMD_INVALID) {
	invalid_cmds++;
      } else if (desc->op == CMD_READY) {
	ready_cmds++;
      } else {
	done_cmds++;
      }
    }
    if (ready_cmds == 0) {
      dmacmd_info_t dummy_dmacmdinfo;
      for (int i = 0; i < done_cmds; i++) {
	do {
	  int ret = fpga_dequeue(pdmainfo, &dummy_dmacmdinfo);
	} while (ret != 0);
      }
      break;
    } else {
      if (count < wait_time) {
	printf("waiting for all commands to get done...\n");
	printf("queue info: readhead=%d, writehead=%d, (invalid, ready, done) = (%d, %d, %d)\n",
	       queue->readhead, queue->writehead, invalid_cmds, ready_cmds, done_cmds);
	sleep(1);
      } else {
	printf("time out, queue needs to be re-initialized because of garbage commands\n");
	break;
      }
      count++;
    }
  }
  log_libfpga(LIBFPGA_LOG_INFO, "//--- TEST END ---//\n");
  ret = fpga_lldma_queue_finish(pdmainfo);
  if (ret < 0) {
    log_libfpga(LIBFPGA_LOG_ERROR, "fpga_lldma_queue_finish error!!!(%d)\n",ret);
  }
  ret = fpga_disable_regrw(dev_id);
  if (ret < 0) {
    log_libfpga(LIBFPGA_LOG_ERROR, "fpga_disable_regrw error!!\n");
  }
  log_libfpga(LIBFPGA_LOG_INFO, "fpga_disable_regrw:ret(%d)\n",ret);
  // FPGA close
  ret = fpga_finish();
  if (ret < 0) {
    log_libfpga(LIBFPGA_LOG_ERROR, "fpga finish error!!\n");
  }
  log_libfpga(LIBFPGA_LOG_INFO, "fpga_finish:ret(%d)\n",ret);
  // finish DPDK shmem
  ret = fpga_shmem_finish();
  if (ret < 0) {
    log_libfpga(LIBFPGA_LOG_ERROR, "fpga shmem finish error!!\n");
  }
  log_libfpga(LIBFPGA_LOG_INFO, "fpga_shmem_finish:ret(%d)\n",ret);
  return 0;
}
