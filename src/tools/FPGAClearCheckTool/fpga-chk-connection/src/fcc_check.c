/*************************************************
* Copyright 2025 NTT Corporation , FUJITSU LIMITED
*************************************************/

#include <fcc_check.h>
#include <fcc_prm.h>
#include <fcc_log.h>

#include <liblogging.h>
#include <libfpgactl.h>
#include <libdma.h>
#include <libchain.h>

#ifndef FPGA_LANE_MAX
#define FPGA_LANE_MAX       2
#endif
#ifndef FPGA_CID_MIN_LLDMA
#define FPGA_CID_MIN_LLDMA  0
#endif
#ifndef FPGA_CID_MAX_LLDMA
#define FPGA_CID_MAX_LLDMA  15
#endif
#ifndef FPGA_CID_MIN_PTU
#define FPGA_CID_MIN_PTU    1
#endif
#ifndef FPGA_CID_MAX_PTU
#define FPGA_CID_MAX_PTU    511
#endif
#ifndef FPGA_FCHID_MIN
#define FPGA_FCHID_MIN      0
#endif
#ifndef FPGA_FCHID_MAX
#define FPGA_FCHID_MAX      511
#endif


static int __fcc_check_initialize(
  void)
{
  int ret;
  libfpga_log_set_level(LIBFPGA_LOG_NOTHING);
  // libfpga_log_set_output_stdout();

  ret = fpga_scan_devices();
  if (ret < 0) {
    fcc_log_errorf(" ! Cannot open FPGAs: Xpcie driver mey not be loaded...\n");
  } else if (ret == 0) {
    fcc_log_printf(" ! FPGA not found: FPGA may not be written with child bitstream...\n");
  }

  return ret;
}

static void __fcc_check_finalize(
  void)
{
  fpga_finish();
}

static int __fcc_check_lldma(
  int *index)
{
  fpga_set_refqueue_polling_timeout(1);
  fpga_set_refqueue_polling_interval(1);

  int session_index = 0;
  int found_df = 0;
  int ret;
  dma_info_t dma_info;

  for (const fcc_prm_lldma_t *ent = fcc_prm_get_lldma_list(); ent && *ent; ent++) {
    // Print header
    fcc_log_printf("- Dataflow-session : %d\n", index ? *index + session_index : session_index);
    fcc_log_printf("\t- Parameters\n");
    fcc_log_printf("\t\t- Extif_id : 0(LLDMA)\n");
    fcc_log_printf("\t\t- Connector_id : %s\n", *ent);

    // Check lldma setting
    ret = fpga_lldma_queue_setup(*ent, &dma_info);
    if (ret == 0) {
      // Dataflow settings found
      fcc_log_printf("\t- Result : Found\n");
      fcc_log_printf("\t- Found Dataflow Settings\n");
      fpga_device_user_info_t info = {.device_file_path="dummy"};
      fpga_get_device_info(dma_info.dev_id, &info);
      fcc_log_printf("\t\t- Device : %s\n", info.device_file_path);
      fcc_log_printf("\t\t- Direction : %s\n", dma_info.dir == DMA_HOST_TO_DEV ? "RX(Ingress)" : "TX(Egress)");
      fcc_log_printf("\t\t- Connection_id : %u\n", dma_info.chid);
      fpga_lldma_queue_finish(&dma_info);
      fcc_prm_push_err_lldma_list(session_index);
      found_df++;
    } else if (ret == -CONNECTOR_ID_MISMATCH) {
      // Dataflow settings Not found
      fcc_log_printf("\t- Result : Not found\n");
    } else {
      // Failed to check lldma setting
      fcc_log_errorf("\t- Result : Error(%d) at fpga_lldma_queue_setup()\n", ret);
    }
    session_index++;
  }

  if (index)
    *index = *index + session_index;

  return found_df;
}

static int __fcc_check_ptu(
  int *index)
{
  int session_index = 0;
  int found_df = 0;
  int ret;
  uint32_t status;
  uint32_t dev_id;

  for (const fcc_prm_ptu_t *ent = fcc_prm_get_ptu_list(); ent && *ent; ent++) {
    // Print header
    fcc_log_printf("- Dataflow-session : %d\n", index ? *index + session_index : session_index);
    fcc_log_printf("\t- Parameters\n");
    fcc_log_printf("\t\t- Device : %s\n", (*ent)->device);
    fcc_log_printf("\t\t- Lane : %u\n", (*ent)->lane);
    fcc_log_printf("\t\t- Extif_id : %u(%s)\n", (*ent)->extif_id, FCC_PRM_EXTIF_TO_STR((*ent)->extif_id));
    fcc_log_printf("\t\t- Connection_id : %u\n", (*ent)->cid);

    // Convert name into dev_id
    ret = fpga_get_dev_id((*ent)->device, &dev_id);
    if (ret) {
      // Failed to convert
      fcc_log_errorf("\t- Result : Error(%d) at fpga_get_dev_id()\n", ret);
    } else {
      // Check ptu session
      ret = fpga_chain_get_con_status(dev_id, (*ent)->lane, (*ent)->extif_id, (*ent)->cid, &status);
      if (ret == 0) {
        if (status == 0) {
          // Dataflow settings Not found
          fcc_log_printf("\t- Result : Not found\n");
        } else {
          // Dataflow settings found
          fcc_log_printf("\t- Result : Found\n");
          fcc_prm_push_err_ptu_list(session_index);
          found_df++;
        }
      } else {
        // Failed to check ptu session
        fcc_log_errorf("\t- Result : Error(%d) at fpga_chain_get_con_status()\n", ret);
      }
    }
    session_index++;
  }

  if (index)
    *index = *index + session_index;

  return found_df;
}

static int __fcc_check_chain(
  int *index)
{
  int chain_index = 0;
  int found_df = 0;
  int ret;

  uint32_t dev_id;
  uint32_t ingr_extif_id;
  uint32_t ingr_cid;
  uint32_t egr_extif_id;
  uint32_t egr_cid;

  for (const fcc_prm_chain_t *ent = fcc_prm_get_chain_list(); ent && *ent; ent++) {
    // Print header
    fcc_log_printf("- Dataflow-chain : %d\n", index ? *index + chain_index : chain_index);
    fcc_log_printf("\t- Parameters\n");
    fcc_log_printf("\t\t- Device : %s\n", (*ent)->device);
    fcc_log_printf("\t\t- Lane : %u\n", (*ent)->lane);
    fcc_log_printf("\t\t- Function_channel_id : %u\n", (*ent)->fchid);
    fcc_log_printf("\t\t- Direction : %u(%s)\n", (*ent)->dir, FCC_PRM_DIR_TO_STR((*ent)->dir));

    // Convert name into dev_id
    ret = fpga_get_dev_id((*ent)->device, &dev_id);
    if (ret) {
      // Failed to convert
      fcc_log_errorf("\t- Result : Error(%d) at fpga_get_dev_id()\n", ret);
    } else {
      // Check chain setting
      ret = fpga_chain_read_soft_table(dev_id, (*ent)->lane, (*ent)->fchid,
        &ingr_extif_id, &ingr_cid, &egr_extif_id, &egr_cid);
      if (ret == 0) {
        if ((*ent)->dir == FCC_PRM_DIR_INGR) {
          // When checking ingress chain
          if (ingr_extif_id == -1 && ingr_cid == -1) {
            // Dataflow settings Not found
            fcc_log_printf("\t- Result : Not found\n");
          } else {
            // Dataflow settings found
            fcc_log_printf("\t- Result : Found\n");
            fcc_log_printf("\t\t- Extif_id : %u(%s)\n", ingr_extif_id, FCC_PRM_EXTIF_TO_STR(ingr_extif_id));
            fcc_log_printf("\t\t- Connection_id : %u\n", ingr_cid);
            fcc_prm_push_err_chain_list(chain_index);
            found_df++;
          }
        } else {
          // When checking egress chain
          if (egr_extif_id == -1 && egr_cid == -1) {
            // Dataflow settings Not found
            fcc_log_printf("\t- Result : Not found\n");
          } else {
            // Dataflow settings found
            fcc_log_printf("\t- Result : Found\n");
            fcc_log_printf("\t\t- Extif_id : %u(%s)\n", egr_extif_id, FCC_PRM_EXTIF_TO_STR(egr_extif_id));
            fcc_log_printf("\t\t- Connection_id : %u\n", egr_cid);
            fcc_prm_push_err_chain_list(chain_index);
            found_df++;
          }
        }
      } else {
        // Error
        fcc_log_errorf("\t- Result : Error(%d) at fpga_chain_read_soft_table()\n", ret);
      }
    }
    chain_index++;
  }

  if (index)
    *index = *index + chain_index;

  return found_df;
}

static int __fcc_check_print_summary(
  void)
{
  fcc_log_printf("- Found Dataflow Index\n");

  // Print Dataflow session
  fcc_log_printf("\t- Dataflow-session : [");
  const fcc_prm_index_t *heads[2] = {
    fcc_prm_get_err_lldma_list(),
    fcc_prm_get_err_ptu_list()};
  int offsets[2] = {
    0,
    fcc_prm_get_lldma_list_size()};
  for (int idx_head = 0; idx_head < 2;idx_head++) {
    for (const fcc_prm_index_t *ent = heads[idx_head]; ent && *ent; ent++) {
      if ((ent != heads[0]) && (heads[0] || (ent != heads[1])))
        fcc_log_printf(",");
      fcc_log_printf("%u", (**ent) + offsets[idx_head]);
    }
  }
  fcc_log_printf("]\n");

  // Print Dataflow chain
  fcc_log_printf("\t- Dataflow-chain : [");
  const fcc_prm_index_t *head = fcc_prm_get_err_chain_list();
  for (const fcc_prm_index_t *ent = head; ent && *ent; ent++) {
    if (ent != head)
      fcc_log_printf(",");
    fcc_log_printf("%u", **ent);
  }
  fcc_log_printf("]\n");

  return 0;
}


int fcc_check_prm(
  void)
{
  int ret;
  int index = 0;
  int found_df = 0;

  // Initialize FPGA
  ret = __fcc_check_initialize();
  if (ret <= 0) return ret;

  // Check exteranal interface session settings
  fcc_log_printf("# Extif Settings\n");
  ret = __fcc_check_lldma(&index);
  found_df += ret;
  ret = __fcc_check_ptu(&index);
  found_df += ret;
  fcc_log_printf("\n");

  // Check chain settings
  fcc_log_printf("# Chain Settings\n");
  ret = __fcc_check_chain(NULL);
  found_df += ret;
  fcc_log_printf("\n");

  // Print summary
  fcc_log_printf("# Summary\n");
  __fcc_check_print_summary();
  fcc_log_printf("\n");

  // Finalize FPGA
  __fcc_check_finalize();

  return found_df;
}


static int __fcc_check_dump_lldma(
  void)
{
  int found_df = 0;
  int ret;

#ifdef SUPPORT_OTHER_THAN_LIBFPGA_API
  for (uint32_t dev_id = 0; dev_id < fpga_get_num(); dev_id++) {
    // Prepare
    fpga_ioctl_chsts_t ioctl_chsts;
    int is_print_comma;
    fpga_device_t *dev = fpga_get_device(dev_id);
    if (!dev) {
      fcc_log_errorf(" ! Failed to get device[%u]\n", dev_id);
      continue;
    }
    // Get lldma channel status RX
    ioctl_chsts.dir = DMA_HOST_TO_DEV;
    ret = ioctl(dev->fd, XPCIE_DEV_LLDMA_GET_CH_STAT, &ioctl_chsts);
    if (ret) {
      fcc_log_errorf(" ! Failed XPCIE_DEV_LLDMA_GET_CH_STAT\n");
      continue;
    }
    // Print lldma channel status RX avail
    is_print_comma = 0;
    fcc_log_printf("- LLDMA channels(%s : RX : Avail)  : [", dev->name);
    for (int shift = 0; shift < sizeof(uint32_t) * 8; shift++) {
      if ((ioctl_chsts.avail_status >> shift) & 0b1) {
        if (is_print_comma)
          fcc_log_printf(",");
        else
          is_print_comma = 1;
        fcc_log_printf("%u", shift);
      }
    }
    fcc_log_printf("]\n");
    // Print lldma channel status RX active
    is_print_comma = 0;
    fcc_log_printf("- LLDMA channels(%s : RX : Active) : [", dev->name);
    for (int shift = 0; shift < sizeof(uint32_t) * 8; shift++) {
      if ((ioctl_chsts.active_status >> shift) & 0b1) {
        if (is_print_comma)
          fcc_log_printf(",");
        else
          is_print_comma = 1;
        fcc_log_printf("%u", shift);
      }
    }
    fcc_log_printf("]\n");
    // Get lldma channel status TX
    ioctl_chsts.dir = DMA_DEV_TO_HOST;
    ret = ioctl(dev->fd, XPCIE_DEV_LLDMA_GET_CH_STAT, &ioctl_chsts);
    if (ret) {
      fcc_log_errorf(" ! Failed XPCIE_DEV_LLDMA_GET_CH_STAT\n");
      continue;
    }
    // Print lldma channel status TX avail
    fcc_log_printf("- LLDMA channels(%s : TX : Avail)  : [", dev->name);
    is_print_comma = 0;
    for (int shift = 0; shift < sizeof(uint32_t) * 8; shift++) {
      if ((ioctl_chsts.avail_status >> shift) & 0b1) {
        if (is_print_comma)
          fcc_log_printf(",");
        else
          is_print_comma = 1;
        fcc_log_printf("%u", shift);
      }
    }
    fcc_log_printf("]\n");
    // Print lldma channel status TX active
    is_print_comma = 0;
    fcc_log_printf("- LLDMA channels(%s : TX : Active) : [", dev->name);
    for (int shift = 0; shift < sizeof(uint32_t) * 8; shift++) {
      if ((ioctl_chsts.active_status >> shift) & 0b1) {
        if (is_print_comma)
          fcc_log_printf(",");
        else
          is_print_comma = 1;
        fcc_log_printf("%u", shift);
      }
    }
    fcc_log_printf("]\n");
  }
#endif

  for (uint32_t dev_id = 0; dev_id < fpga_get_num(); dev_id++) {
    for (uint32_t lane = 0; lane < FPGA_LANE_MAX; lane++) {
      for (uint32_t cid = FPGA_CID_MIN_LLDMA; cid <= FPGA_CID_MAX_LLDMA; cid++) {
        uint32_t status;
        uint32_t extif_id = 0;
        ret = fpga_chain_get_con_status(dev_id, lane, extif_id, cid, &status);
        if (ret == 0 && status == 1) {
          fcc_log_printf("- LLDMA-session : %d\n", found_df);
          fpga_device_user_info_t info = {.device_file_path="dummy"};
          fpga_get_device_info(dev_id, &info);
          fcc_log_printf("\t- Device : %s\n", info.device_file_path);
          fcc_log_printf("\t- Lane : %u\n", lane);
          fcc_log_printf("\t- Extif_id : 0(LLDMA)\n");
          fcc_log_printf("\t- Connection_id : %u\n", cid);
          found_df++;
        }
      }
    }
  }

  return found_df;
}

static int __fcc_check_dump_ptu(
  void)
{
  int found_df = 0;
  int ret;

  for (uint32_t dev_id = 0; dev_id < fpga_get_num(); dev_id++) {
    for (uint32_t lane = 0; lane < FPGA_LANE_MAX; lane++) {
      for (uint32_t cid = FPGA_CID_MIN_PTU; cid <= FPGA_CID_MAX_PTU; cid++) {
        uint32_t status;
        uint32_t extif_id = 1;
        ret = fpga_chain_get_con_status(dev_id, lane, extif_id, cid, &status);
        if (ret == 0 && status == 1) {
          fcc_log_printf("- PTU-session : %d\n", found_df);
          fpga_device_user_info_t info = {.device_file_path="dummy"};
          fpga_get_device_info(dev_id, &info);
          fcc_log_printf("\t- Device : %s\n", info.device_file_path);
          fcc_log_printf("\t- Lane : %u\n", lane);
          fcc_log_printf("\t- Extif_id : 1(PTU)\n");
          fcc_log_printf("\t- Connection_id : %u\n", cid);
          found_df++;
        }
      }
    }
  }

  return found_df;
}

static int __fcc_check_dump_chain(
  void)
{
  int found_df = 0;
  int ret;

  for (uint32_t dev_id = 0; dev_id < fpga_get_num(); dev_id++) {
    for (uint32_t lane = 0; lane < FPGA_LANE_MAX; lane++) {
      for (uint32_t fchid = FPGA_FCHID_MIN; fchid <= FPGA_FCHID_MAX; fchid++) {
        uint32_t ingr_extif_id, ingr_cid;
        uint32_t egr_extif_id, egr_cid;
        ret = fpga_chain_read_soft_table(dev_id, lane, fchid,
          &ingr_extif_id, &ingr_cid, &egr_extif_id, &egr_cid);
        if (ret == 0) {
          if ((ingr_extif_id & ingr_cid) != -1) {
            uint32_t direction = FCC_PRM_DIR_INGR;
            fcc_log_printf("- Chain-settings : %d\n", found_df);
            fpga_device_user_info_t info = {.device_file_path="dummy"};
            fpga_get_device_info(dev_id, &info);
            fcc_log_printf("\t- Device : %s\n", info.device_file_path);
            fcc_log_printf("\t- Lane : %u\n", lane);
            fcc_log_printf("\t- Function_channel_id : %u\n", fchid);
            fcc_log_printf("\t- Direction : %u(%s)\n", direction, FCC_PRM_DIR_TO_STR(direction));
            fcc_log_printf("\t- Extif_id : %u(%s)\n", ingr_extif_id, FCC_PRM_EXTIF_TO_STR(ingr_extif_id));
            fcc_log_printf("\t- Connection_id : %u\n", ingr_cid);
            found_df++;
          }
          if ((egr_extif_id & egr_cid) != -1) {
            uint32_t direction = FCC_PRM_DIR_EGR;
            fcc_log_printf("- Chain-settings : %d\n", found_df);
            fpga_device_user_info_t info = {.device_file_path="dummy"};
            fpga_get_device_info(dev_id, &info);
            fcc_log_printf("\t- Device : %s\n", info.device_file_path);
            fcc_log_printf("\t- Lane : %u\n", lane);
            fcc_log_printf("\t- Function_channel_id : %u\n", fchid);
            fcc_log_printf("\t- Direction : %u(%s)\n", direction, FCC_PRM_DIR_TO_STR(direction));
            fcc_log_printf("\t- Extif_id : %u(%s)\n", egr_extif_id, FCC_PRM_EXTIF_TO_STR(egr_extif_id));
            fcc_log_printf("\t- Connection_id : %u\n", egr_cid);
            found_df++;
          }
        }
      }
    }
  }

  return found_df;
}

int fcc_check_dump(
  void)
{
  int ret;
  int found_df = 0;

  // Initialize FPGA
  ret = __fcc_check_initialize();
  if (ret <= 0) return ret;

  // Check exteranal interface session settings
  fcc_log_printf("# Extif Settings\n");
  ret = __fcc_check_dump_lldma();
  found_df += ret;
  ret = __fcc_check_dump_ptu();
  found_df += ret;
  fcc_log_printf("\n");

  // Check chain settings
  fcc_log_printf("# Chain Settings\n");
  ret = __fcc_check_dump_chain();
  found_df += ret;
  fcc_log_printf("\n");

  // Finalize FPGA
  __fcc_check_finalize();

  return found_df;
}
