/*****************************************************************
 * Copyright 2024 NTT Corporation, FUJITSU LIMITED
 *****************************************************************/

#define	ALIGNMENT_ADDR     4096
#define	DATA_SIZE_1KB      1024
#define	ALIGNMENT_DMA_SIZE   64

#define	ENQUEUE_WAIT_MSEC			60000
#define	ENQUEUE_WAIT_MSEC_PERIOD		100
#define	DEQUEUE_WAIT_MSEC			60000
#define	DEQUEUE_WAIT_MSEC_PERIOD		100

// FPGA Frame Header
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

typedef struct array_buffer_info {
  void*    top_addr;
  size_t   buffer_size;
  void**   buffer_addrs;
} array_buffer_info_t;

int32_t wait_fpga_enqueue(dma_info_t *dmainfo, dmacmd_info_t *dmacmdinfo, const int32_t msec);
int32_t wait_fpga_dequeue(dma_info_t *dmainfo, dmacmd_info_t *dmacmdinfo, const int32_t msec);
void* alloc_array_buffer_host(size_t payload_size,
			      size_t array_size,
			      size_t *transfer_size_p,
			      size_t *buffer_size_p);
int32_t free_array_buffer_host(int32_t *addr);
void set_array_buffer_info(array_buffer_info_t *abi, size_t buffer_size, size_t array_size, void* addr);
