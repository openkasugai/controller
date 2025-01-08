/*****************************************************************
 * Copyright 2024 NTT Corporation, FUJITSU LIMITED
 *****************************************************************/

#include <stddef.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#include "libshmem.h"
#include "libdma.h"
#include "libdmacommon.h"
#include "liblogging.h"

#include "util.h"

int32_t wait_fpga_enqueue(dma_info_t *dmainfo, dmacmd_info_t *dmacmdinfo, const int32_t msec) {
  int32_t wait_msec = ENQUEUE_WAIT_MSEC_PERIOD;
  int32_t ret = 0;
  int32_t cnt = msec/wait_msec;
  for (size_t i=0; i < cnt; i++) {
    ret = fpga_enqueue(dmainfo, dmacmdinfo);
    if (ret == 0) {
      return 0;
    } else if (ret == -ENQUEUE_QUEFULL) {
      usleep(wait_msec * 1000);
    } else {
      return ret;
    }
  }
  return -1;
}

int32_t wait_fpga_dequeue(dma_info_t *dmainfo, dmacmd_info_t *dmacmdinfo, const int32_t msec) {
  int32_t timeout = DEQUEUE_WAIT_MSEC_PERIOD;
  int32_t cnt = msec/timeout;
  for (size_t i=0; i < cnt; i++) {
    int32_t ret = fpga_dequeue(dmainfo, dmacmdinfo);
    if (ret == 0) {
      return 0;
    } else if (ret != -DEQUEUE_TIMEOUT) {
      return ret;
    }
  }
  return -1;
}

void* alloc_array_buffer_host(size_t payload_size,
			      size_t array_size,
			      size_t *transfer_size_p,
			      size_t *buffer_size_p) {

  size_t hd_size = sizeof(frameheader_t);
  size_t transfer_size = hd_size + payload_size;

  // Because the minimum transfer size is 1KB
  // If the size of the frame header + data is less than 1KB, correct it to 1KB
  if (transfer_size < DATA_SIZE_1KB){
    transfer_size = DATA_SIZE_1KB;
  } else {
    // Corrects to 64B alignment if > 1KB
    // size = ((size + 63)>>6)<<6)
    transfer_size = ((transfer_size + (ALIGNMENT_DMA_SIZE -1)) & (~(ALIGNMENT_DMA_SIZE -1)));
  }

  // 1 Record Transfer Size Per Data
  *transfer_size_p = transfer_size;
  
  // The destination address of the CPU must be 4KB aligned.
  // Round up transfer data size to 4KB to expand alignment space
  size_t buffer_size = ((transfer_size + (ALIGNMENT_ADDR -1))&(~(ALIGNMENT_ADDR -1)));
  *buffer_size_p = buffer_size;
  
  // Add data size by queue length
  size_t alloc_size = buffer_size * array_size;

  // Extra space for 4KB alignment of addresses to be registered in FPGA
  // (*) To shift the transfer destination backward when the first address is not aligned by 4KB
  alloc_size = alloc_size + ALIGNMENT_ADDR;
  void *addr = fpga_shmem_alloc(alloc_size);
  if (addr == NULL){
    printf("fpga_shmem_alloc failed\n");
    return NULL;
  }

  memset(addr, 0, alloc_size);
  return addr;
}

int32_t free_array_buffer_host(int32_t *addr) {
  fpga_shmem_free(addr);
  return 0;
}

void set_array_buffer_info(array_buffer_info_t *abi, size_t buffer_size, size_t array_size, void* addr) {
  abi->top_addr = addr;
  abi->buffer_size = buffer_size;
  abi->buffer_addrs = malloc(sizeof(void*)*array_size);
  void *top_addr_aligned = (void *)(((uint64_t)addr + (ALIGNMENT_ADDR-1)) & ~(ALIGNMENT_ADDR-1));
  for(int queue_id = 0; queue_id < array_size; queue_id++){
    void *buffer_addr = (void*)((uint64_t)top_addr_aligned + queue_id * buffer_size);
    abi->buffer_addrs[queue_id] = buffer_addr;
  }
}
