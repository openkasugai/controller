/*****************************************************************
 * Copyright 2024 NTT Corporation, FUJITSU LIMITED
 *****************************************************************/
#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <string.h>
#include <unistd.h>
#include <getopt.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <assert.h>
#include <time.h>
#include "libshmem.h"
#include "liblldma.h"
#include "liblogging.h"
#include "common.h"
#include "glue.h"
//#include "cppfunc.h"
#include "glue_func.h"

static void* gmm[CH_NUM_MAX][SHMEMALLOC_NUM_MAX];
static bool gdeqshmstate[CH_NUM_MAX][SHMEMALLOC_NUM_MAX];

static bool g_ch_en[CH_NUM_MAX];
static int32_t g_ch_num[LANE_NUM_MAX];
static int32_t g_enq_num = 0;
static int32_t g_loglevel=LOG_ERROR;
static int32_t g_shmalloc_num = 0;
static divide_que_t g_divide_que;


int32_t getopt_loglevel(void)
{
	return g_loglevel;
}

uint32_t getopt_shmalloc_num(void)
{
	return g_shmalloc_num;
}

bool getopt_ch_en(uint32_t i)
{
	return g_ch_en[i];
}


uint32_t getopt_enq_num(void)
{
	return g_enq_num;
}


void glue_setting(void){

	//For debug
	g_loglevel = 2;

	//Setting the number of channels in a lane
	g_ch_num[0] = 1;
	//ch enable setting
	g_ch_en[0] = true;
	//Setting the number of allocated areas for dmacmd issuance
	g_enq_num = 100;
	//Setting the number of areas allocated for data
	g_shmalloc_num = SHMEMALLOC_NUM_MAX;

	g_divide_que.que_num = g_enq_num;
	g_divide_que.que_num_rem = 0;
	g_divide_que.div_num = 1;

}

//--------------------------
// device id
//--------------------------
static uint32_t dev_id_list[FPGA_MAX_DEVICES];
static bool set_dev_id_state = false;

int32_t set_dev_id_list()
{
	int32_t ret = 0;

	// get device list
	char **device_list;
	ret = fpga_get_device_list(&device_list);
	if (ret < 0) {
		printf("fpga_get_device_list:ret(%d) error!!\n", ret);
		//rslt2file("fpga_get_device_list error!!\n");
		//logfile(LOG_ERROR, "fpga_get_device_list:ret(%d) error!!\n", ret);
		return -1;
	}
	printf("fpga_get_device_list:ret(%d)\n", ret);
	//logfile(LOG_DEBUG, "fpga_get_device_list:ret(%d)\n", ret);

	// get dev_id from device list
	for (size_t i=0; i < fpga_get_num(); i++ ) {
		char str[64];
		sprintf(str, "%s%s", FPGA_DEVICE_PREFIX, device_list[i]);
		ret = fpga_get_dev_id(str, &dev_id_list[i]);
		if (ret < 0) {
			printf("fpga_get_dev_id:ret(%d) error!!\n", ret);
			//rslt2file("fpga_get_dev_id error!!\n");
			//logfile(LOG_ERROR, "fpga_get_dev_id:ret(%d) error!!\n", ret);
			return -1;
		}
		printf("fpga_get_dev_id:ret(%d)\n", ret);
		printf("  %s dev_id(%u)\n", str, dev_id_list[i]);
		//logfile(LOG_DEBUG, "fpga_get_dev_id:ret(%d)\n", ret);
		//logfile(LOG_DEBUG, "  %s dev_id(%u)\n", str, dev_id_list[i]);
	}
	set_dev_id_state = true;

	// release device list
	ret = fpga_release_device_list(device_list);
	if (ret < 0) {
		printf("fpga_release_device_list:ret(%d) error!!\n", ret);
		//rslt2file("fpga_release_device_list error!!\n");
		//logfile(LOG_ERROR, "fpga_release_device_list:ret(%d) error!!\n", ret);
		return -1;
	}
	printf("fpga_release_device_list:ret(%d)\n", ret);
	//logfile(LOG_DEBUG, "fpga_release_device_list:ret(%d)\n", ret);

	return 0;
}

uint32_t* get_dev_id(uint32_t index)
{
	assert(set_dev_id_state);

	uint32_t *p = &dev_id_list[index];

	return p;
}

//--------------------------
//  allocate shared memory
//--------------------------
int32_t shmem_malloc(shmem_mode_t mode, mngque_t* p, uint32_t ch_id, uint32_t width, uint32_t height)
{
	printf("CH(%d) shmem_malloc..(%p)\n", ch_id, p);
	//logfile(LOG_DEBUG, "CH(%d) shmem_malloc..(%p)\n", ch_id, p);

	if (p == NULL) {
		printf(" input error shmem_malloc..nil(%p)\n", p);
		//logfile(LOG_ERROR, " input error shmem_malloc..nil(%p)\n", p);
		return -1;
	}

	uint32_t enq_num = getopt_enq_num();
	uint32_t shmalloc_num = getopt_shmalloc_num();
	uint32_t imgsize_dst1 = height * width * 3;
	uint32_t headsize = sizeof(frameheader_t);
	uint32_t bufsize_dst1 = imgsize_dst1 + headsize;

	//alloc mem for queue
	printf("--- shmem alloc ---\n");
	//logfile(LOG_DEBUG, "--- shmem alloc ---\n");
	uint32_t ss = 0;
	p->enq_num = enq_num;
	p->srcdsize = 0;
	p->dst1dsize = 0;
	p->dst2dsize = 0;
	p->d2ddsize = 0;

	printf(" enq_num(%d), headsize(%d), imgsize_dst1(%d), bufsize_dst1(%d), shmalloc_num(%d)\n", enq_num, headsize, imgsize_dst1, bufsize_dst1, shmalloc_num);
	//logfile(LOG_DEBUG, " enq_num(%d), headsize(%d), imgsize_dst1(%d), bufsize_dst1(%d), shmalloc_num(%d)\n", enq_num, headsize, imgsize_dst1, bufsize_dst1, shmalloc_num);
	p->dst1dsize = bufsize_dst1;
	//ss = bufsize_dst1;                       //buf(dst1)
	ss = bufsize_dst1  + 0x10000;              //Extra space for 4KB alignment

	memset(&gmm[ch_id][0], 0, sizeof(gmm[ch_id]));
	printf("alloc..\n");
	//logfile(LOG_DEBUG, "alloc..\n");
 	for (size_t i=0; i < shmalloc_num; i++) {
		if (ss != 0) {
			printf("shmem alloc..(%d)\n",ss);
			//logfile(LOG_DEBUG, "shmem alloc..(%d)\n",ss);
			gmm[ch_id][i] = fpga_shmem_alloc(ss);
			if (gmm[ch_id][i] == NULL) {
				printf("shmemlloc error(%d)!\n",(int)i);
				//logfile(LOG_ERROR, "shmemlloc error(%d)!\n",i);
				p->enqbuf[i].srcbufp = NULL;
				p->enqbuf[i].dst1bufp = NULL;
				p->enqbuf[i].dst2bufp = NULL;
				return -1;
			}

			p->enqbuf[i].srcbufp = NULL;
			p->enqbuf[i].dst1bufp = (void*)(((uint64_t)gmm[ch_id][i] & ~0xfff)+0x1000);
			p->enqbuf[i].dst2bufp = NULL;

			printf("srcbufp(%p), dst1bufp(%p), dst2bufp(%p)\n", p->enqbuf[i].srcbufp, p->enqbuf[i].dst1bufp, p->enqbuf[i].dst2bufp);
			//logfile(LOG_DEBUG, "srcbufp(%p), dst1bufp(%p), dst2bufp(%p)\n", p->enqbuf[i].srcbufp, p->enqbuf[i].dst1bufp, p->enqbuf[i].dst2bufp);

		} else {
			p->enqbuf[i].srcbufp = NULL;
			p->enqbuf[i].dst1bufp = NULL;
			p->enqbuf[i].dst2bufp = NULL;
		}
	}


	//initialize data memory area
	for (size_t i=0; i < shmalloc_num; i++) {
		if (p->enqbuf[i].dst1bufp != NULL) {
			init_data((uint8_t*)p->enqbuf[i].dst1bufp, bufsize_dst1, 1); //0xff
		}
	}

	if (p->dst1dsize != 0) {
		if (bufsize_dst1 < DATA_SIZE_1KB) {
			// dequeue dst_len header + payload < 1KB set to 1KB
			p->dst1buflen = DATA_SIZE_1KB;
		} else {
			// dequeue dst_len alignment(ALIGN_DST_LEN byte)
			p->dst1buflen = (bufsize_dst1 + (ALIGN_DST_LEN - 1)) & ~(ALIGN_DST_LEN - 1);
		}
	}

	return 0;
}

int32_t shmem_free(const mngque_t* p, uint32_t ch_id)
{
	printf("CH(%u) shmem_free...\n", ch_id);
	//logfile(LOG_DEBUG, "CH(%u) shmem_free...\n", ch_id);

	uint32_t shmalloc_num = getopt_shmalloc_num();

	if (p != NULL) {
		printf("shmem_free(%p)\n", p);
		//logfile(LOG_DEBUG, "shmem_free(%p)\n", p);
		for (size_t i=0; i < shmalloc_num; i++) {
			if (p->enqbuf[i].srcbufp) {
				printf("shmemfree..(%p)\n",p->enqbuf[i].srcbufp);
				//logfile(LOG_DEBUG, "shmemfree..(%p)\n",p->enqbuf[i].srcbufp);
				fpga_shmem_free((void*)gmm[ch_id][i]);
			} else if (p->enqbuf[i].dst1bufp) {
				printf("shmemfree..(%p)\n",p->enqbuf[i].dst1bufp);
				//logfile(LOG_DEBUG, "shmemfree..(%p)\n",p->enqbuf[i].dst1bufp);
				fpga_shmem_free((void*)gmm[ch_id][i]);
			}
		}

	}

	return 0;
}

bool* get_deq_shmstate(uint32_t ch_id)
{
	bool *p = &gdeqshmstate[ch_id][0];

	return p;
}

//--------------------------
// queue info
//--------------------------
static dma_info_t deqdmainfo_channel[FPGA_MAX_DEVICES][CH_NUM_MAX];
static dma_info_t deqdmainfo[FPGA_MAX_DEVICES][CH_NUM_MAX];
static dmacmd_info_t **deqdmacmdinfo = NULL;

int32_t dmacmdinfo_malloc()
{
	printf("dmacmdinfo_malloc...\n");
	//logfile(LOG_DEBUG, "dmacmdinfo_malloc...\n");

	uint32_t enq_num = getopt_enq_num();

	// for dequeue
	deqdmacmdinfo = (dmacmd_info_t**)malloc(sizeof(dmacmd_info_t) * CH_NUM_MAX);
	if (deqdmacmdinfo == NULL) {
		printf("deqdmacmdinfo malloc error!\n");
		//logfile(LOG_ERROR, "deqdmacmdinfo malloc error!\n");
		return -1;
	}
	for (size_t i=0; i < CH_NUM_MAX; i++) {
		deqdmacmdinfo[i] = (dmacmd_info_t*)malloc(sizeof(dmacmd_info_t) * enq_num);
		if (deqdmacmdinfo[i] == NULL) {
			printf("deqdmacmdinfo[%d] malloc error!\n", (int)i);
			//logfile(LOG_ERROR, "deqdmacmdinfo[%zu] malloc error!\n", i);
			return -1;
		}
	}
	printf("  deqdmacmdinfo malloc(%p)\n", deqdmacmdinfo);
	//logfile(LOG_DEBUG, "  deqdmacmdinfo malloc(%p)\n", deqdmacmdinfo);

	return 0;
}

void dmacmdinfo_free()
{
	printf("dmacmdinfo_free...\n");
	//logfile(LOG_DEBUG, "dmacmdinfo_free...\n");

	// for dequeue
	if (deqdmacmdinfo != NULL) {
		for (size_t i=0; i < CH_NUM_MAX; i++) {
			if (deqdmacmdinfo[i] != NULL) {
				printf("  deqdmacmdinfo[%d] free(%p)\n", (int)i, deqdmacmdinfo[i]);
				//logfile(LOG_DEBUG, "  deqdmacmdinfo[%zu] free(%p)\n", i, deqdmacmdinfo[i]);
				free(deqdmacmdinfo[i]);
			}
		}
		printf("  deqdmacmdinfo free(%p)\n", deqdmacmdinfo);
		//logfile(LOG_DEBUG, "  deqdmacmdinfo free(%p)\n", deqdmacmdinfo);
		free(deqdmacmdinfo);
	} else {
		printf("  deqdmacmdinfo buffer is NULL!\n");
		//logfile(LOG_ERROR, "  deqdmacmdinfo buffer is NULL!\n");
	}
}

dma_info_t* get_deqdmainfo(uint32_t dev_id, uint32_t ch_id)
{
	dma_info_t* p;

	p = &deqdmainfo[dev_id][ch_id];

	return p;
}

dmacmd_info_t* get_deqdmacmdinfo(uint32_t ch_id, uint32_t enq_id)
{
	dmacmd_info_t* p;

	p = &deqdmacmdinfo[ch_id][enq_id];

	return p;
}

const divide_que_t* get_divide_que(void)
{
	divide_que_t* p;

	p = &g_divide_que;

	return p;
}



//-----------------------------------------
// For viewing debug logs
//-----------------------------------------
int32_t prlog_dmacmd_info(const dmacmd_info_t *p, uint32_t ch_id, uint32_t enq_id)
{
	if (p == NULL) {
		return -1;
	}
	printf("CH(%u) ENQ(%u) pr_dmacmd_info(%p)\n", ch_id, enq_id, p);
	printf("  CH(%u) ENQ(%u) task_id(0x%x)\n", ch_id, enq_id, p->task_id);
	printf("  CH(%u) ENQ(%u) data_len(0x%x)\n", ch_id, enq_id, p->data_len);
	printf("  CH(%u) ENQ(%u) data_addr(%p)\n", ch_id, enq_id, p->data_addr);
	printf("  CH(%u) ENQ(%u) desc_addr(%p)\n", ch_id, enq_id, p->desc_addr);
	printf("  CH(%u) ENQ(%u) result_status(%u)\n", ch_id, enq_id, p->result_status);
	printf("  CH(%u) ENQ(%u) result_task_id(0x%x)\n", ch_id, enq_id, p->result_task_id);
	printf("  CH(%u) ENQ(%u) result_data_len(0x%x)\n", ch_id, enq_id, p->result_data_len);
	printf("  CH(%u) ENQ(%u) result_data_addr(%p)\n", ch_id, enq_id, p->result_data_addr);

	//logfile(LOG_DEBUG, "CH(%u) ENQ(%u) pr_dmacmd_info(%p)\n", ch_id, enq_id, p);
	//logfile(LOG_DEBUG, "  CH(%u) ENQ(%u) task_id(0x%x)\n", ch_id, enq_id, p->task_id);
	//logfile(LOG_DEBUG, "  CH(%u) ENQ(%u) data_len(0x%x)\n", ch_id, enq_id, p->data_len);
	//logfile(LOG_DEBUG, "  CH(%u) ENQ(%u) data_addr(%p)\n", ch_id, enq_id, p->data_addr);
	//logfile(LOG_DEBUG, "  CH(%u) ENQ(%u) desc_addr(%p)\n", ch_id, enq_id, p->desc_addr);
	//logfile(LOG_DEBUG, "  CH(%u) ENQ(%u) result_status(%u)\n", ch_id, enq_id, p->result_status);
	//logfile(LOG_DEBUG, "  CH(%u) ENQ(%u) result_task_id(0x%x)\n", ch_id, enq_id, p->result_task_id);
	//logfile(LOG_DEBUG, "  CH(%u) ENQ(%u) result_data_len(0x%x)\n", ch_id, enq_id, p->result_data_len);
	//logfile(LOG_DEBUG, "  CH(%u) ENQ(%u) result_data_addr(%p)\n", ch_id, enq_id, p->result_data_addr);

	return 0;
}


int32_t prlog_dma_info(const dma_info_t *p, uint32_t ch_id)
{
	if (p == NULL) {
		return -1;
	}
	printf("CH(%u) pr_dma_info(%p)\n", ch_id, p);
	printf("  CH(%u) dev_id(0x%x)\n", ch_id, p->dev_id);
	printf("  CH(%u) dir(%d)\n", ch_id, p->dir);
	printf("  CH(%u) chid(0x%x)\n", ch_id, p->chid);
	printf("  CH(%u) queue_addr(%p)\n", ch_id, p->queue_addr);
	printf("  CH(%u) queue_size(%u)\n", ch_id, p->queue_size);

	//logfile(LOG_DEBUG, "CH(%u) pr_dma_info(%p)\n", ch_id, p);
	//logfile(LOG_DEBUG, "  CH(%u) dev_id(0x%x)\n", ch_id, p->dev_id);
	//logfile(LOG_DEBUG, "  CH(%u) dir(%d)\n", ch_id, p->dir);
	//logfile(LOG_DEBUG, "  CH(%u) chid(0x%x)\n", ch_id, p->chid);
	//logfile(LOG_DEBUG, "  CH(%u) queue_addr(%p)\n", ch_id, p->queue_addr);
	//logfile(LOG_DEBUG, "  CH(%u) queue_size(%u)\n", ch_id, p->queue_size);

	if (p->connector_id == NULL) {
		printf("  CH(%u) connector_id(%s)\n", ch_id, "");
		//logfile(LOG_DEBUG, "  CH(%u) connector_id(%s)\n", ch_id, "");
	} else {
		printf("  CH(%u) connector_id(%s)\n", ch_id, p->connector_id);
		//logfile(LOG_DEBUG, "  CH(%u) connector_id(%s)\n", ch_id, p->connector_id);
	}

	return 0;
}

void prlog_mngque(const mngque_t *p)
{
	printf("pr_mngque...\n");
	//logfile(LOG_DEBUG, "pr_mngque...\n");

	if (p != NULL) {
		printf("pr_mngque(%p)\n",p);
		printf("  enq_num(%d)\n",p->enq_num);
		printf("  srcdsize(0x%x)\n",p->srcdsize);
		printf("  dst1dsize(0x%x)\n",p->dst1dsize);
		printf("  dst2dsize(0x%x)\n",p->dst2dsize);
		printf("  d2ddsize(0x%x)\n",p->d2ddsize);
		printf("  srcbuflen(0x%x)\n",p->srcbuflen);
		printf("  dst1buflen(0x%x)\n",p->dst1buflen);
		printf("  dst2buflen(0x%x)\n",p->dst2buflen);
		printf("  d2dbuflen(0x%x)\n",p->d2dbuflen);
		printf("  d2dbufp(%p)\n",p->d2dbufp);

		//logfile(LOG_DEBUG, "pr_mngque(%p)\n",p);
		//logfile(LOG_DEBUG, "  enq_num(%d)\n",p->enq_num);
		//logfile(LOG_DEBUG, "  srcdsize(0x%x)\n",p->srcdsize);
		//logfile(LOG_DEBUG, "  dst1dsize(0x%x)\n",p->dst1dsize);
		//logfile(LOG_DEBUG, "  dst2dsize(0x%x)\n",p->dst2dsize);
		//logfile(LOG_DEBUG, "  d2ddsize(0x%x)\n",p->d2ddsize);
		//logfile(LOG_DEBUG, "  srcbuflen(0x%x)\n",p->srcbuflen);
		//logfile(LOG_DEBUG, "  dst1buflen(0x%x)\n",p->dst1buflen);
		//logfile(LOG_DEBUG, "  dst2buflen(0x%x)\n",p->dst2buflen);
		//logfile(LOG_DEBUG, "  d2dbuflen(0x%x)\n",p->d2dbuflen);
		//logfile(LOG_DEBUG, "  d2dbufp(%p)\n",p->d2dbufp);

		for (size_t i=0; i < getopt_shmalloc_num(); i++) {
			printf("  [%d] srcbufp(%p)\n", (int)i, p->enqbuf[i].srcbufp);
			printf("  [%d] dst1bufp(%p)\n", (int)i, p->enqbuf[i].dst1bufp);
			printf("  [%d] dst2bufp(%p)\n", (int)i, p->enqbuf[i].dst2bufp);
			//logfile(LOG_DEBUG, "  [%zu] srcbufp(%p)\n", i, p->enqbuf[i].srcbufp);
			//logfile(LOG_DEBUG, "  [%zu] dst1bufp(%p)\n", i, p->enqbuf[i].dst1bufp);
			//logfile(LOG_DEBUG, "  [%zu] dst2bufp(%p)\n", i, p->enqbuf[i].dst2bufp);
		}
	}
}
