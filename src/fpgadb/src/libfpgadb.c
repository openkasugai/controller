/*************************************************
* Copyright 2024 NTT Corporation, FUJITSU LIMITED
*************************************************/

#include <liblogging.h>
#include <libfpgactl.h>

#include <libfpgadb.h>

#include <parson.h>

#include <stdio.h>
#include <string.h>
#include <stdbool.h>
#include <errno.h>

#include <semaphore.h>


// LogLibFpga
#undef FPGA_LOGGER_LIBNAME
#define FPGA_LOGGER_LIBNAME "libfpgadb:  "


typedef enum DB_SEM_CMD{
  DB_SEM_INV, /**< Invalid parameter */
  DB_SEM_SH,  /**< Shared lock */
  DB_SEM_EX,  /**< Exclusive lock */
} db_sem_cmd_t;


static struct {
  uint32_t *parent_bitstream_id;
  uint32_t *child_bitstream_id;
} *db_dummy_bsid_list[FPGA_MAX_DEVICES];

static sem_t db_bsid_list_sem_sh;
static sem_t db_bsid_list_sem_ex;


static void __db_bsid_list_sem_init(void) {
  static volatile bool is_init = false;
  if (__sync_bool_compare_and_swap(&is_init, false, true)) {
    int ret;
    if ((ret = sem_init(&db_bsid_list_sem_sh, 0, 0)))
      llf_err(-ret, "Failed to initialize db_bsid_list_sem_sh.\n");
    if ((ret = sem_init(&db_bsid_list_sem_ex, 0, 1)))
      llf_err(-ret, "Failed to initialize db_bsid_list_sem_ex.\n");
  }
}


static inline void __db_bsid_list_sem_lock(db_sem_cmd_t cmd) {
  for (int cnt = 0; cnt < 100; cnt++) {
    sem_wait(&db_bsid_list_sem_ex);
    if (cmd == DB_SEM_SH) {
      sem_post(&db_bsid_list_sem_sh); // Success sh lock
      sem_post(&db_bsid_list_sem_ex); // Success ex unlock
      break;
    } else {
      if (sem_trywait(&db_bsid_list_sem_sh) == EAGAIN) // No sh lock found
        break;  // Break the loop with holding ex lock
      sem_post(&db_bsid_list_sem_sh); // Relock for sem_try_wait
      sem_post(&db_bsid_list_sem_ex);
    }
  }
}


static inline void __db_bsid_list_sem_unlock(db_sem_cmd_t cmd) {
  if (cmd == DB_SEM_SH)
    sem_wait(&db_bsid_list_sem_sh);
  else
    sem_post(&db_bsid_list_sem_ex);
}


static void __db_init(void) {
  static volatile bool is_init = false;
  if (__sync_bool_compare_and_swap(&is_init, false, true)) {
    __db_bsid_list_sem_init();
    memset(db_dummy_bsid_list, 0, sizeof(db_dummy_bsid_list));
  }
}


static int __fpga_db_get_bitstream_id(
  uint32_t dev_id,
  uint32_t *parent_bsid,
  uint32_t *child_bsid)
{
  // Check input
  if (!fpga_get_device(dev_id) || (!parent_bsid && !child_bsid)) {
    llf_err(INVALID_ARGUMENT, "%s(dev_id(%u),parent_bsid(%#llx),child_bsid(%#llx))\n",
      __func__, dev_id, (uintptr_t)parent_bsid, (uintptr_t)child_bsid);
    return -INVALID_ARGUMENT;
  }

  int ret;

  // Get from db_dummy_bsid_list if needed
  __db_bsid_list_sem_lock(DB_SEM_SH);
  if (db_dummy_bsid_list[dev_id]) {
    if (db_dummy_bsid_list[dev_id]->parent_bitstream_id && parent_bsid) {
      *parent_bsid = *db_dummy_bsid_list[dev_id]->parent_bitstream_id;
      parent_bsid = NULL;
    }
    if (db_dummy_bsid_list[dev_id]->child_bitstream_id && child_bsid) {
      *child_bsid = *db_dummy_bsid_list[dev_id]->child_bitstream_id;
      child_bsid = NULL;
    }
  }
  __db_bsid_list_sem_unlock(DB_SEM_SH);
  if (!parent_bsid && !child_bsid)
    return 0;

  // Update child bsid
  ret = fpga_update_info(dev_id);
  if (ret) {
    llf_err(-ret, "Failed fpga_update_info\n");
    return ret;
  }

  // Get parent/child bsid
  fpga_device_user_info_t info;
  ret = fpga_get_device_info(dev_id, &info);
  if (ret) {
    llf_err(-ret, "Failed fpga_get_device_info\n");
    return ret;
  }

  if (parent_bsid) *parent_bsid = info.bitstream_id.parent;
  if (child_bsid) *child_bsid = info.bitstream_id.child;

  return 0;
}


int fpga_db_get_bitstream_id(
  const char *device_name,
  uint32_t *parent_bsid,
  uint32_t *child_bsid)
{
  __db_init();
  // Check input
  if (!device_name || (!parent_bsid && !child_bsid)) {
    llf_err(INVALID_ARGUMENT, "%s(device_name(%s),parent_bsid(%#llx),child_bsid(%#llx))\n",
      __func__, device_name ? device_name : "<null>", (uintptr_t)parent_bsid, (uintptr_t)child_bsid);
    return -INVALID_ARGUMENT;
  }
  llf_dbg("%s(device_name(%s),parent_bsid(%#llx),child_bsid(%#llx))\n",
    __func__, device_name, (uintptr_t)parent_bsid, (uintptr_t)child_bsid);

  uint32_t dev_id;
  int ret;

  // Get dev_id from serial_id
  ret = fpga_get_dev_id(device_name, &dev_id);
  if (ret) {
    llf_err(-ret, "Failed fpga_get_dev_id\n");
    return ret;
  }

  return __fpga_db_get_bitstream_id(
    dev_id,
    parent_bsid,
    child_bsid);
}


int fpga_db_get_bitstream_id_by_dev_id(
  uint32_t dev_id,
  uint32_t *parent_bsid,
  uint32_t *child_bsid)
{
  __db_init();
  // Check input
  if (!parent_bsid && !child_bsid) {
    llf_err(INVALID_ARGUMENT, "%s(dev_id(%u),parent_bsid(%#llx),child_bsid(%#llx))\n",
      __func__, dev_id, (uintptr_t)parent_bsid, (uintptr_t)child_bsid);
    return -INVALID_ARGUMENT;
  }
  llf_dbg("%s(dev_id(%u),parent_bsid(%#llx),child_bsid(%#llx))\n",
    __func__, dev_id, (uintptr_t)parent_bsid, (uintptr_t)child_bsid);

  return __fpga_db_get_bitstream_id(
    dev_id,
    parent_bsid,
    child_bsid);
}


static int __fpga_db_get_device_config(
  const char *parent_bsid,
  const char *child_bsid,
  char **config_json)
{
  // Check input
  if (!config_json || !parent_bsid) {
    llf_err(INVALID_ARGUMENT, "%s(parent_bsid(%s),child_bsid(%s),config_json(%#llx))\n",
      __func__, parent_bsid ? parent_bsid : "<null>", child_bsid ? child_bsid : "<null>", (uintptr_t)config_json);
    return -INVALID_ARGUMENT;
  }

  char *config_file = NULL;
  int ret;

  // Get file name for FPGA configuration database from libfpgactl
  if ((ret = fpga_get_device_config_path(&config_file))) {
    llf_err(-ret, "Failed to get config file...\n");
    return ret;
  }

  // Parse configuration database file
  JSON_Value *root = json_parse_file_with_comments(config_file);
  if (!root) {
    llf_err(INVALID_ARGUMENT, "Failed to parse file: %s\n", config_file);
    free(config_file);
    return -INVALID_ARGUMENT;
  }

  // Get object for root
  JSON_Object *obj = json_value_get_object(root);
  if (!obj) {
    llf_err(INVALID_ARGUMENT, "Failed to get object: %s\n", config_file);
    free(config_file);
    json_value_free(root);
    return -INVALID_ARGUMENT;
  }
  free(config_file);

  // Get array for "configs"
  JSON_Array *array_configs = json_object_get_array(obj, "configs");
  if (!array_configs) {
    llf_err(INVALID_DATA, "Failed to get array configs\n");
    json_value_free(root);
    return -INVALID_DATA;
  }
  for (int index_configs = 0;; index_configs++) {
    // Get object for configs[index_configs]
    JSON_Object *obj_configs = json_array_get_object(array_configs, index_configs);
    if (!obj_configs) {
      ret = -INVALID_DATA;
      llf_err(INVALID_DATA, "Failed to access configs[%d]\n", index_configs);
      break;
    }
    // Get "parent-bitstream-id" from configs[index_configs]
    const char *string_parent_bsid = json_object_get_string(obj_configs, "parent-bitstream-id");
    if (!string_parent_bsid) {
      llf_warn(INVALID_DATA, "Failed to find parent-bitstream-id at configs[%d]\n", index_configs);
      continue;
    }
    // Check if "parent-bitstream-id" is the same with the argument or not
    if (!strcmp(string_parent_bsid, parent_bsid)) {
      // "parent-bitstream-id" is the same with the argument
      if (child_bsid) {
        int is_exist_child_bsid = 0;
        // When child bitstream id is not NULL, Get array for "child-bitstream-ids"
        JSON_Array *array_child_bsids = json_object_get_array(obj_configs, "child-bitstream-ids");
        if (!array_child_bsids) {
          ret = -INVALID_DATA;
          llf_err(INVALID_DATA, "Failed to find child-bitstream-ids at configs[%d]\n", index_configs);
          break;
        }
        for (int index_child_bsids = 0;; index_child_bsids++) {
          // Get object for child-bitstream-ids[index_child_bsids]
          JSON_Object *obj_child_bsids = json_array_get_object(array_child_bsids, index_child_bsids);
          if (!obj_child_bsids) break;
          // Get "child-bitstream-id" from child-bitstream-ids[index_child_bsids]
          const char *string_child_bsid = json_object_get_string(obj_child_bsids, "child-bitstream-id");
          if (!string_child_bsid) {
            llf_warn(INVALID_DATA, "Failed to find child-bitstream-id at child-bitstream-ids[%d]\n", index_child_bsids);
            json_array_remove(array_child_bsids, index_child_bsids);
            index_child_bsids = 0 - 1;
            continue;
          }
          // Check if "child-bitstream-id" is the same with the argument or not
          if (strcmp(string_child_bsid, child_bsid)) {
            json_array_remove(array_child_bsids, index_child_bsids);
            index_child_bsids = 0 - 1;
          } else {
            is_exist_child_bsid++;
          }
        }
        if (is_exist_child_bsid == 0) {
          llf_err(INVALID_DATA, "Failed to find valid child-bitstream-id in child-bitstream-ids\n");
          ret = -INVALID_DATA;
          break;
        } else if (is_exist_child_bsid > 1) {
          llf_warn(INVALID_DATA, "child-bitstream-id found %d times...\n", is_exist_child_bsid);
        }
      }
      if (ret) break;
      // serialize to string from config[index_configs]
      JSON_Value *val_config = json_object_get_wrapping_value(obj_configs);
      char *string_config = json_serialize_to_string_pretty(val_config);
      *config_json = strdup(string_config);
      json_free_serialized_string(string_config);
      break;
    }
  }

  json_value_free(root);

  return ret;
}


int fpga_db_get_device_config(
  const char *device_name,
  char **config_json)
{
  __db_init();
  // Check input
  if (!config_json || !device_name) {
    llf_err(INVALID_ARGUMENT, "%s(device_name(%s),config_json(%#llx))\n",
      __func__, device_name ? device_name : "<null>", (uintptr_t)config_json);
    return -INVALID_ARGUMENT;
  }
  llf_dbg("%s(device_name(%s),config_json(%#llx))\n", __func__, device_name , (uintptr_t)config_json);

  int ret;
  uint32_t dev_id;
  uint32_t parent;
  uint32_t child;
  char parent_bsid[9];
  char child_bsid[9];

  // Get Bitstream id from FPGA
  ret = fpga_get_dev_id(device_name, &dev_id);
  if (ret) {
    llf_err(-ret, "Failed fpga_get_dev_id\n");
    return ret;
  }
  ret = __fpga_db_get_bitstream_id(dev_id, &parent, &child);
  if (ret) {
    llf_err(-ret, "Failed to get Bitstream id for device[%u]\n", dev_id);
    return ret;
  }
  sprintf(parent_bsid, "%08x", parent);
  sprintf(child_bsid, "%08x", child);

  return __fpga_db_get_device_config(parent_bsid, child_bsid, config_json);
}


int fpga_db_get_device_config_by_dev_id(
  uint32_t dev_id,
  char **config_json)
{
  __db_init();
  // Check input
  if (!config_json) {
    llf_err(INVALID_ARGUMENT, "%s(dev_id(%u),config_json(%#llx))\n",
      __func__, dev_id, (uintptr_t)config_json);
    return -INVALID_ARGUMENT;
  }
  llf_dbg("%s(dev_id(%u),config_json(%#llx))\n", __func__, dev_id, (uintptr_t)config_json);

  int ret;
  uint32_t parent;
  uint32_t child;
  char parent_bsid[9];
  char child_bsid[9];

  // Get Bitstream id from FPGA
  ret = __fpga_db_get_bitstream_id(dev_id, &parent, &child);
  if (ret) {
    llf_err(-ret, "Failed to get Bitstream id for device[%u]\n", dev_id);
    return ret;
  }
  sprintf(parent_bsid, "%08x", parent);
  sprintf(child_bsid, "%08x", child);

  return __fpga_db_get_device_config(parent_bsid, child_bsid, config_json);
}


int fpga_db_get_device_config_by_bitstream_id(
  const char *parent_bsid,
  const char *child_bsid,
  char **config_json)
{
  // Check input
  if (!config_json || !parent_bsid) {
    llf_err(INVALID_ARGUMENT, "%s(parent_bsid(%s),child_bsid(%s),config_json(%#llx))\n",
      __func__, parent_bsid ? parent_bsid : "<null>", child_bsid ? child_bsid : "<null>", (uintptr_t)config_json);
    return -INVALID_ARGUMENT;
  }
  llf_dbg("%s(parent_bsid(%s),child_bsid(%s),config_json(%#llx))\n",
    __func__, parent_bsid, child_bsid ? child_bsid : "<null>", (uintptr_t)config_json);

  return __fpga_db_get_device_config(parent_bsid, child_bsid, config_json);
}


static int __fpga_db_get_child_bitstream_ids(
  const char *parent_bsid,
  char **child_bsid_list[])
{
  // Check input
  if (!parent_bsid || !child_bsid_list) {
    llf_err(INVALID_ARGUMENT, "%s(parent_bsid(%s),child_bsid_list(%#llx))\n",
      __func__, parent_bsid ? parent_bsid : "<null>", (uintptr_t)child_bsid_list);
    return -INVALID_ARGUMENT;
  }

  char *config_file = NULL;
  int ret;

  // Get file name for FPGA configuration database from libfpgactl
  if ((ret = fpga_get_device_config_path(&config_file))) {
    llf_err(-ret, "Failed to get config file...\n");
    return ret;
  }

  // Parse configuration database file
  JSON_Value *root = json_parse_file_with_comments(config_file);
  if (!root) {
    llf_err(INVALID_ARGUMENT, "Failed to parse file: %s\n", config_file);
    free(config_file);
    return -INVALID_ARGUMENT;
  }

  // Get object for root
  JSON_Object *obj = json_value_get_object(root);
  if (!obj) {
    llf_err(INVALID_ARGUMENT, "Failed to get object: %s\n", config_file);
    free(config_file);
    json_value_free(root);
    return -INVALID_ARGUMENT;
  }
  free(config_file);

  // Get array for "configs"
  JSON_Array *array_configs = json_object_get_array(obj, "configs");
  if (!array_configs) {
    llf_err(INVALID_DATA, "Failed to get array configs\n");
    json_value_free(root);
    return -INVALID_DATA;
  }
  for (int index_configs = 0;; index_configs++) {
    // Get object for configs[index_configs]
    JSON_Object *obj_configs = json_array_get_object(array_configs, index_configs);
    if (!obj_configs) {
      ret = -INVALID_DATA;
      llf_err(INVALID_DATA, "Failed to access configs[%d]\n", index_configs);
      break;
    }
    // Get "parent-bitstream-id" from configs[index_configs]
    const char *string_parent_bsid = json_object_get_string(obj_configs, "parent-bitstream-id");
    if (!string_parent_bsid) {
      llf_warn(INVALID_DATA, "Failed to find parent-bitstream-id at configs[%d]\n", index_configs);
      continue;
    }
    // Check if "parent-bitstream-id" is the same with the argument or not
    if (!strcmp(string_parent_bsid, parent_bsid)) {
      // "parent-bitstream-id" is the same with the argument
      // When child bitstream id is not NULL, Get array for "child-bitstream-ids"
      JSON_Array *array_child_bsids = json_object_get_array(obj_configs, "child-bitstream-ids");
      if (!array_child_bsids) {
        ret = -INVALID_DATA;
        llf_err(INVALID_DATA, "Failed to find child-bitstream-ids at configs[%d]\n", index_configs);
        break;
      }
      int max_child_bsids = json_array_get_count(array_child_bsids);
      if (max_child_bsids == 0) {
        ret = -INVALID_DATA;
        llf_err(INVALID_DATA, "Invalid data: child-bitstream-ids may be empty array.\n");
        break;
      }
      *child_bsid_list = (char**)malloc(sizeof(char*) * max_child_bsids + 1);
      if (!*child_bsid_list) {
        ret = -FAILURE_MEMORY_ALLOC;
        llf_err(FAILURE_MEMORY_ALLOC, "Failed to allocate memory for list:child_bsid_list[%d].\n",
          max_child_bsids + 1);
        break;
      }
      memset(*child_bsid_list, 0, sizeof(char*) * max_child_bsids + 1);
      for (int index_child_bsids = 0; index_child_bsids < max_child_bsids; index_child_bsids++) {
        // Get object for child-bitstream-ids[index_child_bsids]
        JSON_Object *obj_child_bsids = json_array_get_object(array_child_bsids, index_child_bsids);
        if (!obj_child_bsids) {
          ret = -INVALID_DATA;
          llf_err(INVALID_DATA, "Failed to access child-bitstream-ids[%d]\n", index_child_bsids);
          break;
        }
        // Get "child-bitstream-id" from child-bitstream-ids[index_child_bsids]
        const char *string_child_bsid = json_object_get_string(obj_child_bsids, "child-bitstream-id");
        if (!string_child_bsid) {
          ret = -INVALID_DATA;
          llf_err(INVALID_DATA, "Failed to find child-bitstream-id at child-bitstream-ids[%d]\n", index_child_bsids);
          break;
        }
        // Get child bsids
        (*child_bsid_list)[index_child_bsids] = strdup(string_child_bsid);
        if (!(*child_bsid_list)[index_child_bsids]) {
          ret = -FAILURE_MEMORY_ALLOC;
          llf_err(FAILURE_MEMORY_ALLOC, "Failed to allocate memory for child_bsid_list[%d].\n",
            index_child_bsids);
          break;
        }
      }
      if (ret)
        fpga_db_free_child_bitstream_ids(*child_bsid_list);
      break;
    }
  }

  json_value_free(root);

  return ret;
}


int fpga_db_get_child_bitstream_ids(
  const char *device_name,
  char **child_bsid_list[])
{
  __db_init();
  // Check input
  if (!device_name || !child_bsid_list) {
    llf_err(INVALID_ARGUMENT, "%s(device_name(%s),child_bsid_list(%#llx))\n",
      __func__, device_name ? device_name : "<null>", (uintptr_t)child_bsid_list);
    return -INVALID_ARGUMENT;
  }
  llf_dbg("%s(device_name(%s),child_bsid_list(%#llx))\n",
    __func__, device_name, (uintptr_t)child_bsid_list);

  int ret;
  uint32_t dev_id;
  uint32_t parent;
  char parent_bsid[9];

  // Get Bitstream id from FPGA
  ret = fpga_get_dev_id(device_name, &dev_id);
  if (ret) {
    llf_err(-ret, "Failed fpga_get_dev_id\n");
    return ret;
  }
  ret = __fpga_db_get_bitstream_id(dev_id, &parent, NULL);
  if (ret) {
    llf_err(-ret, "Failed to get Bitstream id for device[%u]\n", dev_id);
    return ret;
  }
  sprintf(parent_bsid, "%08x", parent);

  return __fpga_db_get_child_bitstream_ids(parent_bsid, child_bsid_list);
}


int fpga_db_get_child_bitstream_ids_by_dev_id(
  uint32_t dev_id,
  char **child_bsid_list[])
{
  __db_init();
  // Check input
  if (!child_bsid_list) {
    llf_err(INVALID_ARGUMENT, "%s(dev_id(%u),child_bsid_list(%#llx))\n",
      __func__, dev_id, (uintptr_t)child_bsid_list);
    return -INVALID_ARGUMENT;
  }
  llf_dbg("%s(dev_id(%u),child_bsid_list(%#llx))\n",
    __func__, dev_id, (uintptr_t)child_bsid_list);

  int ret;
  uint32_t parent;
  char parent_bsid[9];

  // Get Bitstream id from FPGA
  ret = __fpga_db_get_bitstream_id(dev_id, &parent, NULL);
  if (ret) {
    llf_err(-ret, "Failed to get Bitstream id for device[%u]\n", dev_id);
    return ret;
  }
  sprintf(parent_bsid, "%08x", parent);

  return __fpga_db_get_child_bitstream_ids(parent_bsid, child_bsid_list);
}


int fpga_db_get_child_bitstream_ids_by_parent(
  const char *parent_bsid,
  char **child_bsid_list[])
{
  // Check input
  if (!parent_bsid || !child_bsid_list) {
    llf_err(INVALID_ARGUMENT, "%s(parent_bsid(%s),child_bsid_list(%#llx))\n",
      __func__, parent_bsid ? parent_bsid : "<null>", (uintptr_t)child_bsid_list);
    return -INVALID_ARGUMENT;
  }
  llf_dbg("%s(parent_bsid(%s),child_bsid_list(%#llx))\n",
    __func__, parent_bsid, (uintptr_t)child_bsid_list);

  return __fpga_db_get_child_bitstream_ids(parent_bsid, child_bsid_list);
}


int fpga_db_free_child_bitstream_ids(
  char *child_bsid_list[])
{
  // Check input
  if (!child_bsid_list) {
    llf_err(INVALID_ARGUMENT, "%s(child_bsid_list(%#llx))\n",
      __func__, (uintptr_t)child_bsid_list);
    return -INVALID_ARGUMENT;
  }
  llf_dbg("%s(child_bsid_list(%#llx))\n",
    __func__, (uintptr_t)child_bsid_list);

  for (char **child_bsid = child_bsid_list; *child_bsid; child_bsid++)
    free(*child_bsid);
  free(child_bsid_list);

  return 0;
}


static int __fpga_db_disable_dummy_bitstream_by_dev_id(
  uint32_t dev_id)
{
  // Check input
  if (dev_id >= FPGA_MAX_DEVICES) {
    llf_err(INVALID_ARGUMENT, "%s(dev_id(%u))\n", __func__, dev_id);
    return -INVALID_ARGUMENT;
  }

  // Delete data
  if (!db_dummy_bsid_list[dev_id]) {
    // do nothing
  } else {
    if (db_dummy_bsid_list[dev_id]->parent_bitstream_id) {
      free(db_dummy_bsid_list[dev_id]->parent_bitstream_id);
      db_dummy_bsid_list[dev_id]->parent_bitstream_id = NULL;
    }
    if (db_dummy_bsid_list[dev_id]->child_bitstream_id) {
      free(db_dummy_bsid_list[dev_id]->child_bitstream_id);
      db_dummy_bsid_list[dev_id]->child_bitstream_id = NULL;
    }
    free(db_dummy_bsid_list[dev_id]);
    db_dummy_bsid_list[dev_id] = NULL;
  }

  return 0;
}


int fpga_db_disable_dummy_bitstream(
  const char *device_name)
{
  __db_init();
  // Check input
  if (!device_name) {
    llf_err(INVALID_ARGUMENT, "%s(device_name(%s))\n", __func__, "<null>");
    return -INVALID_ARGUMENT;
  }
  llf_dbg("%s(device_name(%s))\n", __func__, device_name);

  int ret;
  uint32_t dev_id;

  // Get dev_id from serial_id
  if ((ret = fpga_get_dev_id(device_name, &dev_id))) {
    llf_err(-ret, "Failed fpga_get_dev_id\n");
    return ret;
  }

  __db_bsid_list_sem_lock(DB_SEM_EX);
  ret = __fpga_db_disable_dummy_bitstream_by_dev_id(dev_id);
  __db_bsid_list_sem_unlock(DB_SEM_EX);

  return ret;
}


int fpga_db_disable_dummy_bitstream_by_dev_id(
  uint32_t dev_id)
{
  __db_init();
  // Check input
  if (!fpga_get_device(dev_id)) {
    llf_err(INVALID_ARGUMENT, "%s(dev_id(%u))\n", __func__, dev_id);
    return -INVALID_ARGUMENT;
  }
  llf_dbg("%s(dev_id(%u))\n", __func__, dev_id);

  int ret;

  __db_bsid_list_sem_lock(DB_SEM_EX);
  ret = __fpga_db_disable_dummy_bitstream_by_dev_id(dev_id);
  __db_bsid_list_sem_unlock(DB_SEM_EX);

  return ret;
}


static int __fpga_db_enable_dummy_bitstream_by_dev_id(
  uint32_t dev_id,
  const uint32_t *parent_dummy_bsid,
  const uint32_t *child_dummy_bsid)
{
  // Check input
  if ((dev_id >= FPGA_MAX_DEVICES) || (!parent_dummy_bsid && !child_dummy_bsid)) {
    llf_err(INVALID_ARGUMENT,
      "%s(dev_id(%u),parent_dummy_bsid(%#llx),child_dummy_bsid(%#llx))\n",
      __func__, dev_id, (uintptr_t)parent_dummy_bsid, (uintptr_t)child_dummy_bsid);
    return -INVALID_ARGUMENT;
  }

  if (!db_dummy_bsid_list[dev_id]) {
    // Allocate memory for db_dummy_bsid_list
    db_dummy_bsid_list[dev_id] = (typeof(db_dummy_bsid_list[dev_id]))malloc(sizeof(typeof(*db_dummy_bsid_list[dev_id])));
    memset(db_dummy_bsid_list[dev_id], 0, sizeof(*db_dummy_bsid_list[dev_id]));
  } else {
    // Memory for db_dummy_bsid_list is already alocated(=update),
    //  so delete existing data at first.
    if (db_dummy_bsid_list[dev_id]->parent_bitstream_id) {
      free(db_dummy_bsid_list[dev_id]->parent_bitstream_id);
      db_dummy_bsid_list[dev_id]->parent_bitstream_id = NULL;
    }
    if (db_dummy_bsid_list[dev_id]->child_bitstream_id) {
      free(db_dummy_bsid_list[dev_id]->child_bitstream_id);
      db_dummy_bsid_list[dev_id]->child_bitstream_id = NULL;
    }
  }

  if (parent_dummy_bsid) {
    uint32_t *dummy_bsid = db_dummy_bsid_list[dev_id]->parent_bitstream_id
                          = (uint32_t*)malloc(sizeof(uint32_t));
    if (!dummy_bsid) {
      llf_err(FAILURE_MEMORY_ALLOC, "Failed to allocate memory for dummy parent_bsid\n");
      __fpga_db_disable_dummy_bitstream_by_dev_id(dev_id);
      return -FAILURE_MEMORY_ALLOC;
    }
    *dummy_bsid = *parent_dummy_bsid;
  } else {
    db_dummy_bsid_list[dev_id]->parent_bitstream_id = NULL;
  }

  if (child_dummy_bsid) {
    uint32_t *dummy_bsid = db_dummy_bsid_list[dev_id]->child_bitstream_id
                          = (uint32_t*)malloc(sizeof(uint32_t));
    if (!dummy_bsid) {
      llf_err(FAILURE_MEMORY_ALLOC, "Failed to allocate memory for dummy child_bsid\n");
      __fpga_db_disable_dummy_bitstream_by_dev_id(dev_id);
      return -FAILURE_MEMORY_ALLOC;
    }
    *dummy_bsid = *child_dummy_bsid;
  }

  return 0;
}


int fpga_db_enable_dummy_bitstream(
  const char *device_name,
  const uint32_t *parent_dummy_bsid,
  const uint32_t *child_dummy_bsid)
{
  __db_init();
  // Check input
  if (!device_name || (!parent_dummy_bsid && !child_dummy_bsid)) {
    llf_err(INVALID_ARGUMENT,
      "%s(device_name(%s),parent_dummy_bsid(%#llx),child_dummy_bsid(%#llx))\n",
      __func__, device_name ? device_name : "<null>",
      (uintptr_t)parent_dummy_bsid, (uintptr_t)child_dummy_bsid);
    return -INVALID_ARGUMENT;
  }
  llf_dbg("%s(device_name(%s),parent_dummy_bsid(%#llx),child_dummy_bsid(%#llx))\n",
    __func__, device_name ? device_name : "<null>",
    (uintptr_t)parent_dummy_bsid, (uintptr_t)child_dummy_bsid);

  int ret;
  uint32_t dev_id;

  // Get dev_id from serial_id
  if ((ret = fpga_get_dev_id(device_name, &dev_id))) {
    llf_err(-ret, "Failed fpga_get_dev_id\n");
    return ret;
  }

  __db_bsid_list_sem_lock(DB_SEM_EX);
  ret = __fpga_db_enable_dummy_bitstream_by_dev_id(dev_id, parent_dummy_bsid, child_dummy_bsid);
  __db_bsid_list_sem_unlock(DB_SEM_EX);

  return ret;
}


int fpga_db_enable_dummy_bitstream_by_dev_id(
  uint32_t dev_id,
  const uint32_t *parent_dummy_bsid,
  const uint32_t *child_dummy_bsid)
{
  __db_init();
  // Check input
  if (!fpga_get_device(dev_id) || (!parent_dummy_bsid && !child_dummy_bsid)) {
    llf_err(INVALID_ARGUMENT,
      "%s(dev_id(%u),parent_dummy_bsid(%#llx),child_dummy_bsid(%#llx))\n",
      __func__, dev_id, (uintptr_t)parent_dummy_bsid, (uintptr_t)child_dummy_bsid);
    return -INVALID_ARGUMENT;
  }
  llf_dbg("%s(dev_id(%u),parent_dummy_bsid(%#llx),child_dummy_bsid(%#llx))\n",
    __func__, dev_id, (uintptr_t)parent_dummy_bsid, (uintptr_t)child_dummy_bsid);

  int ret;

  __db_bsid_list_sem_lock(DB_SEM_EX);
  ret = __fpga_db_enable_dummy_bitstream_by_dev_id(dev_id, parent_dummy_bsid, child_dummy_bsid);
  __db_bsid_list_sem_unlock(DB_SEM_EX);

  return ret;
}
