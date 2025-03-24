/*************************************************
* Copyright 2025 NTT Corporation , FUJITSU LIMITED
*************************************************/

#include <fcc_log.h>
#include <fcc_arg.h>
#include <fcc_check.h>
#include <fcc_json.h>
#include <fcc_prm.h>


static void fcc_cleanup(void) {
  fcc_prm_free_output_file_path();
  fcc_prm_free_lldma_list();
  fcc_prm_free_ptu_list();
  fcc_prm_free_chain_list();
  fcc_prm_free_err_lldma_list();
  fcc_prm_free_err_ptu_list();
  fcc_prm_free_err_chain_list();
}

int main (
  int argc,
  char **argv)
{
  int ret = 0;
  ret = fcc_arg_parse_args(argc, argv);
  if (ret == -FCC_PRM_ERRNO_HELP) {
    fcc_cleanup();
    return 0;
  }
  if (ret < 0) {
    fcc_log_errorf(" ! Failed to parse arguments\n");
    return -1;
  }

  if (fcc_prm_get_is_dump()) {
    // Call dump function
    ret = fcc_check_dump();
  } else {
    // Check dataflow settings
    ret = fcc_check_prm();
    if (ret > 0) {
      // Write into the output json file
      fcc_json_create_output_file();
    }
  }

  // Free memory
  fcc_cleanup();

  return ret;
}

