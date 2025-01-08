/*****************************************************************
 * Copyright 2024 NTT Corporation, FUJITSU LIMITED
 *****************************************************************/
#ifndef __GLUE_H__
#define __GLUE_H__

#include <stdint.h>
#include <stdatomic.h>
#include "libdmacommon.h"
#include "liblldma.h"
#include "libpower.h"


//-----------------------------------------------------
// define
//-----------------------------------------------------
//#define VERSION	"1.0.0"

#define DATA_SIZE_1KB	1024
#define DATA_SIZE_4KB	4096
#define CH_NUM_MAX	32
#define SHMEMALLOC_NUM_MAX	10
#define ALIGN_DST_LEN	64

#define LANE_NUM_MAX			4

typedef enum SHMEM_MODE {
	SHMEM_SRC,
	SHMEM_DST,
	SHMEM_SRC_DST,
	SHMEM_DST1_DST2,
	SHMEM_SRC_DST1_DST2,
	SHMEM_D2D_SRC,
	SHMEM_D2D_DST,
	SHMEM_D2D_SRC_DST,
	SHMEM_D2D
} shmem_mode_t; 


//-----------------------------------------------------
// frame header
//-----------------------------------------------------
#ifdef MODULE_FPGA
typedef struct frameheader {
	uint32_t marker;
	uint32_t payload_len;
	uint8_t reserved1[4];
	uint32_t sequence_num; // Old frame_index
	uint8_t reserved2[8];
	double timestamp;
	uint32_t data_id; // Old num_ch
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

//-----------------------------------------------------
// data format
//-----------------------------------------------------
typedef struct enqbuf {
	uint64_t *srcbufp;
	uint64_t *dst1bufp;
	uint64_t *dst2bufp;
} enqbuf_t;

typedef struct mngque {
	uint32_t enq_num;
	uint32_t srcdsize;
	uint32_t dst1dsize;
	uint32_t dst2dsize;
	uint32_t d2ddsize;
	uint32_t srcbuflen;
	uint32_t dst1buflen;
	uint32_t dst2buflen;
	uint32_t d2dbuflen;
	uint64_t *d2dbufp;
	enqbuf_t enqbuf[SHMEMALLOC_NUM_MAX];
} mngque_t;

typedef struct divide_que {
	uint32_t que_num; // Queues per Split
	uint32_t que_num_rem; // Divided Queue Remainder Count
	uint32_t div_num; // Number of partitions in the queue
} divide_que_t;

typedef struct tcp_client_info {
	char dst[128];
	char port[128];
} tcp_client_info_t;

//-----------------------------------------------------
// thread args
//-----------------------------------------------------
typedef struct thread_enq_args {
	uint32_t dev_id;
	uint32_t ch_id;
	uint32_t run_id;
	uint32_t enq_num;
	mngque_t *pque;
} thread_enq_args_t;

typedef struct thread_deq_args {
	uint32_t dev_id;
	uint32_t ch_id;
	uint32_t run_id;
	uint32_t enq_num;
} thread_deq_args_t;

typedef struct thread_receive_args {
	uint32_t ch_id;
	uint32_t run_id;
	uint32_t enq_num;
	tcp_client_info_t tcp_client_info;
	uint32_t width;
	uint32_t height;
} thread_send_args_t;

//-----------------------------------------------------
// mutex
//-----------------------------------------------------
extern pthread_mutex_t tx_shmmutex[CH_NUM_MAX][SHMEMALLOC_NUM_MAX];

//-----------------------------------------------------
// function
//-----------------------------------------------------
extern void glue_setting(void);
extern bool getopt_ch_en(uint32_t i);
extern uint32_t getopt_enq_num(void);
extern uint32_t getopt_shmalloc_num(void);
extern int32_t shmem_malloc(shmem_mode_t mode, mngque_t* p, uint32_t ch_id, uint32_t width, uint32_t height);
extern int32_t shmem_free(const mngque_t* p, uint32_t ch_id);
extern int32_t deq_shmem_malloc(mngque_t* p, uint32_t ch_id);
extern int32_t deq_shmem_free(const mngque_t* p, uint32_t ch_id);
extern bool* get_deq_shmstate(uint32_t ch_id);

extern void prlog_mngque(const mngque_t *p);
extern int32_t prlog_dma_info(const dma_info_t *p, uint32_t ch_id);
extern int32_t prlog_dmacmd_info(const dmacmd_info_t *p, uint32_t ch_id, uint32_t enq_id);

extern int32_t set_dev_id_list(void);
extern uint32_t* get_dev_id(uint32_t index);
extern uint32_t dev_id_to_index(uint32_t dev_id);
extern int32_t dmacmdinfo_malloc(void);
extern void dmacmdinfo_free(void);
extern dma_info_t* get_deqdmainfo(uint32_t dev_id, uint32_t ch_id);
extern dmacmd_info_t* get_deqdmacmdinfo(uint32_t ch_id, uint32_t enq_id);
extern const divide_que_t* get_divide_que(void);
extern void thread_dma_tx_deq(thread_deq_args_t *args);
extern void thread_dma_tx_enq(thread_enq_args_t *args);
extern void thread_send(thread_send_args_t *args);
extern void pr_device_info(void);

#endif /* __GLUE_H__ */
