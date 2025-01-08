/*****************************************************************
 * Copyright 2024 NTT Corporation, FUJITSU LIMITED
 *****************************************************************/

#include <assert.h>
#include <stdio.h>
#include <string.h>

#include <rte_common.h>
#include <rte_eal.h>

#include "libchain.h"
#include "libfunction.h"
#include "liblldma.h"
#include "libfpgactl.h"
#include "liblogging.h"
#include "libshmem.h"
#include "libshmem_controller.h"

#include "lib_controller.h"

#define CH_NUM_MAX 16
#define LANE_NUM_MAX 2
#define FRAMEWORK_KRNL_NUM_MAX 2
#define FRAMEWORK_SUB_KRNL_NUM_MAX 2
#define FUNCTION_KRNL_NUM_MAX 2
#define CID_BASE 512

//--------------------------
// Constants for FPGA Kernel Configuration
//--------------------------

static const uint32_t fw_idx = FRAMEWORK_KRNL_NUM_MAX / LANE_NUM_MAX;
static const uint32_t fwsub_idx = FRAMEWORK_SUB_KRNL_NUM_MAX / LANE_NUM_MAX;
static const uint32_t func_idx = FUNCTION_KRNL_NUM_MAX / LANE_NUM_MAX;

static uint32_t dev_id;
struct data_size ds;
static dma_info_t dmainfo_rx;
static dma_info_t dmainfo_tx;

int init_DPDK(char* file_prefix) {
  int32_t ret = 0;
  ret = fpga_shmem_controller_init(0, NULL);
  if (ret < 0) {
    rte_exit(EXIT_FAILURE, "fpga_shmem_controller_init failed: ret=%d\n", ret);
  }
  ret = fpga_shmem_enable(file_prefix, NULL);
  if (ret < 0) {
    rte_exit(EXIT_FAILURE, "fpga_shmem_enable failed: ret=%d\n", ret);
  }
  return 0;
}

int fini_DPDK(char *file_prefix) {
  int ret;
  ret = fpga_shmem_disable(file_prefix);
  if (ret < 0) {
    rte_exit(EXIT_FAILURE, "fpga_shmem_disable failed: ret=%d\n", ret);
  }
  ret = fpga_shmem_controller_finish();
  if (ret < 0) {
    printf("fpga_shmem_controller_finish failed: ret=%d\n", ret);
  }
  return 0;
}

static int32_t set_frame_size(uint32_t dev_id,
			      uint32_t lane_id,
			      uint32_t input_height,
			      uint32_t input_width,
			      uint32_t output_height,
			      uint32_t output_width) {
  int32_t ret = 0;
  uint32_t fwsub_krnl_pid = lane_id * fwsub_idx;
  uint32_t func_krnl_pid = lane_id * func_idx;
  ret = fpga_function_config(dev_id, lane_id, "filter_resize");
  if (ret < 0) {
    printf("fpga_function_config failed: ret=%d\n", ret);
    return -1;
  }
  ret = fpga_function_init(dev_id, lane_id, NULL);
  if (ret < 0) {
    printf("fpga_function_init failed: ret=%d\n", ret);
    return -1;
  }
  char json_txt[256];
  snprintf(json_txt, 256, "{\"i_width\":%u, \"i_height\":%u, \"o_width\":%u,\"o_height\":%u}", input_width, input_height, output_width, output_height);
  ret = fpga_function_set(dev_id, lane_id, json_txt);
  if (ret < 0) {
    printf("fpga_function_set failed: ret=%d\n", ret);
    return -1;
  }
}

static int chain_connect(uint32_t dev_id, int ch_id) {
  int ret;
  uint32_t ch_div_unit = CH_NUM_MAX / FUNCTION_KRNL_NUM_MAX;
  uint32_t func_krnl_pid = ch_id / ch_div_unit;
  uint32_t fchid = ch_id;
  bool ingress_extif_id = false;
  uint32_t ingress_cid = ch_div_unit * func_krnl_pid + ch_id;
  bool egress_extif_id = false;
  uint32_t egress_cid = ch_div_unit * func_krnl_pid + ch_id;
  bool ingress_active_flag = true;
  bool egress_active_flag = true;
  bool direct_flag = false;
  bool egress_virtual_flag = false;
  bool egress_blocking_flag = true;
  ret = fpga_chain_connect(dev_id, func_krnl_pid, fchid, ingress_extif_id, ingress_cid, egress_extif_id, egress_cid, ingress_active_flag, egress_active_flag, direct_flag, egress_virtual_flag, egress_blocking_flag);
  if (ret < 0) {
    printf("fpga_chain_connect failed: ret=%d\n", ret);
    return -1;
  }
  bool extif_id = false;
  ret = fpga_chain_set_ddr(dev_id, func_krnl_pid, extif_id);
  if (ret < 0) {
    printf("fpga_chain_set_ddr failed: ret=%d\n", ret);
    return -1;
  }
  ret = fpga_chain_start(dev_id, func_krnl_pid);
  if (ret < 0) {
    printf("fpga_chain_start failed: ret=%d\n", ret);
    return -1;
  }
  return 0;
}

int init_FPGA(char* fpga_dev, int ch_id) {
  int ret = 0;
  ret = fpga_dev_init(fpga_dev, &dev_id);
  if (ret < 0) {
    printf("fpga_dev_init failed: ret=%d\n", ret);
    return -1;
  }
  ret = fpga_enable_regrw(dev_id);
  if (ret < 0) {
    printf("fpga_enable_regrw failed: ret=%d\n", ret);
    return -1;
  }
  int fpga_num = fpga_get_num();
  if (fpga_num != 1) {
    printf("Num of FPGA error(%d)\n", fpga_num);
    return -1;
  }
  uint32_t lane_id = ch_id / (CH_NUM_MAX / LANE_NUM_MAX);
  ret = set_frame_size(dev_id, lane_id, ds.input_width, ds.input_height, ds.output_width, ds.output_height);
  if (ret) {
    printf("set_frame_size failed: ret=%d\n", ret);
    return -1;
  }
  memset(&dmainfo_rx, 0, sizeof(dma_info_t));
  assert(getenv("CONNECTOR_ID_RX") != NULL);
  char *connector_id_rx = getenv("CONNECTOR_ID_RX");
  ret = fpga_lldma_init(dev_id, DMA_HOST_TO_DEV, ch_id, connector_id_rx, &dmainfo_rx);
  if (ret < 0) {
    printf("fpga_lldma_init faled: ret=%d\n", ret);
    return -1;
  }
  memset(&dmainfo_tx, 0, sizeof(dma_info_t));
  assert(getenv("CONNECTOR_ID_TX") != NULL);
  char *connector_id_tx = getenv("CONNECTOR_ID_TX");
  ret = fpga_lldma_init(dev_id, DMA_DEV_TO_HOST, ch_id, connector_id_tx, &dmainfo_tx);
  if (ret < 0) {
    printf("fpga_lldma_init failed: ret=%d\n", ret);
    return -1;
  }
  ret = chain_connect(dev_id, ch_id);
  if (ret < 0) {
    printf("chain_connect failed: ret=%d\n", ret);
    return -1;
  }
  return 0;
}

int fini_FPGA(int ch_id) {
  int ret;
  uint32_t lane_id = ch_id / (CH_NUM_MAX / LANE_NUM_MAX);
  ret = fpga_chain_stop(dev_id, lane_id);
  if (ret < 0) {
    printf("fpga_chain_stop failed: ret=%d\n", ret);
    return -1;
  }
  ret = fpga_function_finish(dev_id, lane_id, NULL);
  if (ret < 0) {
    printf("fpga_function_finish: ret=%d\n", ret);
  }
  ret = fpga_lldma_finish(&dmainfo_rx);
  if (ret < 0) {
    printf("fpga_lldma_finish: ret=%d\n", ret);
  }
  ret = fpga_lldma_finish(&dmainfo_tx);
  if (ret < 0) {
    printf("fpga_lldma_finish failed: ret=%d\n", ret);
  }
  ret = fpga_disable_regrw(dev_id);
  if (ret < 0) {
    printf("fpga_disable_regrw failed: ret=%d\n", ret);
  }
  ret = fpga_finish();
  if (ret < 0) {
    printf("fpga_finish failed: ret=%d\n", ret);
  }
  return 0;
}
