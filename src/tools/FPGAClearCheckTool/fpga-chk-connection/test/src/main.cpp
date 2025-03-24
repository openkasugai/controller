/*************************************************
* Copyright 2025 NTT Corporation , FUJITSU LIMITED
*************************************************/

#include "connection.hpp"
#include "logger.hpp"

#include <stdio.h>
#include <stdlib.h>
#include <stdarg.h>
#include <string.h>
#include <unistd.h>
#include <getopt.h>
#include <time.h>
#include <errno.h>
#include <arpa/inet.h>

#include <liblldma.h>

#ifndef APP_NAME
#define APP_NAME "dummy"
#endif

#ifndef LOG_FILENAME
#define LOG_FILENAME APP_NAME ".log"
#endif

#ifndef LOG_FILENAME_DAEMON
#define LOG_FILENAME_DAEMON APP_NAME "-daemon.log"
#endif

#ifndef LOG_FILENAME_USER
#define LOG_FILENAME_USER APP_NAME "-user.log"
#endif

#define PRM_EXTIF_LLDMA 0
#define PRM_EXTIF_PTU   1

#define PRM_DIR_INGR    0b01
#define PRM_DIR_EGR     0b10
#define PRM_DIR_BOTH    0b11


static const struct option long_options[] = {
  { "device",              required_argument, NULL, 'd'   },
  { "lane",                required_argument, NULL, 'l'   },
  { "fchid",               required_argument, NULL, 'f'   },
  { "function_channel_id", required_argument, NULL, 'f'   }, // dupe
  { "function-channel-id", required_argument, NULL, 'f'   }, // dupe
  { "extif_id",            required_argument, NULL, 'e'   },
  { "extif-id",            required_argument, NULL, 'e'   }, // dupe
  { "cid",                 required_argument, NULL, 'c'   },
  { "connection_id",       required_argument, NULL, 'c'   }, // dupe
  { "connection-id",       required_argument, NULL, 'c'   }, // dupe
  { "connector_id",        required_argument, NULL, 'k'   },
  { "connector-id",        required_argument, NULL, 'k'   }, // dupe
  { "matching-key",        required_argument, NULL, 'k'   }, // dupe
  { "matching_key",        required_argument, NULL, 'k'   }, // dupe
  { "quit",                no_argument,       NULL, 'q'   },
  { "quit-daemon",         no_argument,       NULL, 'q'   }, // dupe
  { "quit_daemon",         no_argument,       NULL, 'q'   }, // dupe
  { "help",                no_argument,       NULL, 'h'   },
  { "user",                no_argument,       NULL, 'u'   },
  { "dir",                 required_argument, NULL, 0x100 },
  { "direction",           required_argument, NULL, 0x100 }, // dupe
  { "daemon",              no_argument,       NULL, 0x101 },
  { "chid",                required_argument, NULL, 0x102 },
  { "fpga_port",           required_argument, NULL, 0x103 },
  { "fpga-port",           required_argument, NULL, 0x103 }, // dupe
  { "host_addr",           required_argument, NULL, 0x104 },
  { "host-addr",           required_argument, NULL, 0x104 }, // dupe
  { "host_port",           required_argument, NULL, 0x105 },
  { "host-port",           required_argument, NULL, 0x105 }, // dupe
  { "delete",              no_argument,       NULL, 0x106 },
  { "create",              no_argument,       NULL, 0x107 },
  { "fpga_addr",           required_argument, NULL, 0x108 },
  { "fpga-addr",           required_argument, NULL, 0x108 }, // dupe
  { NULL,                  0,                 0,    0     }, // sentinel
};

static const char short_options[] = {
    "d:l:f:e:c:k:j:i:o:u:qh"
};


static void print_usage(FILE *fp = stdout) {
  fprintf(fp, "%s\n", APP_NAME);
  fprintf(fp, "usage: %s [-dlfeckjio <PARAMETER>] [--dump] [-h]\n", APP_NAME);
  fprintf(fp, "          --daemon                  : Launch daemon\n");
  fprintf(fp, "       -d/--device <DEVICE>         : Device file path[/dev/xpcie_<UUID>,<UUID>]\n");
  fprintf(fp, "       -l/--lane <LANE>             : Lane number[0-1]\n");
  fprintf(fp, "       -f/--fchid <FCHID>           : Function channel id[0-511]\n");
  fprintf(fp, "       -e/--extif_id <EXTIF_ID>     : External interface id[lldma,LLDMA,0,ptu,PTU,1]\n");
  fprintf(fp, "       -c/--cid <CID>               : Connection id[0-15(LLDMA),1-511(PTU)]\n");
  fprintf(fp, "          --dir <DIRECTION>         : Direction[ingress,1,egress,2,both,3]\n");
  fprintf(fp, "       -k/--connector_id <KEY>      : Connector_id[String]\n");
  fprintf(fp, "          --chid <CHID>             : DMA channel id[0-15]\n");
  fprintf(fp, "          --fpga_addr <ADDR>        : Set FPGA address\n");
  fprintf(fp, "          --fpga_port <PORT>        : Set FPGA port\n");
  fprintf(fp, "          --host_addr <ADDR>        : Set Host address\n");
  fprintf(fp, "          --host_port <PORT>        : Set Host port\n");
  fprintf(fp, "          --create                  : Create a Connection\n");
  fprintf(fp, "          --delete                  : Delete a Connection\n");
  fprintf(fp, "       -u/--user                    : Execute TCP user process\n");
  fprintf(fp, "       -q/--quit                    : Quit daemon\n");
  fprintf(fp, "       -h/--help                    : Print this message\n");
  fprintf(fp, "\n");
}


// param var
static bool g_is_user           = false;
static bool g_is_daemon         = false;
static bool g_is_create         = true;
static bool g_is_quit           = false;
static char *g_device           = NULL;
static char *g_conid            = NULL;
static uint32_t g_direction     = PRM_DIR_INGR;
static uint32_t g_chid          = -1;
static uint32_t g_lane          = -1;
static uint32_t g_fchid         = -1;
static uint32_t g_extif_id      = -1;
static uint32_t g_cid           = -1;
static in_addr_t g_fpga_addr    = 0xC0A8006F;
static in_port_t g_fpga_port    = -1;
static in_addr_t g_host_addr    = -1;
static in_port_t g_host_port    = -1;


static int parse_args_daemon(
  int argc,
  char **argv)
{
  // opt var
  const int  old_optind = optind;
  const int  old_optopt = optopt;
  char* const old_optarg = optarg;
  optind = 1;
  opterr = 0;

  int ret = 0;
  int option;
  while (
    (option = getopt_long(
      argc,
      argv,
      short_options,
      long_options,
      NULL)
    ) != EOF
  ) {
    switch (option) {
      case 0x101: // daemon
        g_is_daemon = true;
        goto out;
      default:
        break;
    }
  }

out:
  optind = old_optind;
  optopt = old_optopt;
  optarg = old_optarg;

  return ret;
}


static int parse_args(
  int argc,
  char **argv)
{
  int ret;

  // opt var
  const int  old_optind = optind;
  const int  old_optopt = optopt;
  char* const old_optarg = optarg;
  optind = 1;
  opterr = 0;

  int option;
  while (
    (option = getopt_long(
      argc,
      argv,
      short_options,
      long_options,
      NULL)
    ) != EOF
  ) {
    switch (option) {
      case 'd':   // device
        g_device = optarg;
        break;
      case 'l':   // lane
        g_lane = atoi(optarg);
        break;
      case 'f':   // fchid
        g_fchid = atoi(optarg);
        break;
      case 'e':   // extif_id
        if (strcmp(optarg, "lldma")      == 0) g_extif_id = PRM_EXTIF_LLDMA;
        else if (strcmp(optarg, "LLDMA") == 0) g_extif_id = PRM_EXTIF_LLDMA;
        else if (strcmp(optarg, "ptu")   == 0) g_extif_id = PRM_EXTIF_PTU;
        else if (strcmp(optarg, "PTU")   == 0) g_extif_id = PRM_EXTIF_PTU;
        else g_extif_id = atoi(optarg);
        break;
      case 'c':   // connection_id
        g_cid = atoi(optarg);
        break;
      case 'k':   // connector_id
        g_conid = optarg;
        break;
      case 'h':   // help
        print_usage();
        ret = 0;
        goto out;
      case 'q':
        g_is_quit = true;
        break;
      case 'u':
        g_is_user = true;
        break;
      case 0x100: // direction
        if (strcmp(optarg, "ingress")      == 0) g_direction = PRM_DIR_INGR;
        else if (strcmp(optarg, "ingr")    == 0) g_direction = PRM_DIR_INGR;
        else if (strcmp(optarg, "Ingress") == 0) g_direction = PRM_DIR_INGR;
        else if (strcmp(optarg, "Ingr")    == 0) g_direction = PRM_DIR_INGR;
        else if (strcmp(optarg, "INGRESS") == 0) g_direction = PRM_DIR_INGR;
        else if (strcmp(optarg, "INGR")    == 0) g_direction = PRM_DIR_INGR;
        else if (strcmp(optarg, "egress")  == 0) g_direction = PRM_DIR_EGR;
        else if (strcmp(optarg, "egr")     == 0) g_direction = PRM_DIR_EGR;
        else if (strcmp(optarg, "Egress")  == 0) g_direction = PRM_DIR_EGR;
        else if (strcmp(optarg, "Egr")     == 0) g_direction = PRM_DIR_EGR;
        else if (strcmp(optarg, "EGRESS")  == 0) g_direction = PRM_DIR_EGR;
        else if (strcmp(optarg, "EGR")     == 0) g_direction = PRM_DIR_EGR;
        else if (strcmp(optarg, "both")    == 0) g_direction = PRM_DIR_BOTH;
        else if (strcmp(optarg, "Both")    == 0) g_direction = PRM_DIR_BOTH;
        else if (strcmp(optarg, "BOTH")    == 0) g_direction = PRM_DIR_BOTH;
        else g_direction = atoi(optarg);
        break;
      case 0x101: // daemon
        break;
      case 0x102: // chid
        g_chid = atoi(optarg);
        break;
      case 0x108: // fpga_addr
        sscanf(optarg, "%x", &g_fpga_addr);
        break;
      case 0x103: // fpga_port
        sscanf(optarg, "%hu", &g_fpga_port);
        break;
      case 0x104: // host_addr
        sscanf(optarg, "%x", &g_host_addr);
        break;
      case 0x105: // host_port
        sscanf(optarg, "%hu", &g_host_port);
        break;
      case 0x106: // delete
        g_is_create = false;
        break;
      case 0x107: // create
        g_is_create = true;
        break;
      default:
        printf("Cannot parse option : %s\n", argv[optind - 1]);
        print_usage();
        ret = -1;
        goto out;
    }
  }

  if (optind >= 0)
    argv[optind - 1] = argv[0];
  ret = optind - 1;

out:
  optind = old_optind;
  optopt = old_optopt;
  optarg = old_optarg;

  return ret;
}


static int connection_server(
  void)
{
  logger *plog = new logger(LOG_FILENAME_DAEMON);
  if (!plog) {
    printf(" ! Failed to new logger\n");
    return -1;
  }
  connection *pdf = new connection(plog);
  if (!pdf) {
    plog->print(" ! Failed to new connection\n");
    return -1;
  }

  pdf->server();

  delete pdf;
  delete plog;

  return 0;
}


static int connection_client(
  void)
{
  logger *plog = new logger(LOG_FILENAME);
  if (!plog) {
    printf(" ! Failed to new logger\n");
    return -1;
  }
  connection *pdf = new connection(plog);
  if (!pdf) {
    plog->print(" ! Failed to new connection\n");
    return -1;
  }

  int ret = -1;
  if (g_is_quit) {
    ret = pdf->client_request_quit();
  } else if (g_is_user) {
    plog->print(" * This process is user, so reopen as %s\n", LOG_FILENAME_USER);
    plog->close();
    plog->open(LOG_FILENAME_USER);
    uint8_t laddr[4], raddr[4];
    char host_addr[17];
    char fpga_addr[17];
    laddr[0] = (g_host_addr >> 0)  & 0xFF;
    laddr[1] = (g_host_addr >> 8)  & 0xFF;
    laddr[2] = (g_host_addr >> 16) & 0xFF;
    laddr[3] = (g_host_addr >> 24) & 0xFF;
    raddr[0] = (g_fpga_addr >> 0)  & 0xFF;
    raddr[1] = (g_fpga_addr >> 8)  & 0xFF;
    raddr[2] = (g_fpga_addr >> 16) & 0xFF;
    raddr[3] = (g_fpga_addr >> 24) & 0xFF;
    sprintf(host_addr, "%u.%u.%u.%u",
      laddr[3], laddr[2], laddr[1], laddr[0]);
    sprintf(fpga_addr, "%u.%u.%u.%u",
      raddr[3], raddr[2], raddr[1], raddr[0]);
    int cnt;
    const int cnt_max = 5;
    for (cnt = 0; cnt < cnt_max; cnt++) {
      if (g_direction == PRM_DIR_INGR)
        ret = pdf->host_connect_fpga(
          fpga_addr,
          g_fpga_port,
          host_addr,
          g_host_port);
      else
        ret = pdf->host_accept_fpga(
          host_addr,
          g_host_port);
      if (ret >= 0) {
        plog->print(" * Got fd : %d\n", ret);
        break;
      }
    }
    if (cnt != cnt_max) {
      sleep(1);
      ret = pdf->host_close_fpga(ret);
    }
  } else if (
    g_device != NULL &&
    g_conid  != NULL &&
    g_chid   != (uint32_t)-1)
  {
    // Set lldma prm
    if (g_is_create)
      ret = pdf->client_request_create(
        g_device,
        g_conid,
        g_direction == PRM_DIR_INGR ? DMA_HOST_TO_DEV : DMA_DEV_TO_HOST,
        g_chid);
    else
      ret = pdf->client_request_delete(
        g_device,
        g_conid,
        g_direction == PRM_DIR_INGR ? DMA_HOST_TO_DEV : DMA_DEV_TO_HOST,
        g_chid);
  } else if (
    g_device    != NULL         &&
    g_lane      != (uint32_t)-1 &&
    g_fpga_port != (uint32_t)-1 &&
    g_host_addr != (uint32_t)-1 &&
    g_host_port != (uint32_t)-1)
  {
    // Set ptu prm
    if (g_is_create)
      ret = pdf->client_request_create(
        g_device,
        g_lane,
        g_fpga_port,
        g_host_addr,
        g_host_port,
        g_direction == PRM_DIR_INGR);
    else
      ret = pdf->client_request_delete(
        g_device,
        g_lane,
        g_fpga_port,
        g_host_addr,
        g_host_port,
        g_direction == PRM_DIR_INGR);
  } else if (
    g_device   != NULL         &&
    g_lane     != (uint32_t)-1 &&
    g_fchid    != (uint32_t)-1 &&
    g_extif_id != (uint32_t)-1 &&
    g_cid      != (uint32_t)-1)
  {
    // Set chain prm
    if (g_is_create)
      ret = pdf->client_request_create(
        g_device,
        g_lane,
        g_fchid,
        g_extif_id,
        g_cid,
        1,
        g_direction == PRM_DIR_INGR ? 1 : 0,
        1,
        g_direction == PRM_DIR_INGR ? 0 : 1);
    else
      ret = pdf->client_request_delete(
        g_device,
        g_lane,
        g_fchid,
        g_extif_id,
        g_cid,
        1,
        g_direction == PRM_DIR_INGR ? 1 : 0,
        1,
        g_direction == PRM_DIR_INGR ? 0 : 1);
  } else {
    plog->print(" ! No matching parameters found...\n");
  }

  delete pdf;
  delete plog;

  return ret;
}

int main (
  int argc,
  char **argv)
{
  int ret;

  ret = parse_args_daemon(argc, argv);
  if (ret < 0) {
    printf(" ! Failed to parse args...\n");
    return -1;
  }

  if (g_is_daemon) {
    printf(" # Execute server as daemon process\n");
    ret = connection_server();
  } else {
    printf(" # Execute client\n");
    ret = parse_args(argc, argv);
    if (ret < 0) {
      printf(" ! Failed to parse args...\n");
      return -1;
    }
    if (ret == 0)
      return 0;
    ret = connection_client();
  }

  return ret;
}
