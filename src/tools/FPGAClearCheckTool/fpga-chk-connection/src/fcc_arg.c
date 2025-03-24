/*************************************************
* Copyright 2025 NTT Corporation , FUJITSU LIMITED
*************************************************/

#include <fcc_arg.h>
#include <fcc_log.h>
#include <fcc_prm.h>
#include <fcc_json.h>

#include <stdlib.h>
#include <string.h>
#include <getopt.h>


static const struct option fcc_long_options[] = {
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
  { "dir",                 required_argument, NULL, 0x100 },
  { "direction",           required_argument, NULL, 0x100 }, // dupe
  { "connector_id",        required_argument, NULL, 'k'   },
  { "connector-id",        required_argument, NULL, 'k'   }, // dupe
  { "matching-key",        required_argument, NULL, 'k'   }, // dupe
  { "matching_key",        required_argument, NULL, 'k'   }, // dupe
  { "json_params",         required_argument, NULL, 'j'   },
  { "json-params",         required_argument, NULL, 'j'   }, // dupe
  { "input_json_file",     required_argument, NULL, 'i'   },
  { "input-json-file",     required_argument, NULL, 'i'   }, // dupe
  { "output_json_file",    required_argument, NULL, 'o'   },
  { "output-json-file",    required_argument, NULL, 'o'   }, // dupe
  { "dump",                no_argument,       NULL, 0x101 },
  { "help",                no_argument,       NULL, 'h'   },
  { NULL,                  0,                 0,    0     }, // sentinel
};

static const char fcc_short_options[] = {
    "d:l:f:e:c:k:j:i:o:h"
};


static void __fcc_arg_print_usage(void) {
  fcc_log_printf("%s : v%08x\n", APP_NAME, fcc_log_get_app_version());
  fcc_log_printf("usage: %s [-dlfeckjio <PARAMETER>] [--dump] [-h]\n", APP_NAME);
  fcc_log_printf("       -d/--device <DEVICE>         : Device file path[/dev/xpcie_<UUID>,<UUID>]\n");
  fcc_log_printf("       -l/--lane <LANE>             : Lane number[0-1]\n");
  fcc_log_printf("       -f/--fchid <FCHID>           : Function channel id[0-511]\n");
  fcc_log_printf("       -e/--extif_id <EXTIF_ID>     : External interface id[lldma,LLDMA,0,ptu,PTU,1]\n");
  fcc_log_printf("       -c/--cid <CID>               : Connection id[0-15(LLDMA),1-511(PTU)]\n");
  fcc_log_printf("          --dir <DIRECTION>         : Direction[ingress,1,egress,2,both,3]\n");
  fcc_log_printf("       -k/--connector_id <KEY>      : Connector_id[String]\n");
  fcc_log_printf("       -j/--json_params <JSON>      : JSON format parameters[String]\n");
  fcc_log_printf("       -i/--input_json_file <FILE>  : JSON format input file path[String]\n");
  fcc_log_printf("       -o/--output_json_file <FILE> : JSON format output file path[String]\n");
  fcc_log_printf("                                    : This file will be created only when settings are found.\n");
  fcc_log_printf("          --dump                    : Dump all settings\n");
  fcc_log_printf("       -h/--help                    : Print this message\n");
  fcc_log_printf("\n");
}


static int __fcc_arg_parse_args(
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

  // param var
  char *device           = NULL;
  uint32_t lane          = -1;
  uint32_t direction     = FCC_PRM_DIR_BOTH;
  uint32_t fchid         = -1;
  uint32_t extif_id      = -1;
  uint32_t cid           = -1;

  int option;
  while (
    (option = getopt_long(
      argc,
      argv,
      fcc_short_options,
      fcc_long_options,
      NULL)
    ) != EOF
  ) {
    switch (option) {
      case 'd':   // device
        device = optarg;
        break;
      case 'l':   // lane
        lane = atoi(optarg);
        break;
      case 'f':   // fchid
        fchid = atoi(optarg);
        break;
      case 'e':   // extif_id
        if (strcmp(optarg, "lldma")      == 0) extif_id = FCC_PRM_EXTIF_LLDMA;
        else if (strcmp(optarg, "LLDMA") == 0) extif_id = FCC_PRM_EXTIF_LLDMA;
        else if (strcmp(optarg, "ptu")   == 0) extif_id = FCC_PRM_EXTIF_PTU;
        else if (strcmp(optarg, "PTU")   == 0) extif_id = FCC_PRM_EXTIF_PTU;
        else extif_id = atoi(optarg);
        break;
      case 'c':   // connection_id
        cid = atoi(optarg);
        break;
      case 0x100: // direction
        if (strcmp(optarg, "ingress")      == 0) direction = FCC_PRM_DIR_INGR;
        else if (strcmp(optarg, "ingr")    == 0) direction = FCC_PRM_DIR_INGR;
        else if (strcmp(optarg, "Ingress") == 0) direction = FCC_PRM_DIR_INGR;
        else if (strcmp(optarg, "Ingr")    == 0) direction = FCC_PRM_DIR_INGR;
        else if (strcmp(optarg, "INGRESS") == 0) direction = FCC_PRM_DIR_INGR;
        else if (strcmp(optarg, "INGR")    == 0) direction = FCC_PRM_DIR_INGR;
        else if (strcmp(optarg, "egress")  == 0) direction = FCC_PRM_DIR_EGR;
        else if (strcmp(optarg, "egr")     == 0) direction = FCC_PRM_DIR_EGR;
        else if (strcmp(optarg, "Egress")  == 0) direction = FCC_PRM_DIR_EGR;
        else if (strcmp(optarg, "Egr")     == 0) direction = FCC_PRM_DIR_EGR;
        else if (strcmp(optarg, "EGRESS")  == 0) direction = FCC_PRM_DIR_EGR;
        else if (strcmp(optarg, "EGR")     == 0) direction = FCC_PRM_DIR_EGR;
        else if (strcmp(optarg, "both")    == 0) direction = FCC_PRM_DIR_BOTH;
        else if (strcmp(optarg, "Both")    == 0) direction = FCC_PRM_DIR_BOTH;
        else if (strcmp(optarg, "BOTH")    == 0) direction = FCC_PRM_DIR_BOTH;
        else direction = atoi(optarg);
        break;
      case 'k':   // connector_id
        // Set lldma prm
        ret = fcc_prm_push_lldma_list(optarg);
        if (ret) goto out;
        break;
      case 'j':   // json_prm
        ret = fcc_json_parse_string(optarg);
        if (ret) goto out;
        break;
      case 'i':   // input_json_file
        ret = fcc_json_parse_file(optarg);
        if (ret) goto out;
        break;
      case 'o':   // output_json_file
        ret = fcc_prm_set_output_file_path(optarg);
        if (ret) goto out;
        break;
      case 0x101: // dump
        fcc_prm_set_is_dump(true);
        ret = 0;
        goto out;
      case 'h':   // help
        __fcc_arg_print_usage();
        ret = -FCC_PRM_ERRNO_HELP;
        goto out;
      default:
        fcc_log_errorf("Cannot parse option : %s\n", argv[optind - 1]);
        __fcc_arg_print_usage();
        ret = -1;
        goto out;
    }
  }

  // Set ptu prm
  if (device != NULL && lane != -1 && extif_id != -1 && cid != -1) {
    ret = fcc_prm_push_ptu_list(device, lane, extif_id, cid);
    if (ret) goto out;
  }

  // Set chain prm
  if (device != NULL && lane != -1 && fchid != -1) {
    if (direction & FCC_PRM_DIR_INGR) {
      ret = fcc_prm_push_chain_list(device, lane, fchid, FCC_PRM_DIR_INGR);
      if (ret) goto out;
    }
    if (direction & FCC_PRM_DIR_EGR) {
      ret = fcc_prm_push_chain_list(device, lane, fchid, FCC_PRM_DIR_EGR);
      if (ret) goto out;
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

int fcc_arg_parse_args(
  int argc,
  char **argv)
{
  int ret;
  int ret_parsed = 0;

  const int  old_optind = optind;
  const int  old_optopt = optopt;
  char* const old_optarg = optarg;

  do {
    ret = __fcc_arg_parse_args(argc, argv);
    if (ret < 0) {
      ret_parsed = ret;
      break;
    } else if (ret == 0) {
      break;
    }
    argc -= ret;
    argv += ret;
    ret_parsed += ret;
  } while (1);

  if (!fcc_prm_get_is_dump() && ret_parsed == 0) {
    __fcc_arg_print_usage();
    ret_parsed = -FCC_PRM_ERRNO_HELP;
  }

  optind = old_optind;
  optopt = old_optopt;
  optarg = old_optarg;

  return ret_parsed;
}
