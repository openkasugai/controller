/*****************************************************************
 * Copyright 2024 NTT Corporation, FUJITSU LIMITED
 *****************************************************************/
#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <stdbool.h>
#include <string.h>
#include <unistd.h>
#include <pthread.h>
#include <netinet/in.h>
#include "libshmem.h"
#include "libfpgactl.h"
#include "libdma.h"
#include "libdmacommon.h"
#include "liblldma.h"
#include "libchain.h"
#include "common.h"
#include "glue.h"
#include "glue_func.h"

int32_t glue(tcp_client_info_t tcp_client_info, char* connector_id, uint32_t width, uint32_t height)
{
	printf("--- glue start!! ---\n");
	//logfile(LOG_DEBUG, "--- glue start!! ---\n");

	int32_t ret = 0;
	int32_t errcnt = 0;
	mngque_t pque[CH_NUM_MAX];
	memset(&pque, 0, sizeof(mngque_t) * CH_NUM_MAX);
	char cmd[128];
	atomic_bool get_power_en = ATOMIC_VAR_INIT(false);

	int32_t fpga_num = fpga_get_num();
	if (fpga_num != 1) {
		printf(" Num of FPGA error(%d)\n", fpga_num);
		//logfile(LOG_ERROR, " Num of FPGA error(%d)\n", fpga_num);
		return -1;
	}
	uint32_t *dev_id = get_dev_id(0);

	// --- get options  ---
	uint32_t enq_num = getopt_enq_num();


	// deq_shmstate init
	for (size_t i=0; i < CH_NUM_MAX; i++) {
		bool *s = get_deq_shmstate(i);
		for (size_t j=0; j < SHMEMALLOC_NUM_MAX; j++) {
			s[j] = false;
			pthread_mutex_init( &tx_shmmutex[i][j], NULL );
		}
	}

	//----------------------------------------------
	// allocate buffer
	//----------------------------------------------
	if (glue_allocate_buffer() < 0) {
		// error
		return -1;
	}

	//----------------------------------------------
	// shared memory allocate
	//----------------------------------------------
	if (glue_shmem_allocate(SHMEM_DST, pque, width, height) < 0) {
		// error
		goto _END1;
	}

	//----------------------------------------------
	// fpga lldma queue setup (set dmainfo)
	//----------------------------------------------
	if (glue_dequeue_lldma_queue_setup(*dev_id, connector_id) < 0) {
		// error
		goto _END2;
	}


	//----------------------------------------------
	// DMA TX enqueue thread start
	//----------------------------------------------
	//printf("DMA TX enqueue thread start\n");
	//logfile(LOG_DEBUG, "--- pthread_create thread_dma_tx_enq ---\n");
	//rslt2file("\n--- dma tx enqueue thread start ---\n");
	pthread_t thread_dma_tx_enq_id[CH_NUM_MAX];
	//thread_enq_args_t th_dma_tx_enq_args[CH_NUM_MAX];
	thread_enq_args_t th_dma_tx_enq_args[CH_NUM_MAX];
	for (size_t i=0; i < CH_NUM_MAX; i++) {
		uint32_t ch_id = i;
		if (getopt_ch_en(ch_id)) {
			//th_dma_tx_enq_args[ch_id] = (thread_enq_args_t){*dev_id, ch_id, 0, enq_num};
			//ret = pthread_create(&thread_dma_tx_enq_id[ch_id], NULL, (void*)thread_dma_tx_enq, &th_dma_tx_enq_args[ch_id]);
			th_dma_tx_enq_args[ch_id] = (thread_enq_args_t){*dev_id, ch_id, 0, enq_num, pque};
			ret = pthread_create(&thread_dma_tx_enq_id[ch_id], NULL, (void*)thread_dma_tx_enq, &th_dma_tx_enq_args[ch_id]);
			if (ret) {
				printf(" CH(%d) create thread_dma_tx_enq error!(%d)\n", ch_id, ret);
				//logfile(LOG_ERROR, " CH(%zu) create thread_dma_tx_enq error!(%d)\n", ch_id, ret);
				// error
				goto _END3;
			}
			printf("CH(%d) thread_dma_tx_enq_id(%lx),\n", ch_id, thread_dma_tx_enq_id[ch_id]);
			//logfile(LOG_DEBUG, "CH(%zu) thread_dma_tx_enq_id(%lx),\n", ch_id, thread_dma_tx_enq_id[ch_id]);
		}
	}

	sleep(1);
	//----------------------------------------------
	// Send thread start
	//----------------------------------------------
	//printf("Send thread start\n");
	//logfile(LOG_DEBUG, "--- pthread_create thread_send ---\n");
	//rslt2file("\n--- send thread start ---\n");
	pthread_t thread_send_id[CH_NUM_MAX];
	//thread_receive_args_t thsend_args[CH_NUM_MAX];
	thread_send_args_t thsend_args[CH_NUM_MAX];
	for (size_t i=0; i < CH_NUM_MAX; i++) {
		uint32_t ch_id = i;
		if (getopt_ch_en(ch_id)) {
			//thsend_args[ch_id] = (thread_receive_args_t){ch_id, 0, enq_num};
			thsend_args[ch_id] = (thread_send_args_t){ch_id, 0, enq_num, tcp_client_info, width, height};
			//ret = pthread_create(&thread_send_id[ch_id], NULL, (void*)thread_receive, &thsend_args[ch_id]);
			ret = pthread_create(&thread_send_id[ch_id], NULL, (void*)thread_send, &thsend_args[ch_id]);
			if (ret) {
				printf(" CH(%d) create thread_send error!(%d)\n", ch_id, ret);
				//logfile(LOG_ERROR, " CH(%zu) create thread_send error!(%d)\n", ch_id, ret);
				// error
				goto _END3;
			}
			printf("CH(%d) thread_send_id(%lx),\n", ch_id, thread_send_id[ch_id]);
			//logfile(LOG_DEBUG, "CH(%zu) thread_send_id(%lx),\n", ch_id, thread_send_id[ch_id]);
		}
	}

	sleep(1);
	//----------------------------------------------
	// DMA TX dequeue thread start
	//----------------------------------------------
	//printf("DMA TX dequeue thread start\n");
	//logfile(LOG_DEBUG, "--- pthread_create thread_dma_tx_deq ---\n");
	//rslt2file("\n--- dma tx dequeue thread start ---\n");
	pthread_t thread_dma_tx_deq_id[CH_NUM_MAX];
	thread_deq_args_t th_dma_tx_deq_args[CH_NUM_MAX];
	for (size_t i=0; i < CH_NUM_MAX; i++) {
		uint32_t ch_id = i;
		if (getopt_ch_en(ch_id)) {
			th_dma_tx_deq_args[ch_id] = (thread_deq_args_t){*dev_id, ch_id, 0, enq_num};
            //ret = pthread_create(&thread_dma_tx_deq_id[ch_id], NULL, (void*)thread_dma_tx_deq, &th_dma_tx_deq_args[ch_id]);
            ret = pthread_create(&thread_dma_tx_deq_id[ch_id], NULL, (void*)thread_dma_tx_deq, &th_dma_tx_deq_args[ch_id]);
			if (ret) {
				printf(" CH(%d) create thread_dma_tx_deq error!(%d)\n", ch_id, ret);
				//logfile(LOG_ERROR, " CH(%zu) create thread_dma_tx_deq error!(%d)\n", ch_id, ret);
				// error
				goto _END3;
			}
			printf("CH(%d) thread_dma_tx_deq_id(%lx),\n", ch_id, thread_dma_tx_deq_id[ch_id]);
			//logfile(LOG_DEBUG, "CH(%zu) thread_dma_tx_deq_id(%lx),\n", ch_id, thread_dma_tx_deq_id[ch_id]);
		}
	}

	//----------------------------------------------
	// waitting... all finish
	//----------------------------------------------
	printf(" All threads started\n");
	//logfile(LOG_DEBUG, " ...waitting for all dequeue process to finish\n");
	//rslt2file("\n...waitting for all dequeue process to finish\n");

	// DMA TX enqueue thread end
	for (size_t i=0; i < CH_NUM_MAX; i++) {
		uint32_t ch_id = i;
		if (getopt_ch_en(ch_id)) {
			printf("CH(%d) pthread_join(thread_dma_tx_enq: %lx)\n", ch_id, thread_dma_tx_enq_id[ch_id]);
			//logfile(LOG_DEBUG, "CH(%zu) pthread_join(thread_dma_tx_enq: %lx)\n", ch_id, thread_dma_tx_enq_id[ch_id]);
			ret = pthread_join(thread_dma_tx_enq_id[ch_id], NULL);
			if (ret) {
				printf(" CH(%d) pthread_join thread_dma_tx_enq error!(%d)\n", ch_id, ret);
				//logfile(LOG_ERROR, " CH(%zu) pthread_join thread_dma_tx_enq error!(%d)\n", ch_id, ret);
			}
		}
	}

	// send thread end
	for (size_t i=0; i < CH_NUM_MAX; i++) {
		uint32_t ch_id = i;
		if (getopt_ch_en(ch_id)) {
			printf("CH(%d) pthread_join(thread_send: %lx)\n", ch_id, thread_send_id[ch_id]);
			//logfile(LOG_DEBUG, "CH(%zu) pthread_join(thread_send: %lx)\n", ch_id, thread_send_id[ch_id]);
			ret = pthread_join(thread_send_id[ch_id], NULL);
			if (ret) {
				printf(" CH(%d) pthread_join thread_send error!(%d)\n", ch_id, ret);
				//logfile(LOG_ERROR, " CH(%zu) pthread_join thread_send error!(%d)\n", ch_id, ret);
			}
		}
	}

	// DMA TX dequeue thread end
	for (size_t i=0; i < CH_NUM_MAX; i++) {
		uint32_t ch_id = i;
		if (getopt_ch_en(ch_id)) {
			printf("CH(%d) pthread_join(thread_dma_tx_deq: %lx)\n", ch_id, thread_dma_tx_deq_id[ch_id]);
			//logfile(LOG_DEBUG, "CH(%zu) pthread_join(thread_dma_tx_deq: %lx)\n", ch_id, thread_dma_tx_deq_id[ch_id]);
			ret = pthread_join(thread_dma_tx_deq_id[ch_id], NULL);
			if (ret) {
				printf(" CH(%d) pthread_join thread_dma_tx_deq error!(%d)\n", ch_id, ret);
				//logfile(LOG_ERROR, " CH(%zu) pthread_join thread_dma_tx_deq error!(%d)\n", ch_id, ret);
			}
		}
	}

	//----------------------------------------------
	// end processing
	//----------------------------------------------
_END3:
	// lldma queue finish
	glue_dequeue_lldma_queue_finish(*dev_id);

_END2:
	// shared memory free
	glue_shmem_free(pque);

_END1:
	// free buffer
	glue_free_buffer();


	printf("--- glue end!! ---\n");
	//logfile(LOG_DEBUG, "...glue end\n");

	return 0;
}
