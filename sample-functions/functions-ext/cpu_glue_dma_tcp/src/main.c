/*****************************************************************
 * Copyright 2024 NTT Corporation, FUJITSU LIMITED
 *****************************************************************/
#include <unistd.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/time.h>
#include <netinet/in.h>
#include <rte_common.h>
#include <rte_eal.h>
#include "libshmem.h"
#include "libfpgabs.h"
#include "libdma.h"
#include "libdmacommon.h"
#include "liblldma.h"
#include "libchain.h"
#include "liblogging.h"
#include "common.h"
#include "glue.h"
#include "glue_func.h"
#include <string.h>

static char* fpga_dev;
static char* file_prefix;
static char* secondary_mode;
char* connector_id;
uint32_t dev_id = 0;
bool shmem_secondary = true;

int init_mem(char* file_prefix, bool shmem_secondary) {
	// initialize DPDK
	int ret;
	if (shmem_secondary) {
		printf("start fpga_shmem_init @secondary\n");
    	ret = fpga_shmem_init(file_prefix, NULL, 0);
    	if (ret < 0) {
			printf("fpga_shmem_init error\n");
    		//log_libfpga(LIBFPGA_LOG_ERROR, "fpga_shmem_init error!! ret=%d\n", ret);
    		return -1;
    	}
	} else {
		//for debug
		printf("start fpga_shmem_init_sys @primary\n");
    	ret = fpga_shmem_init_sys(file_prefix, NULL, NULL, NULL, 0);
    	if (ret < 0) {
			printf("fpga_shmem_init_sys error\n");
			//log_libfpga(LIBFPGA_LOG_ERROR, "fpga_shmem_init_sys error!! ret=%d\n", ret);
			return -1;
		}
  	}
	printf("fpga_shmem_init:ret(%d)\n",ret);
	//log_libfpga(LIBFPGA_LOG_INFO, "fpga_shmem_init:ret(%d)\n",ret);
	return 0;
}


int main(int argc, char **argv)
{
	int32_t ret = 0;
	char *strenv;


	if (argc != 4) {
		printf("Invalid Argments!\n");
		exit(0);
	}

	//Allocating memory for storing TCP communication information
	tcp_client_info_t tcp_client_info;

	//Acquisition of argument information
	//Acquisition of information for TCP transmission
	char tcp_dst_info[128];
	strcpy(tcp_dst_info,argv[1]);
	printf("tcp_dst_info : %s\n",tcp_dst_info);
	char* colon_pos = strchr(tcp_dst_info,':');
	strncpy(tcp_client_info.dst,tcp_dst_info,(size_t)(colon_pos-tcp_dst_info));
	tcp_client_info.dst[(size_t)(colon_pos-tcp_dst_info)]='\0';
	//strcpy(tcp_client_info.dst,"127.0.0.1");
	strcpy(tcp_client_info.port,colon_pos+1);
	//Acquisition of input width information
	uint32_t input_width = atoi(argv[2]);
	//Acquisition of input height information
	uint32_t input_height = atoi(argv[3]);

	printf("TCP CLIENT INFO (ip:%s, port:%d)\n",tcp_client_info.dst,atoi(tcp_client_info.port));
	printf("INPUT WIDTH  (%d)\n",input_width);
	printf("INPUT HEIGHT (%d)\n",input_height);


	libfpga_log_set_level(LIBFPGA_LOG_ERROR);
	//rslt2file("\nVersion: %s\n", VERSION);

	//Acquisition of environment variable information
	//for debug
	strenv = getenv("SHMEM_SECONDARY");
	secondary_mode = strenv;

	//Acquisition of environment variable information
	strenv = getenv("GLUEENV_DPDK_FILE_PREFIX");
	if (strenv == NULL) {
		printf("GLUEENV_DPDK_FILE_PREFIX is NULL\n");
		exit(0);
	}
	printf("GLUEENV_DPDK_FILE_PREFIX=%s\n",strenv);
	file_prefix = strenv;
	
	//Acquisition of environment variable information
	strenv = getenv("GLUEENV_FPGA_DEV_NAME");
	if (strenv == NULL) {
		printf("GLUEENV_FPGA_DEV_NAME is NULL\n");
		exit(0);
	}
	printf("GLUEENV_FPGA_DEV_NAME=%s\n",strenv);
	fpga_dev = strenv;

	//Acquisition of environment variable information
	strenv = getenv("GLUEENV_FPGA_DMA_DEV_TO_HOST_CONNECTOR_ID");
	if (strenv == NULL) {
		printf("GLUEENV_FPGA_DMA_DEV_TO_HOST_CONNECTOR_ID is NULL\n");
		exit(0);
	}
	printf("GLUEENV_FPGA_DMA_DEV_TO_HOST_CONNECTOR_ID=%s\n",strenv);
	connector_id = strenv;
	



	// initialize DPDK
	//for debug
	if (secondary_mode != NULL && strcmp(secondary_mode, "0") == 0) {
		printf("start shmem primary mode\n");
		//log_libfpga(LIBFPGA_LOG_INFO,"start shmem primary mode\n");
		shmem_secondary = false;
	}	

	printf("start init_mem\n");
	ret = init_mem(file_prefix, shmem_secondary);
	if (ret < 0) {
		printf("init_mem failed ret=%d\n", ret);
		//log_libfpga(LIBFPGA_LOG_ERROR, "init_mem failed ret=%d\n", ret);
		return -1;
	}


	// initialize FPGA
	printf("fpga_dev_init called!\n");
	ret = fpga_dev_init(fpga_dev, &dev_id);
	if (ret < 0) {
		printf("fpga_dev_init error!\n");
		return -1;
	}

	//Information Set
	glue_setting();

	// set dev_id list
	ret = set_dev_id_list();
	if (ret < 0) {
		goto _END2;
	}

	// lock FPGA
	for (size_t i=0; i < fpga_get_num(); i++ ) {
		uint32_t *dev_id = get_dev_id(i);
		ret = fpga_ref_acquire(*dev_id);
		if (ret < 0) {
			printf("dev(%u) fpga_ref_acquire:ret(%d) error!!\n", *dev_id, ret);
			//logfile(LOG_ERROR, "dev(%u) fpga_ref_acquire:ret(%d) error!!\n", *dev_id, ret);
			//rslt2file("dev(%u) fpga_ref_acquire error!!\n", *dev_id, ret);
			goto _END3;
		}
		printf("dev(%u) fpga_ref_acquire:ret(%d)\n", *dev_id, ret);
		//logfile(LOG_DEBUG, "dev(%u) fpga_ref_acquire:ret(%d)\n", *dev_id, ret);
	}

	// fpga lldma setup buffer
	for (size_t i=0; i < fpga_get_num(); i++ ) {
		uint32_t *dev_id = get_dev_id(i);
		ret = fpga_lldma_setup_buffer(*dev_id);
		if (ret < 0) {
			printf("dev(%u) fpga_lldma_setup_buffer:ret(%d) error!!\n", *dev_id, ret);
			//logfile(LOG_ERROR, "dev(%u) fpga_lldma_setup_buffer:ret(%d) error!!\n", *dev_id, ret);
			//rslt2file("dev(%u) fpga_lldma_setup_buffer error!!\n", *dev_id);
			goto _END3;
		}
		printf("dev(%u) fpga_lldma_setup_buffer:ret(%d)\n", *dev_id, ret);
		//logfile(LOG_DEBUG, "dev(%u) fpga_lldma_setup_buffer:ret(%d)\n", *dev_id, ret);
	}

	// fpga enable regrw
	for (size_t i=0; i < fpga_get_num(); i++ ) {
		uint32_t *dev_id = get_dev_id(i);
		ret = fpga_enable_regrw(*dev_id);
		if (ret < 0) {
			printf("dev(%u) fpga_enable_regrw:ret(%d) error!!\n", *dev_id, ret);
			//logfile(LOG_ERROR, "dev(%u) fpga_enable_regrw:ret(%d) error!!\n", *dev_id, ret);
			//rslt2file("dev(%u) fpga_enable_regrw error!!\n", *dev_id);
			goto _END3;
		}
		printf("dev(%u) fpga_enable_regrw:ret(%d)\n", *dev_id, ret);
		//logfile(LOG_DEBUG, "dev(%u) fpga_enable_regrw:ret(%d)\n", *dev_id, ret);
	}


	// execute Glue
	//printf("//--- GLUE START ---\n");
	//rslt2file("//--- GLUE START ---\n");

	ret = glue(tcp_client_info,connector_id,input_width,input_height);
	if (ret < 0) {
		printf("glue error(%d)\n", ret);
		//logfile(LOG_ERROR, "glue error(%d)\n", ret);
		//rslt2file("glue error(%d)\n", ret);
	}

	printf("//--- GLUE END ---//\n");
	//rslt2file("//--- GLUE END ---//\n");

_END3:
	// unlock FPGA
	for (size_t i=0; i < fpga_get_num(); i++ ) {
		uint32_t *dev_id = get_dev_id(i);
		ret = fpga_ref_release(*dev_id);
		if (ret < 0) {
			printf("dev(%u) fpga_ref_release:ret(%d) error!!\n", *dev_id, ret);
			//logfile(LOG_ERROR, "dev(%u) fpga_ref_release:ret(%d) error!!\n", *dev_id, ret);
			//rslt2file("dev(%u) fpga_ref_release error!!\n", *dev_id);
		}
		printf("dev(%u) fpga_ref_release:ret(%d)\n", *dev_id, ret);
		//logfile(LOG_DEBUG, "dev(%u) fpga_ref_release:ret(%d)\n", *dev_id, ret);
	}

_END2:
	// finish FPGA
	ret = fpga_finish();
	if (ret < 0) {
		printf("fpga finish error!!\n");
		//logfile(LOG_ERROR, "fpga finish error!!\n");
		//rslt2file("fpga finish error!!\n");
	}
	printf("fpga_finish:ret(%d)\n",ret);
	//logfile(LOG_DEBUG, "fpga_finish:ret(%d)\n",ret);

_END1:
	// finish DPDK shmem
	ret = fpga_shmem_finish();
	if (ret < 0) {
		printf("fpga shmem finish error!!\n");
		//logfile(LOG_ERROR, "fpga shmem finish error!!\n");
		//rslt2file("fpga shmem finish error!!\n");
	}
	printf("fpga_shmem_finish:ret(%d)\n",ret);
	//logfile(LOG_DEBUG, "fpga_shmem_finish:ret(%d)\n",ret);


	return 0;
}
