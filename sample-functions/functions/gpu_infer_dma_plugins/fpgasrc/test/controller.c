/*****************************************************************
 * Copyright 2024 NTT Corporation, FUJITSU LIMITED
 *****************************************************************/

#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#include "lib_controller.h"
#include "liblogging.h"

static int ch_id;
static char* file_prefix;

void int_handler(int sig) {
  int ret;
  ret = fini_FPGA(ch_id);
  if (ret < 0) {
    exit(-1);
  }
  ret = fini_DPDK(file_prefix);
  if (ret < 0) {
    exit(-1);
  }
  printf("controller finished\n");
  fflush(stdout);
  exit(0);
}

int main() {

  int ret = 0;
  char* fpga_dev;
  extern struct data_size ds;
  char *strenv;

  int log_level;
  strenv = getenv("LOG_LEVEL");
  if (strenv == NULL) {
    log_level = LIBFPGA_LOG_INFO;
  } else {
    log_level = atoi(strenv);
  }
  libfpga_log_set_level(log_level);
  
  strenv = getenv("FPGA_DEV");
  if (strenv == NULL) {
    printf("environment variable FPGA_DEV must be specified\n");
    return -1;
  }
  fpga_dev = strenv;

  strenv = getenv("FILE_PREFIX");
  if (strenv == NULL) {
    printf("environment variable FILE_PREFIX must be specified\n");
    return -1;
  }
  file_prefix = strenv;

  strenv = getenv("CH_ID");
  if (strenv == NULL) {
    printf("environment variable CH_ID must be specified\n");
    return -1;
  }
  ch_id = atoi(strenv);

  strenv = getenv("INPUT_HEIGHT");
  if (strenv == NULL) {
    printf("environment variable INPUT_HEIGHT must be specified\n");
    return -1;
  }
  ds.input_height = atoi(strenv);

  strenv = getenv("INPUT_WIDTH");
  if (strenv == NULL) {
    printf("environment variable INPUT_WIDTH must be specified\n");
    return -1;
  }
  ds.input_width = atoi(strenv);

  strenv = getenv("OUTPUT_HEIGHT");
  if (strenv == NULL) {
    printf("environment variable OUTPUT_HEIGHT must be specified\n");
    return -1;
  }
  ds.output_height = atoi(strenv);

  strenv = getenv("OUTPUT_WIDTH");
  if (strenv == NULL) {
    printf("environment variable OUTPUT_WIDTH must be specified\n");
    return -1;
  }
  ds.output_width = atoi(strenv);
  
  sigset_t set;
  sigemptyset(&set);
  sigaddset(&set, SIGINT);
  sigprocmask(SIG_BLOCK, &set, NULL);

  ret = init_DPDK(file_prefix);
  if (ret < 0) {
    printf("init_DPDK failed: ret=%d\n", ret);
    return -1;
  }

  ret = init_FPGA(fpga_dev, ch_id);
  if (ret < 0) {
    printf("init_FPGA failed: ret=%d\n", ret);
    return -1;
  }

  if (signal(SIGINT, int_handler) == SIG_ERR) {
    return -1;
  }

  sigprocmask(SIG_UNBLOCK, &set, NULL);

  printf("controller started\n");
  while (true) {
    sleep(60);
  }

  return 0;
}
