/*****************************************************************
 * Copyright 2024 NTT Corporation, FUJITSU LIMITED
 *****************************************************************/
#include <stdio.h>
#include <stdint.h>
#include <stdbool.h>
#include <string.h>
#include <unistd.h>
#include <netinet/in.h>
#include <pthread.h>
#include "libshmem.h"
#include "libfpgactl.h"
#include "libdma.h"
#include "libdmacommon.h"
#include "liblldma.h"
#include "libchain.h"
#include "libfunction.h"
#include "common.h"
#include "glue.h"
#include "glue_func.h"

int32_t glue_shmem_allocate(shmem_mode_t shmem_mode, mngque_t *pque, uint32_t width, uint32_t height)
{
	int32_t ret = 0;

	printf("--- shmem_malloc ---\n");
	//logfile(LOG_DEBUG, "--- shmem_malloc ---\n");
	for (size_t i=0; i < CH_NUM_MAX; i++) {
		uint32_t ch_id = i;
		if (getopt_ch_en(ch_id)) {
			ret = shmem_malloc(shmem_mode, &pque[ch_id], ch_id, width, height);
			if (ret < 0) {
				printf("shmem_malloc error(%d)\n",ret);
				//logfile(LOG_ERROR, "shmem_malloc error(%d)\n",ret);
				return -1;
			}
			prlog_mngque(&pque[ch_id]);
		}
	}
}

int32_t glue_shmem_free(mngque_t *pque)
{
	printf("--- shmem_free ---\n");
	//logfile(LOG_DEBUG, "--- shmem_free ---\n");
	for (size_t i=0; i < CH_NUM_MAX; i++) {
		uint32_t ch_id = i;
		if (getopt_ch_en(ch_id)) {
			shmem_free(&pque[ch_id], ch_id);
		}
	}
}

int32_t glue_allocate_buffer(void)
{
	int32_t ret = 0;

	// allocate dmacmdinfo buffer
	printf("--- dmacmdinfo_malloc ---\n");
	//logfile(LOG_DEBUG, "--- dmacmdinfo_malloc ---\n");
	ret = dmacmdinfo_malloc();
	if (ret < 0) {
		printf("dmacmdinfo_alloc error!!!(%d)\n",ret);
		//logfile(LOG_ERROR, "dmacmdinfo_alloc error!!!(%d)\n",ret);
		return -1;
	}

	return 0;
}

void glue_free_buffer(void)
{
	int32_t ret = 0;

	// free dmacmdinfo buffer
	printf("--- dmacmdinfo_free ---\n");
	//logfile(LOG_DEBUG, "--- dmacmdinfo_free ---\n");
	dmacmdinfo_free();

}


int32_t glue_dequeue_lldma_queue_setup(uint32_t dev_id, char* connect_id)
{
	int32_t ret = 0;

	printf("--- dequeue fpga_lldma_queue_setup ---\n");
	//logfile(LOG_DEBUG, "--- dequeue fpga_lldma_queue_setup ---\n");
	for (size_t i=0; i < CH_NUM_MAX; i++) {
		uint32_t ch_id = i;
		if (getopt_ch_en(ch_id)) {
			//dma_info_t *pdmainfo_ch = get_deqdmainfo_channel(dev_id, ch_id);
			dma_info_t *pdmainfo = get_deqdmainfo(dev_id, ch_id);
			memset(pdmainfo, 0, sizeof(dma_info_t));
			char *connector_id = connect_id;
			//char *connector_id = pdmainfo_ch->connector_id;
			printf("dev(%d) CH(%d) dequeue fpga_lldma_queue_setup\n", dev_id, ch_id);
			//logfile(LOG_DEBUG, "dev(%zu) CH(%zu) dequeue fpga_lldma_queue_setup\n", dev_id, ch_id);
			ret = fpga_lldma_queue_setup(connector_id, pdmainfo);
			if (ret < 0) {
				printf("dequeue fpga_lldma_queue_setup error!!!(%d)\n",ret);
				//logfile(LOG_ERROR, "dequeue fpga_lldma_queue_setup error!!!(%d)\n",ret);
				// error
				return -1;
			}
			prlog_dma_info(pdmainfo, ch_id);
		}
	}

	return 0;
}


//int32_t glue_dequeue_set_dma_cmd(uint32_t run_id, uint32_t enq_num, mngque_t *pque)
//{
//	int32_t ret = 0;
//	const divide_que_t *div_que = get_divide_que();
//
//	printf("--- dequeue set_dma_cmd ---\n");
//	//logfile(LOG_DEBUG, "--- dequeue set_dma_cmd ---\n");
//	//rslt2file("\n--- dequeue set dma cmd ---\n");
//	for (size_t i=0; i < CH_NUM_MAX; i++) {
//		uint32_t ch_id = i;
//		if (getopt_ch_en(ch_id)) {
//			uint32_t data_len = pque[ch_id].dst1buflen;
//			uint32_t dsize = pque[ch_id].dst1dsize;
//			uint32_t dstidx = 0;
//			uint16_t taskidx = 1 + run_id * div_que->que_num;
//			printf("CH(%zu) dma tx data size=%d Byte\n", ch_id, dsize);
//			//rslt2file("CH(%zu) dma tx data size=%d Byte\n", ch_id, dsize);
//			for (size_t k=0; k < enq_num; k++) {
//				uint32_t enq_id = k + run_id * div_que->que_num;
//				uint16_t task_id = taskidx;
//				if (dstidx >= getopt_shmalloc_num()) {
//					dstidx = 0;
//				}
//				void *data_addr = pque[ch_id].enqbuf[dstidx].dst1bufp;
//				dstidx++;
//				dmacmd_info_t *pdmacmdinfo = get_deqdmacmdinfo(ch_id, enq_id);
//				memset(pdmacmdinfo, 0, sizeof(dmacmd_info_t));
//				printf("CH(%zu) DEQ(%zu) set_dma_cmd\n", ch_id, enq_id);
//				//logfile(LOG_DEBUG, "CH(%zu) DEQ(%zu) set_dma_cmd\n", ch_id, enq_id);
//				ret = set_dma_cmd(pdmacmdinfo, task_id, data_addr, data_len);
//				if (ret < 0) {
//					printf("dequeue set_dma_cmd error!!!(%d)\n",ret);
//					//logfile(LOG_ERROR, "dequeue set_dma_cmd error!!!(%d)\n",ret);
//					// error
//					return -1;
//				}
//				prlog_dmacmd_info(pdmacmdinfo, ch_id, enq_id);
//        
//				if (taskidx == 0xFFFF) {
//					taskidx = 1;
//				} else {
//					taskidx++;
//				}
//			}
//		}
//	}
//
//	return 0;
//}


void glue_dequeue_lldma_queue_finish(uint32_t dev_id)
{
	int32_t ret = 0;

	printf("--- dequeue fpga_lldma_queue_finish ---\n");
	//logfile(LOG_DEBUG, "--- dequeue fpga_lldma_queue_finish ---\n");
	for (size_t i=0; i < CH_NUM_MAX; i++) {
		uint32_t ch_id = i;
		if (getopt_ch_en(ch_id)) {
			printf("dev(%d) CH(%d) dequeue fpga_lldma_queue_finish\n", dev_id, ch_id);
			//logfile(LOG_DEBUG, "dev(%zu) CH(%zu) dequeue fpga_lldma_queue_finish\n", dev_id, ch_id);
			dma_info_t *pdmainfo = get_deqdmainfo(dev_id, ch_id);
			ret = fpga_lldma_queue_finish(pdmainfo);
			if (ret < 0) {
				printf("dequeue fpga_lldma_queue_finish error!!!(%d)\n",ret);
				//logfile(LOG_ERROR, "dequeue fpga_lldma_queue_finish error!!!(%d)\n",ret);
				// error
			}
			pdmainfo->connector_id = NULL;
			prlog_dma_info(pdmainfo, ch_id);
		}
	}
}
