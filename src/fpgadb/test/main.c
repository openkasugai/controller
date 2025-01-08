/*************************************************
* Copyright 2024 NTT Corporation, FUJITSU LIMITED
*************************************************/

#include <libfpgactl.h>
#include <liblogging.h>

#include <libfpgadb.h>

#include <stdio.h>
#include <string.h>


static int execute_with_fpgas(void)
{
  int ret;
  char **device_list = NULL;
  uint32_t u32_parent;
  char char_parent[9];
  char output_file_name[512];
  char *config_json = NULL;
  FILE *fp = NULL;
  char **child_list = NULL;
  const char *database_file_path = "test-bitstream_id-config-table.json";

  // Initialize FPGA
  ret = fpga_scan_devices();
  if (ret < 0) {
    printf(" ! Failed fpga_scan_devices in %d\n", ret);
    return -1;
  }
  printf(" * Execute with FPGAs\n");
  printf(" *  Use '%s'\n", database_file_path);
  if (fpga_set_device_config_path(database_file_path))
    printf(" !  Failed to set '%s'\n", database_file_path);

  // Get device name list
  if ((ret = fpga_get_device_list(&device_list))) {
    printf(" ! Failed fpga_get_device_list in %d\n", ret);
    fpga_finish();
    return -1;
  }
  printf(" * Get FPGAs' serial ids\n");

  for (char **device_name = device_list; *device_name; device_name++) {
    printf(" * Device[%s]\n", *device_name);
    /**
     * Serial_id(path) -> Config json
     */
    if ((ret = fpga_db_get_device_config(*device_name, &config_json))) {
      printf(" ! Failed fpga_db_get_device_config in %d\n", ret);
      fpga_finish();
      fpga_release_device_list(device_list);
      return -1;
    }
    // Output data to file
    sprintf(output_file_name, "config-device-%s.json", *device_name);
    fp = fopen(output_file_name, "w");
    if (!fp) fp = stdout;
    fprintf(fp, "%s\n", config_json);
    free(config_json);
    if (fp != stdout) {
      fclose(fp);
      printf(" *    Create '%s'\n", output_file_name);
    }

    /**
     * Serial_id(path) -> BSID
     */
    // Get Bitstream-IDs from FPGA
    if ((ret = fpga_db_get_bitstream_id(*device_name, &u32_parent, NULL))) {
      printf(" ! Failed fpga_db_get_bitstream_id in %d\n", ret);
      fpga_finish();
      fpga_release_device_list(device_list);
      return -1;
    }
    // Print result
    printf(" * Parent BSID : %08x\n", u32_parent);
    // Create bitstream_id from num to string to match "bitstream_id-config-table.json"
    sprintf(char_parent, "%08x", u32_parent);

    /**
     * BSID(parent) -> available BSIDs(child)
     */
    ret = fpga_db_get_child_bitstream_ids_by_parent(char_parent, &child_list);
    if (ret) {
      printf(" ! Failed fpga_db_get_child_bitstream_ids_by_parent in %d\n", ret);
      fpga_finish();
      fpga_release_device_list(device_list);
      return -1;
    }
    printf(" * This parent-bitstream-id(%s) can use the below child-bitstream-id:\n", char_parent);
    for (char **child_bsid = child_list; *child_bsid; child_bsid++) {
      printf(" * - %s\n", *child_bsid);
      /**
       * BSID(parent+child) -> Config json
       */
      // Get config json string from json file
      ret = fpga_db_get_device_config_by_bitstream_id(char_parent, *child_bsid, &config_json);
      if (ret) {
        printf(" ! Failed fpga_db_get_device_config_by_bitstream_id in %d\n", ret);
        fpga_finish();
        fpga_db_free_child_bitstream_ids(child_list);
        fpga_release_device_list(device_list);
        return -1;
      }
      // Output data to file
      sprintf(output_file_name, "config-available-parent(%s)-child(%s).json",
        char_parent, *child_bsid);
      fp = fopen(output_file_name, "w");
      if (!fp) fp = stdout;
      fprintf(fp, "%s\n", config_json);
      free(config_json);
      if (fp != stdout) {
        fclose(fp);
        printf(" *    Create '%s'\n", output_file_name);
      }
    }
    fpga_db_free_child_bitstream_ids(child_list);
  }

  fpga_release_device_list(device_list);

  // Finalize FPGA
  fpga_finish();

  return 0;
}


static int execute_without_fpgas(void)
{
  int ret;
  char char_parent[9];
  char output_file_name[512];
  char *config_json = NULL;
  FILE *fp = NULL;
  char **child_list = NULL;
  uint32_t parent_bsid;
  const char *database_file_path = "dummy-bitstream_id-config-table.json";

  printf(" * Execute without FPGAs\n");
  printf(" *  Use '%s'\n", database_file_path);
  if (fpga_set_device_config_path(database_file_path))
    printf(" !  Failed to set '%s'\n", database_file_path);

  printf(" * Input Parent Bitstream ID in hex.\n > ");
  scanf("%08x", &parent_bsid);

  printf(" * Parent BSID : %08x\n", parent_bsid);
  // Create bitstream_id from num to string to match "bitstream_id-config-table.json"
  sprintf(char_parent, "%08x", parent_bsid);

  /**
   * BSID(parent) -> available BSIDs(child)
   */
  ret = fpga_db_get_child_bitstream_ids_by_parent(char_parent, &child_list);
  if (ret) {
    printf(" ! Failed fpga_db_get_child_bitstream_ids_by_parent in %d\n", ret);
    return -1;
  }
  printf(" * This parent-bitstream-id(%s) can use the below child-bitstream-id:\n", char_parent);
  for (char **child_bsid = child_list; *child_bsid; child_bsid++) {
    printf(" * - %s\n", *child_bsid);
    /**
     * BSID(parent+child) -> Config json
     */
    // Get config json string from json file
    ret = fpga_db_get_device_config_by_bitstream_id(char_parent, *child_bsid, &config_json);
    if (ret) {
      printf(" ! Failed fpga_db_get_device_config_by_bitstream_id in %d\n", ret);
      fpga_db_free_child_bitstream_ids(child_list);
      return -1;
    }
    // Output data to file
    sprintf(output_file_name, "config-available-parent(%s)-child(%s).json",
      char_parent, *child_bsid);
    fp = fopen(output_file_name, "w");
    if (!fp) fp = stdout;
    fprintf(fp, "%s\n", config_json);
    free(config_json);
    if (fp != stdout) {
      fclose(fp);
      printf(" *    Create '%s'\n", output_file_name);
    }
  }
  fpga_db_free_child_bitstream_ids(child_list);

  return 0;
}


int main(
  int argc,
  char **argv)
{
  int ret;

  if ((ret = libfpga_log_parse_args(argc, argv)) < 0) {
    printf(" ! Failed libfpga_log_parse_args in %d\n", ret);
    return -1;
  }
  argc-=ret;
  argv+=ret;

  if (argc != 1 && !strcmp(argv[1], "--without-fpga"))
    return execute_without_fpgas();
  else
    return execute_with_fpgas();
}
