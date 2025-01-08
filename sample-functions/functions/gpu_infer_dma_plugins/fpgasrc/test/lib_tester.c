/*****************************************************************
 * Copyright 2024 NTT Corporation, FUJITSU LIMITED
 *****************************************************************/

#include <assert.h>
#include <unistd.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/time.h>
#include <pthread.h>
#include <netinet/in.h>
#include <rte_eal.h>
#include "libshmem.h"
#include "libfpgactl.h"
#include "libdma.h"
#include "libdmacommon.h"
#include "liblldma.h"
#include "libptu.h"
#include "libchain.h"
#include "common.h"
#include "param_tables.h"
#include "bcdbg.h"

#include "lib_tester.h"

#define GST_PLUGIN

#define WAIT_TIME_DEQUEUE       5000    //msec

pthread_mutex_t shmmutex[CH_NUM_MAX][SHMEMALLOC_NUM_MAX];

static int32_t wait_fpga_dequeue(dma_info_t *dmainfo, dmacmd_info_t *dmacmdinfo, const uint32_t deq_id, const int32_t msec) {
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
        logfile(LOG_ERROR, "  CH(%u) deq(%u) task_id(%u) dequeue timeout!!!\n", ch_id, deq_id, task_id);
        return -1;
}

#ifdef USE_DIRECT_TRANSFER
static int32_t wait_fpga_dequeue_physical(dma_info_t *dmainfo, dmacmd_info_t *dmacmdinfo, const uint32_t deq_id, const int32_t msec) {
        int32_t timeout = 100; //fpga_deuque timeout 100msec
        uint32_t ch_id = dmainfo->chid;
        uint16_t task_id = dmacmdinfo->task_id;
        int32_t cnt = msec/timeout;
        for (size_t i=0; i < cnt; i++) {
	  int32_t ret = fpga_dequeue_physical(dmainfo, dmacmdinfo);
                if (ret == 0) {
                        return 0;
                }
        }
        logfile(LOG_ERROR, "  CH(%u) deq(%u) task_id(%u) dequeue timeout!!!\n", ch_id, deq_id, task_id);
        return -1;
}
#endif /* USE_DIRECT_TRANSFER */

#define MARKER 0xE0FF10AD

static void end1();
static void end2();
static void end3();
static void end_mem(uint32_t ch_id);
static void end_ptu_lldma(int ch_id);
static void end_queue(uint32_t ch_id);

struct data_size ds;
static uint32_t lane_num = 1;
static uint32_t lane_id = 0;
static mngque_t pque[CH_NUM_MAX];
static int debug_read_counter = 0;
static pthread_t thread_ptu_id[CH_NUM_MAX];
static uint32_t tcp_cid[CH_NUM_MAX];

unsigned long get_current_time() {
  struct timeval tv;
  int rc = gettimeofday(&tv, NULL);
  assert(rc == 0);
  unsigned long current_time = 1UL*1000*1000*1000*tv.tv_sec + 1UL*1000*tv.tv_usec;
  return current_time;
}

static int parse_app_args(int argc, char **argv)
{
  int32_t ret = 0;

  ret = parse_app_args_func(argc, argv);
  if (ret < 0) {
    return ret;
  }

  ret = check_options();
  if (ret < 0) {
    return ret;
  }

  return 0;
}

static void end1() {

  uint32_t ch_num = getopt_ch_num();

  // free shmem buffer
  log_libfpga(LIBFPGA_LOG_INFO, "--- deq_shmem_free ---\n");
  for (size_t i=0; i < ch_num; i++) {
    deq_shmem_free(&pque[i], i);
  }

  // free dmacmdinfo buffer
  log_libfpga(LIBFPGA_LOG_INFO, "--- dmacmdinfo_free ---\n");
  dmacmdinfo_free();

  log_libfpga(LIBFPGA_LOG_INFO, "...test end\n");
}

static void end2() {

  int32_t ret = 0;
  uint32_t ch_num = getopt_ch_num();
  uint32_t dev_id = 0; // kari
  uint32_t ptu_idx = PTU_KRNL_NUM_MAX / LANE_NUM_MAX;

  // ptu exit
  log_libfpga(LIBFPGA_LOG_INFO, "--- fpga_ptu_exit ---\n");
  for (size_t i=0; i < lane_num; i++) {
    uint32_t ptu_krnl_pid = i * ptu_idx;
    log_libfpga(LIBFPGA_LOG_INFO, "Lane(%zu) fpga_ptu_exit\n", i);
    ret = fpga_ptu_exit(dev_id, ptu_krnl_pid);
    if (ret < 0) {
      log_libfpga(LIBFPGA_LOG_INFO, "fpga_ptu_exit error!!!(%d)\n",ret);
      //error
    }
  }

  // lldma finish
  log_libfpga(LIBFPGA_LOG_INFO, "--- fpga_lldma_finish ---\n");
  for (size_t i=0; i < ch_num; i++) {
    log_libfpga(LIBFPGA_LOG_INFO, "CH(%zu) fpga_dma_finish\n", i);
    dma_info_t *pdmainfo_ch = get_dmainfo_channel(i);
    ret = fpga_lldma_finish(pdmainfo_ch);
    if (ret < 0) {
      log_libfpga(LIBFPGA_LOG_ERROR, "fpga_dma_finish error!!!(%d)\n",ret);
      // error
    }
    pdmainfo_ch->connector_id = NULL;
    prlog_dma_info(pdmainfo_ch, i);
  }
  end1();
}

static void end3() {

  int32_t ret = 0;
  uint32_t ch_num = getopt_ch_num();

  //----------------------------------------------
  // end processing
  //----------------------------------------------
  // lldma queue finish
  log_libfpga(LIBFPGA_LOG_INFO, "--- fpga_lldma_queue_finish ---\n");
  for (size_t i=0; i < ch_num; i++) {
    log_libfpga(LIBFPGA_LOG_INFO, "CH(%zu) fpga_lldma_queue_finish\n", i);
    dma_info_t *pdmainfo = get_dmainfo(i);
    ret = fpga_lldma_queue_finish(pdmainfo);
    if (ret < 0) {
      log_libfpga(LIBFPGA_LOG_ERROR, "fpga_lldma_queue_finish error!!!(%d)\n",ret);
      // error
    }
    pdmainfo->connector_id = NULL;
    prlog_dma_info(pdmainfo, i);
  }
  end2();
}

#ifdef USE_DIRECT_TRANSFER

int read_buffer_physical(uint32_t cmd_idx) {
  int dummy_chid = 0;
  int32_t ret;
  dma_info_t *pdmainfo = get_dmainfo(dummy_chid);
  dmacmd_info_t *pdmacmdinfo = get_dmacmdinfo(dummy_chid, cmd_idx);
  ret = wait_fpga_dequeue_physical(pdmainfo, pdmacmdinfo, cmd_idx, WAIT_TIME_DEQUEUE);
  if (ret < 0) {
    // timeout
    log_libfpga(LIBFPGA_LOG_WARN, "read_buffer timed out!\n");
    return -1;
  }
  uint32_t next_cmd_idx = cmd_idx + 1;
  uint32_t deq_num = getopt_deq_num();
  return next_cmd_idx % deq_num;
}

#endif /* USE_DIRECT_TRANSFER */

int write_buffer(uint32_t cmd_idx) {
  int dummy_chid = 0;
  int32_t ret;
  dma_info_t *pdmainfo = get_dmainfo(dummy_chid);
  dmacmd_info_t *pdmacmdinfo = get_dmacmdinfo(dummy_chid, cmd_idx);
  frameheader_t *fh = (frameheader_t*)((uint64_t)pdmacmdinfo->data_addr);
  uint64_t image_size = ds.input_height * ds.input_width * 3;
  fh->payload_len = image_size;
  ret = fpga_enqueue(pdmainfo, pdmacmdinfo);
  if (ret < 0) {
    log_libfpga(LIBFPGA_LOG_ERROR, "fpga_enqueue error!!!(%d)\n",ret);
    // error
    end3();
    return -1;
  }
  uint32_t next_cmd_idx = cmd_idx + 1;
  uint32_t deq_num = getopt_deq_num();
  return next_cmd_idx % deq_num;
}

int free_buffer(uint32_t cmd_idx) {
  int dummy_chid = 0;
  int32_t ret;
  dma_info_t *pdmainfo = get_dmainfo(dummy_chid);
  dmacmd_info_t *pdmacmdinfo = get_dmacmdinfo(dummy_chid, cmd_idx);
  ret = wait_fpga_dequeue(pdmainfo, pdmacmdinfo, cmd_idx, WAIT_TIME_DEQUEUE);
  if (ret < 0) {
    // timeout
    log_libfpga(LIBFPGA_LOG_WARN, "read_buffer timed out!\n");
    return -1;
  }
  return 0;
}

static void end_mem(uint32_t ch_id) {

  // free shmem buffer
  log_libfpga(LIBFPGA_LOG_INFO, "--- deq_shmem_free ---\n");
  deq_shmem_free(&pque[ch_id], ch_id);

  // free dmacmdinfo buffer
  log_libfpga(LIBFPGA_LOG_INFO, "--- dmacmdinfo_free ---\n");
  dmacmdinfo_free();

  log_libfpga(LIBFPGA_LOG_INFO, "...test end\n");
}

static void end_ptu_lldma(int ch_id) {

  int32_t ret = 0;
  uint32_t ch_num;
  if (ch_id == -1) {
    ch_id = 0;
    ch_num = getopt_ch_num();
  } else {
    ch_num = 1;
  }
  uint32_t dev_id = 0; // kari
  uint32_t ptu_idx = PTU_KRNL_NUM_MAX / LANE_NUM_MAX;

  // ptu exit
  log_libfpga(LIBFPGA_LOG_INFO, "--- fpga_ptu_exit ---\n");
  for (size_t i=lane_id; i < lane_id + lane_num; i++) {
    uint32_t ptu_krnl_pid = i * ptu_idx;
    log_libfpga(LIBFPGA_LOG_INFO, "Lane(%zu) fpga_ptu_exit\n", i);
    ret = fpga_ptu_exit(dev_id, ptu_krnl_pid);
    if (ret < 0) {
      log_libfpga(LIBFPGA_LOG_ERROR, "fpga_ptu_exit error!!!(%d)\n",ret);
      //error
    }
  }

  // lldma finish
  log_libfpga(LIBFPGA_LOG_INFO, "--- fpga_lldma_finish ---\n");
  for (size_t i=ch_id; i < ch_id + ch_num; i++) {
    log_libfpga(LIBFPGA_LOG_INFO, "CH(%zu) fpga_dma_finish\n", i);
    dma_info_t *pdmainfo_ch = get_dmainfo_channel(i);
    ret = fpga_lldma_finish(pdmainfo_ch);
    if (ret < 0) {
      log_libfpga(LIBFPGA_LOG_ERROR, "fpga_dma_finish error!!!(%d)\n",ret);
      // error
    }
    pdmainfo_ch->connector_id = NULL;
    prlog_dma_info(pdmainfo_ch, i);
  }
}

static void end_queue(uint32_t ch_id) {

  int32_t ret = 0;

  //----------------------------------------------
  // end processing
  //----------------------------------------------
  // lldma queue finish
  log_libfpga(LIBFPGA_LOG_INFO, "--- fpga_lldma_queue_finish ---\n");
  log_libfpga(LIBFPGA_LOG_INFO, "CH(%zu) fpga_lldma_queue_finish\n", ch_id);
  dma_info_t *pdmainfo = get_dmainfo(ch_id);
  ret = fpga_lldma_queue_finish(pdmainfo);
  if (ret < 0) {
    log_libfpga(LIBFPGA_LOG_ERROR, "fpga_lldma_queue_finish error!!!(%d)\n",ret);
    // error
  }
  log_libfpga(LIBFPGA_LOG_INFO, "fpga_lldma_queue_finish called! chid=%d\n", ch_id);
  pdmainfo->connector_id = NULL;
  prlog_dma_info(pdmainfo, ch_id);
  //end_ptu_lldma(config_fpga, ch_id);
}

int init_tester(int argc, char **argv, char* connector_id, char* fpga_device, char* file_prefix) {

  int ret;
  int dummy_chid = 0;
  uint32_t dev_id = 0; // kari

  //------------------
  // arguments check
  //------------------
  set_cmdname(argv[0]);

  if (argc == 1) {
    print_usage();
    return -1;
  }

  // initialize DPDK
#if 1
  ret = fpga_shmem_init(file_prefix, NULL, 0);
#else // debug
  ret = fpga_shmem_init_sys(file_prefix, NULL, NULL, NULL, 0);
#endif
  if (ret < 0) {
    rte_exit(EXIT_FAILURE, "Initialize failed\n");
  }
  log_libfpga(LIBFPGA_LOG_INFO, "fpga_shmem_init:ret(%d)\n",ret);

  // initialize FPGA
  ret = fpga_dev_init(fpga_device);
  printf("fpga_init called!\n");
  if (ret < 0) {
    log_libfpga(LIBFPGA_LOG_ERROR, "fpga init error!! ret=%d\n", ret);
    return -1;
  }
  log_libfpga(LIBFPGA_LOG_INFO, "fpga_init:ret(%d)\n",ret);
  argc -= ret;
  argv += ret;

  ret = parse_app_args(argc, argv);
  if (ret < 0) {
    log_libfpga(LIBFPGA_LOG_ERROR, "app option error!!\n");
    return -1;
  }
  log_libfpga(LIBFPGA_LOG_INFO, "parse_app_options:ret(%d)\n",ret);

  log_libfpga(LIBFPGA_LOG_INFO, "//--- TEST START ---\n");

  int32_t errcnt = 0;
  memset(&pque, 0, sizeof(mngque_t) * CH_NUM_MAX);
  char cmd[128];

  int32_t fpga_num = fpga_get_num();
  if (fpga_num != 1) {
    log_libfpga(LIBFPGA_LOG_ERROR, " Num of FPGA error(%d)\n", fpga_num);
    return -1;
  }

  // --- get options  ---
  log_libfpga(LIBFPGA_LOG_INFO, "--- get options  ---\n");
  uint32_t ch_num = getopt_ch_num();
  //uint32_t fps = getopt_fps();
  //uint32_t frame_num = getopt_frame_num();
  uint32_t deq_num = getopt_deq_num();

  if (ch_num > (CH_NUM_MAX / LANE_NUM_MAX)) {
    lane_num = 2;
  }

  uint32_t ptu_idx = PTU_KRNL_NUM_MAX / LANE_NUM_MAX;
  uint32_t fw_idx = FRAMEWORK_KRNL_NUM_MAX / LANE_NUM_MAX;
  uint32_t fwsub_idx = FRAMEWORK_SUB_KRNL_NUM_MAX / LANE_NUM_MAX;
  uint32_t func_idx = FUNCTION_KRNL_NUM_MAX / LANE_NUM_MAX;

  // deq_shmstate init
  bool *s = get_deq_shmstate(dummy_chid);
  for (size_t j=0; j < SHMEMALLOC_NUM_MAX; j++) {
    s[j] = false;
    pthread_mutex_init( &shmmutex[dummy_chid][j], NULL );
  }

  // allocate dmacmdinfo buffer
  log_libfpga(LIBFPGA_LOG_INFO, "--- dmacmdinfo_malloc ---\n");
  ret = dmacmdinfo_malloc();
  if (ret < 0) {
    log_libfpga(LIBFPGA_LOG_ERROR, "dmacmdinfo_alloc error!!!(%d)\n",ret);
    return -1;
  }

  //----------------------------------------------
  // allocate data buffer
  //----------------------------------------------
  log_libfpga(LIBFPGA_LOG_INFO, "--- deq_shmem_malloc ---\n");
  ret = deq_shmem_malloc(&pque[dummy_chid], dummy_chid);
  if (ret < 0) {
    log_libfpga(LIBFPGA_LOG_ERROR, "deq_shmem_malloc error(%d)\n",ret);
    return -1;
  }
  prlog_mngque(&pque[dummy_chid]);

#if 0 // for debug

  // set frame size
  log_libfpga(LIBFPGA_LOG_INFO, "--- fpga_set_frame_size ---\n");
  for (size_t i=0; i < lane_num; i++) {
    uint32_t fwsub_krnl_pid = i * fwsub_idx;
    uint32_t func_krnl_pid = i * func_idx;
    uint32_t fwsub_pid_ch_idx = i * (CH_NUM_MAX / FRAMEWORK_SUB_KRNL_NUM_MAX);
    log_libfpga(LIBFPGA_LOG_INFO, "Lane(%zu) fpga_set_frame_size [FUNCTION_%u]\n", i, func_krnl_pid);
    ret = fpga_set_frame_size(dev_id, FUNC, func_krnl_pid, ds.output_width, ds.output_height, ds.output_width, ds.output_height);
    printf("fpga_set_frame_size: (%d, %d) -> (%d, %d)\n", ds.output_width, ds.output_height, ds.output_width, ds.output_height);
    if (ret) {
      log_libfpga(LIBFPGA_LOG_ERROR, "fpga_set_frame_size error!!!(%d)\n",ret);
      //error
      //end_mem();
      return -1;
    }
    log_libfpga(LIBFPGA_LOG_INFO, "Lane(%zu) fpga_set_frame_size [FRAMESUB_%u]\n", i, fwsub_krnl_pid);
    ret = fpga_set_frame_size(dev_id, SUB, fwsub_krnl_pid, ds.output_width, ds.output_height, ds.output_width, ds.output_height);
    if (ret) {
      log_libfpga(LIBFPGA_LOG_ERROR, "fpga_set_frame_size error!!!(%d)\n",ret);
      //error
      //end_mem();
      return -1;
    }
  }

  // start module
  log_libfpga(LIBFPGA_LOG_INFO, "--- fpga_start_module ---\n");
  for (size_t i=0; i < lane_num; i++) {
    uint32_t fw_krnl_pid = i * fw_idx;
    uint32_t fwsub_krnl_pid = i * fwsub_idx;
    uint32_t func_krnl_pid = i * func_idx;

    log_libfpga(LIBFPGA_LOG_INFO, "Lane(%zu) fpga_start_module [FRAME_%u]\n", i, fw_krnl_pid);
    ret = fpga_start_module(dev_id, FRAME, fw_krnl_pid);
    if (ret) {
      log_libfpga(LIBFPGA_LOG_ERROR, "fpga_start_module error!!!(%d)\n",ret);
      //error
      //end_mem();
      return -1;
    }
    log_libfpga(LIBFPGA_LOG_INFO, "Lane(%zu) fpga_start_module [FRAMESUB_%u]\n", i, fwsub_krnl_pid);
    ret = fpga_start_module(dev_id, SUB, fwsub_krnl_pid);
    if (ret) {
      log_libfpga(LIBFPGA_LOG_ERROR, "fpga_start_module error!!!(%d)\n",ret);
      //error
      //end_mem();
      return -1;
    }
    log_libfpga(LIBFPGA_LOG_INFO, "Lane(%zu) fpga_start_module [FUNCTION_%u]\n", i, func_krnl_pid);
    ret = fpga_start_module(dev_id, FUNC, func_krnl_pid);
    if (ret) {
      log_libfpga(LIBFPGA_LOG_ERROR, "fpga_start_module error!!!(%d)\n",ret);
      //error
      //end_mem();
      return -1;
    }
  }

  //----------------------------------------------
  // fpga lldma init
  //----------------------------------------------
  log_libfpga(LIBFPGA_LOG_INFO, "--- fpga_lldma_init ---\n");
  {
    int i = dummy_chid;
    dma_info_t *pdmainfo_ch = get_dmainfo_channel(i);
    memset(pdmainfo_ch, 0, sizeof(dma_info_t));
    char *connector_id = (char*)getparam_connector_id_rx(i);
    printf("fpga_lldma_init: chid=%d, connector_id=%s\n", i, connector_id);
    log_libfpga(LIBFPGA_LOG_INFO, "CH(%d) fpga_lldma_init\n", i);
    ret = fpga_lldma_init(dev_id, DMA_DEV_TO_HOST, i, connector_id, pdmainfo_ch);
    if (ret < 0) {
      log_libfpga(LIBFPGA_LOG_ERROR, "fpga_lldma_init error!!!(%d)\n",ret);
      // error
      //end_mem();
      return -1;
    }
    prlog_dma_info(pdmainfo_ch, i);
  }
#endif

  //----------------------------------------------
  // fpga lldma queue setup (set dmainfo)
  //----------------------------------------------
  log_libfpga(LIBFPGA_LOG_INFO, "--- fpga_lldma_queue_setup ---\n");
  dma_info_t *pdmainfo = get_dmainfo(dummy_chid);
  memset(pdmainfo, 0, sizeof(dma_info_t));
  log_libfpga(LIBFPGA_LOG_INFO, "CH(%zu) fpga_lldma_queue_setup\n", dummy_chid);
  ret = fpga_lldma_queue_setup(connector_id, pdmainfo);
  if (ret < 0) {
    printf("fpga_lldma_queue_setup error!!!(%d)\n",ret);
    log_libfpga(LIBFPGA_LOG_ERROR, "fpga_lldma_queue_setup error!!!(%d)\n",ret);
    // error
    end_mem(dummy_chid);
    return -1;
  }
  printf("fpga_lldma_queue_setup called! connector_id=%s\n", connector_id);
  prlog_dma_info(pdmainfo, dummy_chid);

  fpga_queue_t *queue = pdmainfo->queue_addr;
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

  //----------------------------------------------
  // set dmacmd info
  //----------------------------------------------
  log_libfpga(LIBFPGA_LOG_INFO, "--- set_dma_cmd ---\n");
  uint32_t data_len = pque[dummy_chid].deq_dstlen;
  uint32_t dsize = pque[dummy_chid].dsize;
  uint32_t dstidx = 0;
  log_libfpga(LIBFPGA_LOG_INFO, "CH(%zu) data size=%d Byte\n", dummy_chid, dsize);
  for (size_t j=0; j < deq_num; j++) {
    uint32_t task_id = j + 1;
    uint32_t graph_id = 0;
    if (dstidx >= getopt_shmalloc_num()) {
      dstidx = 0;
    }
    void *data_addr = pque[dummy_chid].dstbufp[dstidx];
    dstidx++;
    dmacmd_info_t *pdmacmdinfo = get_dmacmdinfo(dummy_chid, j);
    memset(pdmacmdinfo, 0, sizeof(dmacmd_info_t));
    log_libfpga(LIBFPGA_LOG_INFO, "CH(%zu) DEQ(%zu) set_dma_cmd\n", dummy_chid, j);
    ret = set_dma_cmd(pdmacmdinfo, task_id, data_addr, data_len);
    if (ret < 0) {
      log_libfpga(LIBFPGA_LOG_ERROR, "set_dma_cmd error!!!(%d)\n",ret);
      // error
      end_queue(dummy_chid);
      end_mem(dummy_chid);
      return -1;
    }
    prlog_dmacmd_info(pdmacmdinfo, dummy_chid, j);
  }

  //----------------------------------------------
  // fpga enqueue
  //----------------------------------------------
  log_libfpga(LIBFPGA_LOG_INFO, "--- fpga_enqueue ---\n");
  //dma_info_t *pdmainfo = get_dmainfo(dummy_chid);
  for (size_t j=0; j < deq_num; j++) {
    dmacmd_info_t *pdmacmdinfo = get_dmacmdinfo(dummy_chid, j);
    frameheader_t *fh = (frameheader_t*)((uint64_t)pdmacmdinfo->data_addr);
    fh->payload_len = 0;
    ret = fpga_enqueue(pdmainfo, pdmacmdinfo);
    if (ret < 0) {
      printf("fpga_enqueue error!!!(%d)\n",ret);
      log_libfpga(LIBFPGA_LOG_ERROR, "fpga_enqueue error!!!(%d)\n",ret);
      // error
      end_queue(dummy_chid);
      end_mem(dummy_chid);
      return -1;
    }
  }

#if 0 // for debug
  //----------------------------------------------
  // function chain control
  //----------------------------------------------
  resl2file("\n--- function chain connect ---\n");
  log_libfpga(LIBFPGA_LOG_INFO, "--- fpga_chain_connect ---\n");
  uint32_t kernel_num = get_function_krnl_id(dummy_chid);
  uint32_t fchid = getparam_function_chid(dummy_chid);
  // HACK: assume ingress cid = chid + 1 (work-around for libptu without multi-process support)
  uint32_t ingress_cid = dummy_chid + 1;
  uint32_t egress_cid = getparam_lldma_cid(dummy_chid);
  printf("CH(%d) fpga_chain_connect\n", dummy_chid);
  printf("  kernel_num(%u), fchid(%u) ingress_cid(%u) egress_cid(%u)\n", kernel_num, fchid, ingress_cid, egress_cid);
  log_libfpga(LIBFPGA_LOG_INFO, "CH(%u) func_kernel_id(%u), fchid(%u) ingress_cid(%u) egress_cid(%u)\n", dummy_chid, kernel_num, fchid, ingress_cid, egress_cid);
  ret = fpga_chain_connect(dev_id, kernel_num, fchid, ingress_cid, egress_cid);
  if (ret < 0) {
    printf("fpga_chain_connect error!!!(%d)\n",ret);
    // error
    // end_queue(chid);
    return -1;
  }
#endif

  return 0;
}

#ifdef USE_DIRECT_TRANSFER

int init_tester_physical(int argc, char **argv, char* connector_id, unsigned long paddr) {

  int ret;
  int dummy_chid = 0;
  uint32_t dev_id = 0; // kari

  //------------------
  // arguments check
  //------------------
  set_cmdname(argv[0]);

  if (argc == 1) {
    print_usage();
    return -1;
  }

  // initialize DPDK
  ret = fpga_shmem_init(argc, argv);
  if (ret < 0) {
    rte_exit(EXIT_FAILURE, "Initialize failed\n");
  }
  log_libfpga(LIBFPGA_LOG_INFO, "fpga_shmem_init:ret(%d)\n",ret);

  // initialize FPGA
  ret = fpga_init(argc, argv);
  printf("fpga_init called!\n");
  if (ret < 0) {
    log_libfpga(LIBFPGA_LOG_ERROR, "fpga init error!! ret=%d\n", ret);
    return -1;
  }
  log_libfpga(LIBFPGA_LOG_INFO, "fpga_init:ret(%d)\n",ret);
  argc -= ret;
  argv += ret;

  ret = parse_app_args(argc, argv);
  if (ret < 0) {
    log_libfpga(LIBFPGA_LOG_ERROR, "app option error!!\n");
    return -1;
  }
  log_libfpga(LIBFPGA_LOG_INFO, "parse_app_options:ret(%d)\n",ret);

  log_libfpga(LIBFPGA_LOG_INFO, "//--- TEST START ---\n");

  int32_t errcnt = 0;
  memset(&pque, 0, sizeof(mngque_t) * CH_NUM_MAX);
  char cmd[128];

  int32_t fpga_num = fpga_get_num();
  if (fpga_num != 1) {
    log_libfpga(LIBFPGA_LOG_ERROR, " Num of FPGA error(%d)\n", fpga_num);
    return -1;
  }

  // --- get options  ---
  log_libfpga(LIBFPGA_LOG_INFO, "--- get options  ---\n");
  uint32_t ch_num = getopt_ch_num();
  //uint32_t fps = getopt_fps();
  //uint32_t frame_num = getopt_frame_num();
  uint32_t deq_num = getopt_deq_num();

  if (ch_num > (CH_NUM_MAX / LANE_NUM_MAX)) {
    lane_num = 2;
  }

  uint32_t ptu_idx = PTU_KRNL_NUM_MAX / LANE_NUM_MAX;
  uint32_t fw_idx = FRAMEWORK_KRNL_NUM_MAX / LANE_NUM_MAX;
  uint32_t fwsub_idx = FRAMEWORK_SUB_KRNL_NUM_MAX / LANE_NUM_MAX;
  uint32_t func_idx = FUNCTION_KRNL_NUM_MAX / LANE_NUM_MAX;

  // deq_shmstate init
  bool *s = get_deq_shmstate(dummy_chid);
  for (size_t j=0; j < SHMEMALLOC_NUM_MAX; j++) {
    s[j] = false;
    pthread_mutex_init( &shmmutex[dummy_chid][j], NULL );
  }

  // allocate dmacmdinfo buffer
  log_libfpga(LIBFPGA_LOG_INFO, "--- dmacmdinfo_malloc ---\n");
  ret = dmacmdinfo_malloc();
  if (ret < 0) {
    log_libfpga(LIBFPGA_LOG_ERROR, "dmacmdinfo_alloc error!!!(%d)\n",ret);
    return -1;
  }

  //----------------------------------------------
  // allocate data buffer
  //----------------------------------------------
  log_libfpga(LIBFPGA_LOG_INFO, "--- deq_shmem_malloc ---\n");
  ret = deq_shmem_malloc(&pque[dummy_chid], dummy_chid);
  if (ret < 0) {
    log_libfpga(LIBFPGA_LOG_ERROR, "deq_shmem_malloc error(%d)\n",ret);
    return -1;
  }
  prlog_mngque(&pque[dummy_chid]);

  //----------------------------------------------
  // fpga lldma queue setup (set dmainfo)
  //----------------------------------------------
  log_libfpga(LIBFPGA_LOG_INFO, "--- fpga_lldma_queue_setup ---\n");
  dma_info_t *pdmainfo = get_dmainfo(dummy_chid);
  memset(pdmainfo, 0, sizeof(dma_info_t));
  log_libfpga(LIBFPGA_LOG_INFO, "CH(%zu) fpga_lldma_queue_setup\n", dummy_chid);
  ret = fpga_lldma_queue_setup(connector_id, pdmainfo);
  if (ret < 0) {
    printf("fpga_lldma_queue_setup error!!!(%d)\n",ret);
    log_libfpga(LIBFPGA_LOG_ERROR, "fpga_lldma_queue_setup error!!!(%d)\n",ret);
    // error
    end_mem(dummy_chid);
    return -1;
  }
  printf("fpga_lldma_queue_setup called! connector_id=%s\n", connector_id);
  prlog_dma_info(pdmainfo, dummy_chid);

  fpga_queue_t *queue = pdmainfo->queue_addr;
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

  //----------------------------------------------
  // set dmacmd info
  //----------------------------------------------
  log_libfpga(LIBFPGA_LOG_INFO, "--- set_dma_cmd ---\n");
  uint32_t data_len = pque[dummy_chid].deq_dstlen;
  uint32_t dsize = pque[dummy_chid].dsize;
  uint32_t dstidx = 0;
  log_libfpga(LIBFPGA_LOG_INFO, "CH(%zu) data size=%d Byte\n", dummy_chid, dsize);
  for (size_t j=0; j < deq_num; j++) {
    uint32_t task_id = j + 1;
    uint32_t graph_id = 0;
    if (dstidx >= getopt_shmalloc_num()) {
      dstidx = 0;
    }
    void *data_addr = pque[dummy_chid].dstbufp[dstidx];
    dstidx++;
    dmacmd_info_t *pdmacmdinfo = get_dmacmdinfo(dummy_chid, j);
    memset(pdmacmdinfo, 0, sizeof(dmacmd_info_t));
    log_libfpga(LIBFPGA_LOG_INFO, "CH(%zu) DEQ(%zu) set_dma_cmd\n", dummy_chid, j);
    ret = set_dma_cmd(pdmacmdinfo, task_id, data_addr, data_len);
    if (ret < 0) {
      log_libfpga(LIBFPGA_LOG_ERROR, "set_dma_cmd error!!!(%d)\n",ret);
      // error
      end_queue(dummy_chid);
      end_mem(dummy_chid);
      return -1;
    }
    prlog_dmacmd_info(pdmacmdinfo, dummy_chid, j);
  }

  //----------------------------------------------
  // fpga enqueue
  //----------------------------------------------
  log_libfpga(LIBFPGA_LOG_INFO, "--- fpga_enqueue ---\n");
  //dma_info_t *pdmainfo = get_dmainfo(dummy_chid);
  for (size_t j=0; j < deq_num; j++) {
    dmacmd_info_t *pdmacmdinfo = get_dmacmdinfo(dummy_chid, j);
    frameheader_t *fh = (frameheader_t*)((uint64_t)pdmacmdinfo->data_addr);
    fh->payload_len = 0;
    ret = fpga_enqueue_physical(pdmainfo, pdmacmdinfo, paddr);
    if (ret < 0) {
      printf("fpga_enqueue error!!!(%d)\n",ret);
      log_libfpga(LIBFPGA_LOG_ERROR, "fpga_enqueue error!!!(%d)\n",ret);
      // error
      end_queue(dummy_chid);
      end_mem(dummy_chid);
      return -1;
    }
  }

  return 0;
}

#endif /* USE_DIRECT_TRANSFER */

int finish_tester() {

  int32_t ret = 0;
  int dummy_chid = 0;

#if 0 /* doesn't work */
  // clean up queue
  fpga_queue_t *queue = pdmainfo->queue_addr;
  queue->readhead = 0;
  queue->writehead = 0;
  fpga_desc_t *desc;
  for (int i = 0; i < queue->size; i++) {
    desc = &queue->ring[i];
    memset(desc, 0, sizeof(fpga_desc_t));
  }
#endif

#if 0
  // HACK
  int wait_time = 30;
  int count = 0;
  dma_info_t *pdmainfo = get_dmainfo(dummy_chid);
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

  end_queue(dummy_chid);
  end_mem(dummy_chid);
#endif

  log_libfpga(LIBFPGA_LOG_INFO, "//--- TEST END ---//\n");

  // finish DPDK shmem
  ret = fpga_shmem_finish();
  if (ret < 0) {
    log_libfpga(LIBFPGA_LOG_ERROR, "fpga shmem finish error!!\n");
  }
  log_libfpga(LIBFPGA_LOG_INFO, "fpga_shmem_finish:ret(%d)\n",ret);

  // FPGA close
  ret = fpga_finish();
  if (ret < 0) {
    log_libfpga(LIBFPGA_LOG_ERROR, "fpga finish error!!\n");
  }
  log_libfpga(LIBFPGA_LOG_INFO, "fpga_finish:ret(%d)\n",ret);

  return 0;
}
