/*****************************************************************
 * Copyright 2024 NTT Corporation, FUJITSU LIMITED
 *****************************************************************/
#define _GNU_SOURCE
#include <unistd.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>
#include <sched.h>              //CPU_XX
#include <pthread.h>
#include <sys/types.h>          //gettid
#include <sys/syscall.h>        //syscall
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include "libdma.h"
#include "libdmacommon.h"
#include "libpower.h"
#include "liblogging.h"
#include "common.h"
#include "glue.h"
//#include "cppfunc.h"

#define WAIT_TIME_DMA_TX_ENQUEUE	300000	//msec
#define WAIT_TIME_DMA_TX_DEQUEUE	300000	//msec

#define SHMEM_POLLING_INTERVAL	100 //usec

pthread_mutex_t tx_shmmutex[CH_NUM_MAX][SHMEMALLOC_NUM_MAX];

static int32_t wait_dma_tx_fpga_dequeue(dma_info_t *dmainfo, dmacmd_info_t *dmacmdinfo, const uint32_t enq_id, const int32_t msec);
static int32_t wait_dma_tx_fpga_enqueue(dma_info_t *dmainfo, dmacmd_info_t *dmacmdinfo, const uint32_t enq_id, const int32_t msec);

static int32_t wait_dma_tx_fpga_dequeue(dma_info_t *dmainfo, dmacmd_info_t *dmacmdinfo, const uint32_t enq_id, const int32_t msec)
{
	int32_t timeout = 100; //fpga_deuque timeout 100msec

	uint32_t ch_id = dmainfo->chid;
	uint16_t task_id = dmacmdinfo->task_id;

	int32_t cnt = msec/timeout;
	for (size_t i=0; i < cnt; i++) {
		int32_t ret = fpga_dequeue(dmainfo, dmacmdinfo);
		if (ret == 0) {
			return 0;
		}
	}

	printf("  CH(%u) deq(%u) task_id(%u) DMA TX dequeue timeout!!!\n", ch_id, enq_id, task_id);
	//logfile(LOG_ERROR, "  CH(%u) deq(%u) task_id(%u) DMA TX dequeue timeout!!!\n", ch_id, enq_id, task_id);

	return -1;
}

static int32_t wait_dma_tx_fpga_enqueue(dma_info_t *dmainfo, dmacmd_info_t *dmacmdinfo, const uint32_t enq_id, const int32_t msec)
{
	int32_t wait_msec = 100; // 100msec
	int32_t ret = 0;

	uint32_t ch_id = dmainfo->chid;
	uint32_t task_id = dmacmdinfo->task_id;

	int32_t cnt = msec/wait_msec;
	for (size_t i=0; i < cnt; i++) {
		ret = fpga_enqueue(dmainfo, dmacmdinfo);
		if (ret == 0) {
			return 0;
		} else if (ret == -ENQUEUE_QUEFULL) {
			//printf("  CH(%u) deq(%u) task_id(%u) DMA TX fpga_enqueue que full(%d)\n", ch_id, enq_id, task_id, ret);
			//logfile(LOG_DEBUG, "  CH(%u) deq(%u) task_id(%u) DMA TX fpga_enqueue que full(%d)\n", ch_id, enq_id, task_id, ret);
			usleep(wait_msec * 1000);
		} else {
			printf("  CH(%u) deq(%u) task_id(%u) DMA TX fpga_enqueue error!!!(%d)\n", ch_id, enq_id, task_id, ret);
			//logfile(LOG_ERROR, "  CH(%u) deq(%u) task_id(%u) DMA TX fpga_enqueue error!!!(%d)\n", ch_id, enq_id, task_id, ret);
			return ret;
		}
	}

	printf("  CH(%u) deq(%u) task_id(%u) DMA TX enqueue timeout!!!\n", ch_id, enq_id, task_id);
	//logfile(LOG_ERROR, "  CH(%u) deq(%u) task_id(%u) DMA TX enqueue timeout!!!\n", ch_id, enq_id, task_id);

	return -1;
}


//----------------------------------
// Send Thread
//----------------------------------
void thread_send(thread_send_args_t *args)
{
	int sockfd;
	struct sockaddr_in addr;
	int send_ret = 0;
	int send_ret2 = 0;

	uint32_t ch_id = args->ch_id;
	printf("CH(%u) ...thread_send start...\n", ch_id);
	//logfile(LOG_DEBUG, "CH(%u) ...thread_send start...\n", ch_id);


	// Destination address and port number setting
	addr.sin_family = AF_INET;
	addr.sin_port = htons(atoi(args->tcp_client_info.port));
	addr.sin_addr.s_addr = inet_addr( args->tcp_client_info.dst );

	//printf("TCP CLIENT (ip:%s, port:%d)\n",args->tcp_client_info.dst,atoi(args->tcp_client_info.port));


	// socket generation
	if( (sockfd = socket( AF_INET, SOCK_STREAM, 0) ) < 0 ) {
		perror( "socket" );
	}
	
	// server connection
	while (true) {
		// Wait for the connection to complete
		if(!connect( sockfd, (struct sockaddr *)&addr, sizeof( struct sockaddr_in ) )) {
			break;
		}
	}
	
	//const divide_que_t *div_que = get_divide_que();

	uint32_t *dev_id = get_dev_id(fpga_get_num() - 1);
    size_t height = args->height;
    size_t width = args->width;

	bool *deq_shms = get_deq_shmstate(ch_id);
	uint32_t ring = 0;

	uint32_t run_id = args->run_id;
	uint32_t enq_num = args->enq_num;

	size_t i = 0;
	while(1){
	//for (size_t i=0; i < enq_num; i++) {
		uint32_t enq_id = i;
		//uint32_t enq_id = i + run_id * div_que->que_num;
		const dmacmd_info_t *pdmacmdinfo = get_deqdmacmdinfo(ch_id, enq_id);
		uint32_t task_id = pdmacmdinfo-> task_id;
		void *data_addr = (void*)pdmacmdinfo->data_addr;

		bool lp = true;
		while (lp) {
			pthread_mutex_lock( &tx_shmmutex[ch_id][ring] );
			// Determine if shared memory has been dequeued [if true, shared memory can be read]
			bool ds = deq_shms[ring];
			pthread_mutex_unlock( &tx_shmmutex[ch_id][ring] );
			if (ds) {
				//if (! getopt_is_performance_meas())
				printf("CH(%u) deq(%d) task_id(%u) deq_shms[%u]=true send start\n", ch_id, enq_id, task_id, ring);
				//logfile(LOG_DEBUG, "CH(%u) deq(%zu) task_id(%u) deq_shms[%u]=true send start\n", ch_id, enq_id, task_id, ring);
				lp = false;
			}
			usleep(SHMEM_POLLING_INTERVAL);
		}

		if (pdmacmdinfo->result_task_id != 0) {
			// If not dequeue error
			//----------------------------------------------
			// send imagedata
			//----------------------------------------------
			send_ret = 0;
			send_ret2 = 0;
			size_t head_len = sizeof(frameheader_t);
			size_t img_len = height * width * 3;
			//printf("Glue TCP Send Image Size: width=%d, height=%d, 3ch, header_size=%d, total=%d\n",(int)width,(int)height,(int)head_len,(int)(img_len+head_len));
			send_ret = send( sockfd, data_addr, head_len+img_len, 0 );
			if(send_ret  < 0 ) {
				printf( "TCP send error\n" );
			}
			else if (send_ret < head_len+img_len) {
				// If the specified size cannot be sent, send the rest
				do{
					send_ret2 = send( sockfd, data_addr+send_ret, head_len+img_len-send_ret, 0 );
					if(send_ret2  < 0 ) {
						printf( "TCP send error\n" );
						break;
					} else {
						send_ret += send_ret2;
					}
				} while(send_ret < head_len+img_len);
			}

		}
		
		// deq shmstate update
		pthread_mutex_lock( &tx_shmmutex[ch_id][ring] );
		deq_shms[ring] = false; // Set to Unused Shared Memory [Change to false]
		pthread_mutex_unlock( &tx_shmmutex[ch_id][ring] );

		ring++;

		if (ring >= getopt_shmalloc_num()) {
			ring = 0;
		}

		//if (! getopt_is_performance_meas())
		printf("  CH(%u) deq(%d) task_id(%u) send end\n", ch_id, enq_id, task_id);
		//logfile(LOG_DEBUG, "  CH(%u) deq(%zu) task_id(%u) send end\n", ch_id, enq_id, task_id);

		i++;
		if(i==enq_num){
			i=0;
		}			
	}

	close(sockfd);
	printf("CH(%u) ...thread_glue end...\n", ch_id);
	//logfile(LOG_DEBUG, "CH(%u) ...thread_glue end...\n", ch_id);
}


//----------------------------------
// DMA TX Enqueue Thread
//----------------------------------
void thread_dma_tx_enq(thread_enq_args_t *args)
{
	uint32_t ch_id = args->ch_id;

	printf("CH(%u) ...thread_dma_tx_enq start...\n", ch_id);
	//logfile(LOG_DEBUG, "CH(%u) ...thread_dma_tx_enq start...\n", ch_id);

	//const divide_que_t *div_que = get_divide_que();

	uint32_t dev_id = args->dev_id;
	dma_info_t *pdmainfo = get_deqdmainfo(dev_id, ch_id);
	uint32_t run_id = args->run_id;
	uint32_t enq_num = args->enq_num;

	uint32_t enqueue_loop_flag = 0;
	size_t k = 0;
	uint32_t dstidx = 0;
	uint16_t taskidx = 1;
	while(1){
		int32_t ret = 0;
		//const divide_que_t *div_que = get_divide_que();

		//Issue DMA command
		//printf("--- dequeue set_dma_cmd ---\n");
		//logfile(LOG_DEBUG, "--- dequeue set_dma_cmd ---\n");
		//rslt2file("\n--- dequeue set dma cmd ---\n");
		//for (size_t i=0; i < CH_NUM_MAX; i++) {
		//uint32_t ch_id = i;
		if (getopt_ch_en(ch_id)) {
			uint32_t data_len = args->pque[ch_id].dst1buflen;
			uint32_t dsize = args->pque[ch_id].dst1dsize;
			uint32_t enq_id = k;//Enqueue ID loops at 0~99 (DMACMD space allocated)
			//uint32_t enq_id = k + run_id * div_que->que_num;//Enqueue ID loops at 0~99 (DMACMD space allocated)
			uint16_t task_id = taskidx;
			if (dstidx >= getopt_shmalloc_num()) {
				dstidx = 0;
			}
			void *data_addr = args->pque[ch_id].enqbuf[dstidx].dst1bufp;
			dstidx++;
			dmacmd_info_t *pdmacmdinfo = get_deqdmacmdinfo(ch_id, enq_id);
			memset(pdmacmdinfo, 0, sizeof(dmacmd_info_t));
			//printf("CH(%d) DEQ(%d) set_dma_cmd\n", ch_id, enq_id);
			//logfile(LOG_DEBUG, "CH(%zu) DEQ(%zu) set_dma_cmd\n", ch_id, enq_id);
			ret = set_dma_cmd(pdmacmdinfo, task_id, data_addr, data_len);
			if (ret < 0) {
				// error
				printf("enqueue set_dma_cmd error!!!(%d)\n",ret);
				//logfile(LOG_ERROR, "dequeue set_dma_cmd error!!!(%d)\n",ret);
				exit(0);
			}
			//prlog_dmacmd_info(pdmacmdinfo, ch_id, enq_id);
			if (taskidx == 0xFFFF) {
				taskidx = 1;
			} else {
				taskidx++;
			}

			//DMA Command Enqueue
			//uint32_t enq_id = k + run_id * div_que->que_num;
			//dmacmd_info_t *pdmacmdinfo = get_deqdmacmdinfo(ch_id, enq_id);
			ret = wait_dma_tx_fpga_enqueue(pdmainfo, pdmacmdinfo, enq_id, WAIT_TIME_DMA_TX_ENQUEUE);
			if (ret < 0) {
				printf("DMA TX enqerror CH(%u) enq(%d)\n", ch_id, enq_id);
				//logfile(LOG_ERROR, "DMA TX enqerror CH(%u) enq(%zu)\n", ch_id, enq_id);
			}

			k++;
			if(k==enq_num){
				k=0;
			}

		}

	}

	printf("CH(%u) ...thread_dma_tx_enq end...\n", ch_id);
	//logfile(LOG_DEBUG, "CH(%u) ...thread_dma_tx_enq end...\n", ch_id);
}



//----------------------------------
// DMA TX Dequeue Thread
//----------------------------------
void thread_dma_tx_deq(thread_deq_args_t *args)
{
	uint32_t ch_id = args->ch_id;
	printf("CH(%u) ...thread_dma_tx_deq start...\n", ch_id);
	//logfile(LOG_DEBUG, "CH(%u) ...thread_dma_tx_deq start...\n", ch_id);

	//const divide_que_t *div_que = get_divide_que();
	bool *deq_shms = get_deq_shmstate(ch_id);
	uint32_t ring = 0;

	uint32_t dev_id = args->dev_id;
	dma_info_t *pdmainfo = get_deqdmainfo(dev_id, ch_id);
	printf("CH(%u) DMA TX dma_info: dir(%d) chid(%u) queue_addr(%p) queue_size(%u)\n", ch_id, pdmainfo->dir, pdmainfo->chid, pdmainfo->queue_addr, pdmainfo->queue_size);
	//rslt2file("CH(%u) DMA TX dma_info: dir(%d) chid(%u) queue_addr(%p) queue_size(%u)\n", ch_id, pdmainfo->dir, pdmainfo->chid, pdmainfo->queue_addr, pdmainfo->queue_size);
	uint32_t run_id = args->run_id;

	size_t i = 0;
	while(1){

		uint32_t enq_num = args->enq_num;
		uint32_t enq_id = i;
		//uint32_t enq_id = i + run_id * div_que->que_num;
		//printf(" thread_dma_tx_deq(%u): deq(%d)\n", ch_id, enq_id);
		//logfile(LOG_DEBUG, " thread_dma_tx_deq(%u): deq(%zu)\n", ch_id, enq_id);
		dmacmd_info_t *pdmacmdinfo = get_deqdmacmdinfo(ch_id, enq_id);
		//prlog_dma_info(pdmainfo, ch_id);
		//prlog_dmacmd_info(pdmacmdinfo, ch_id, enq_id);
		uint16_t task_id = pdmacmdinfo-> task_id;
		bool lp = true;
		while (lp) {
			pthread_mutex_lock( &tx_shmmutex[ch_id][ring] );
			// Check shared memory usage [if false, shared memory can be written]
			bool ds = deq_shms[ring];
			pthread_mutex_unlock( &tx_shmmutex[ch_id][ring] );
			if (!ds) {
				//if (! getopt_is_performance_meas())
				printf("CH(%u) DMA TX deq(%d) task_id(%u) deq_shms[%u]=false dequeue start\n", ch_id, enq_id, task_id, ring);
				//logfile(LOG_DEBUG, "CH(%u) DMA TX deq(%zu) task_id(%u) deq_shms[%u]=false dequeue start\n", ch_id, enq_id, task_id, ring);
				lp = false;
			}
			usleep(SHMEM_POLLING_INTERVAL);
		}
		// dequeue data set
		//printf("CH(%u) DMA TX dmacmd_info: deq(%d) task_id(%u) dst_len(%u) dst_addr(%p)\n", ch_id, enq_id, pdmacmdinfo->task_id, pdmacmdinfo->data_len, pdmacmdinfo->data_addr);
		//rslt2file("CH(%u) DMA TX dmacmd_info: deq(%zu) task_id(%u) dst_len(%u) dst_addr(%p)\n", ch_id, enq_id, pdmacmdinfo->task_id, pdmacmdinfo->data_len, pdmacmdinfo->data_addr);
		int32_t ret = wait_dma_tx_fpga_dequeue(pdmainfo, pdmacmdinfo, enq_id, WAIT_TIME_DMA_TX_DEQUEUE);
		if (ret < 0) {
			printf("DMA TX deqerror CH(%u) deq(%d)\n", ch_id, enq_id);
			//logfile(LOG_ERROR, "DMA TX deqerror CH(%u) deq(%zu)\n", ch_id, enq_id);
		}
		// deq shmstate update
		pthread_mutex_lock( &tx_shmmutex[ch_id][ring] );
		deq_shms[ring] = true; // Shared memory dequeued [changed to true]
		pthread_mutex_unlock( &tx_shmmutex[ch_id][ring] );
		ring++;
		if (ring >= getopt_shmalloc_num()) {
			ring = 0;
		}
		//prlog_dma_info(pdmainfo, ch_id);
		//prlog_dmacmd_info(pdmacmdinfo, ch_id, enq_id);
		//printf("CH(%u) DMA TX dmacmd_info: deq(%d) result_task_id(%u) result_status(%u) result_data_len(%u)\n", ch_id, enq_id, pdmacmdinfo->result_task_id, pdmacmdinfo->result_status, pdmacmdinfo->result_data_len);
		//rslt2file("CH(%u) DMA TX dmacmd_info: deq(%zu) result_task_id(%u) result_status(%u) result_data_len(%u)\n", ch_id, enq_id, pdmacmdinfo->result_task_id, pdmacmdinfo->result_status, pdmacmdinfo->result_data_len);

		i++;
		if(i==enq_num){
			i=0;
		}			

	}

	printf("CH(%u) ...thread_dma_tx_deq end...\n", ch_id);
	//logfile(LOG_DEBUG, "CH(%u) ...thread_dma_tx_deq end...\n", ch_id);
}
