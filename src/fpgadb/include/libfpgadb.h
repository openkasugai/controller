/*************************************************
* Copyright 2024 NTT Corporation, FUJITSU LIMITED
*************************************************/
/**
 * @file libfpgadb.h
 * @brief Header file for plain database APIs for libfpga
 */

#ifndef LIBFPGADB_INCLUDE_LIBFPGADB_H_
#define LIBFPGADB_INCLUDE_LIBFPGADB_H_

#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

/**
 * @brief Get FPGA's bitstream id
 * @param[in] device_name
 *   FPGA's device path or serial id
 * @param[out] parent_bsid
 *   Parent bitstream id
 * @param[out] child_bsid
 *   Child bitstream id
 * @retval 0
 *   success
 * @retval -INVALID_ARGUMENT
 *   e.g.) `device_name` is null
 * @return retval of fpga_get_dev_id()
 * @return retval of fpga_update_info()
 * @return retval of fpga_get_device_info()
 *
 * @details
 *   Get `parent_bsid` and `child_bsid` from FPGA's register value.@n
 *   User should allocate memory for at least one of `parent_bsid` or `child_bsid`.
 *
 * @sa fpga_db_enable_dummy_bitstream
 * @sa fpga_db_enable_dummy_bitstream_by_dev_id
 * @sa fpga_db_disable_dummy_bitstream
 * @sa fpga_db_disable_dummy_bitstream_by_dev_id
 */
int fpga_db_get_bitstream_id(
        const char *device_name,
        uint32_t *parent_bsid,
        uint32_t *child_bsid);

/**
 * @brief Get FPGA's bitstream id by dev_id
 * @param[in] dev_id
 *   device id obtained by fpga_dev_init()
 * @param[out] parent_bsid
 *   Parent bitstream id
 * @param[out] child_bsid
 *   Child bitstream id
 * @retval 0
 *   success
 * @retval -INVALID_ARGUMENT
 *   e.g.) both `parent_bsid` and `child_bsid` are null
 * @return retval of fpga_get_dev_id()
 * @return retval of fpga_update_info()
 * @return retval of fpga_get_device_info()
 *
 * @details
 *   Get `parent_bsid` and `child_bsid` from FPGA's register value.@n
 *   User should allocate memory for at least one of `parent_bsid` or `child_bsid`.
 *
 * @sa fpga_db_enable_dummy_bitstream
 * @sa fpga_db_enable_dummy_bitstream_by_dev_id
 * @sa fpga_db_disable_dummy_bitstream
 * @sa fpga_db_disable_dummy_bitstream_by_dev_id
 */
int fpga_db_get_bitstream_id_by_dev_id(
        uint32_t dev_id,
        uint32_t *parent_bsid,
        uint32_t *child_bsid);

/**
 * @brief Get FPGA's configuration information
 * @param[in] device_name
 *   FPGA's device path or serial id
 * @param[out] config_json
 *   Configuration information as json format
 * @retval 0
 *   success
 * @retval -INVALID_ARGUMENT
 *   e.g.) `device_name` is null
 * @retval -INVALID_DATA
 *   e.g.) data parameter is missing
 * @return retval of fpga_db_get_bitstream_id()
 *
 * @details
 *   Get configuration information from json file with matching FPGA's bitstream id@n
 *   `*config_json` will be allocated by this API,
 *    so please free() `*config_json` explicitly.
 *
 * @sa fpga_db_enable_dummy_bitstream
 * @sa fpga_db_enable_dummy_bitstream_by_dev_id
 * @sa fpga_db_disable_dummy_bitstream
 * @sa fpga_db_disable_dummy_bitstream_by_dev_id
 */
int fpga_db_get_device_config(
        const char *device_name,
        char **config_json);

/**
 * @brief Get FPGA's configuration information by dev_id
 * @param[in] dev_id
 *   device id obtained by fpga_dev_init().
 * @param[out] config_json
 *   Configuration information as json format
 * @retval 0
 *   success
 * @retval -INVALID_ARGUMENT
 *   e.g.) `config_json` is null
 * @retval -INVALID_DATA
 *   e.g.) data parameter is missing
 * @return retval of fpga_db_get_bitstream_id()
 *
 * @details
 *   Get configuration information from json file with matching FPGA's bitstream id@n
 *   `*config_json` will be allocated by this API,
 *    so please free() `*config_json` explicitly.
 *
 * @sa fpga_db_enable_dummy_bitstream
 * @sa fpga_db_enable_dummy_bitstream_by_dev_id
 * @sa fpga_db_disable_dummy_bitstream
 * @sa fpga_db_disable_dummy_bitstream_by_dev_id
 */
int fpga_db_get_device_config_by_dev_id(
        uint32_t dev_id,
        char **config_json);

/**
 * @brief Get configuration information matching parent/child bitstream ids
 * @param[in] parent_bsid
 *   Parent bitstream id
 * @param[in] child_bsid
 *   Child bitstream id
 * @param[out] config_json
 *   Configuration information as json format
 * @retval 0
 *   success
 * @retval -INVALID_ARGUMENT
 *   e.g.) `config_json` is null
 * @retval -INVALID_DATA
 *   e.g.) data parameter is missing
 * @return retval of fpga_db_get_bitstream_id()
 *
 * @details
 *   Get configuration information from json file with matching arguments's bitstream id@n
 *   `parent_bsid` is mandatory, but `child_bsid`is arbitary.@n
 *   When `child_bsid` is NULL, `config_json` will return the all configurations for `parent_bsid`.
 *   `*config_json` will be allocated by this API,
 *    so please free() `*config_json` explicitly.
 * @remarks
 *   This API does not need FPGA.
 */
int fpga_db_get_device_config_by_bitstream_id(
        const char *parent_bsid,
        const char *child_bsid,
        char **config_json);

/**
 * @brief Get FPGA's available child bitstream ids
 * @param[in] device_name
 *   FPGA's device path or serial id
 * @param[out] child_bsid_list
 *   Child bitstream id list
 * @retval 0
 *   success
 * @retval -INVALID_ARGUMENT
 *   e.g.) `device_name` is null
 * @retval -INVALID_DATA
 *   e.g.) data parameter is missing
 * @retval -FAILURE_MEMORY_ALLOC
 *   Failed to allocate memory
 * @return retval of fpga_get_dev_id()
 * @return retval of fpga_db_get_bitstream_id()
 *
 * @details
 *   Get all child bitstream ids for the target FPGA.@n
 *   `*child_bsid_list` and `(*child_bsid_list)[*]` will be allocated by this API,
 *    so please fpga_db_free_child_bitstream_ids() `child_bsid_list` explicitly.@n
 *   The sentinel of `*child_bsid_list` is NULL.
 *
 * @sa fpga_db_enable_dummy_bitstream
 * @sa fpga_db_enable_dummy_bitstream_by_dev_id
 * @sa fpga_db_disable_dummy_bitstream
 * @sa fpga_db_disable_dummy_bitstream_by_dev_id
 */
int fpga_db_get_child_bitstream_ids(
        const char *device_name,
        char **child_bsid_list[]);

/**
 * @brief Get FPGA's available child bitstream ids by dev_id
 * @param[in] dev_id
 *   device id obtained by fpga_dev_init().
 * @param[out] child_bsid_list
 *   Child bitstream id list
 * @retval 0
 *   success
 * @retval -INVALID_ARGUMENT
 *   e.g.) `child_bsid_list` is null
 * @retval -INVALID_DATA
 *   e.g.) data parameter is missing
 * @retval -FAILURE_MEMORY_ALLOC
 *   Failed to allocate memory
 * @return retval of fpga_get_dev_id()
 * @return retval of fpga_db_get_bitstream_id()
 *
 * @details
 *   Get all child bitstream ids for the target FPGA.@n
 *   `*child_bsid_list` and `(*child_bsid_list)[*]` will be allocated by this API,
 *    so please fpga_db_free_child_bitstream_ids() `child_bsid_list` explicitly.@n
 *   The sentinel of `*child_bsid_list` is NULL.
 *
 * @sa fpga_db_enable_dummy_bitstream
 * @sa fpga_db_enable_dummy_bitstream_by_dev_id
 * @sa fpga_db_disable_dummy_bitstream
 * @sa fpga_db_disable_dummy_bitstream_by_dev_id
 */
int fpga_db_get_child_bitstream_ids_by_dev_id(
        uint32_t dev_id,
        char **child_bsid_list[]);

/**
 * @brief Get available child bitstream ids by parent bitstream id
 * @param[in] parent_bsid
 *   Parent bitstream id
 * @param[out] child_bsid_list
 *   Child bitstream id list
 * @retval 0
 *   success
 * @retval -INVALID_ARGUMENT
 *   e.g.) `parent_bsid` is null
 * @retval -INVALID_DATA
 *   e.g.) data parameter is missing
 * @retval -FAILURE_MEMORY_ALLOC
 *   Failed to allocate memory
 * @return retval of fpga_get_dev_id()
 * @return retval of fpga_db_get_bitstream_id()
 *
 * @details
 *   Get all child bitstream ids for `parent_bsid`.@n
 *   `*child_bsid_list` and `(*child_bsid_list)[*]` will be allocated by this API,
 *    so please fpga_db_free_child_bitstream_ids() `child_bsid_list` explicitly.@n
 *   The sentinel of `*child_bsid_list` is NULL.
 * @remarks
 *   This API does not need FPGA.
 */
int fpga_db_get_child_bitstream_ids_by_parent(
        const char *parent_bsid,
        char **child_bsid_list[]);


/**
 * @brief Free buffer allocated by fpga_db_get_child_bitstream_ids***()
 * @param[in] child_bsid_list
 *   Child bitstream id list
 * @retval 0
 *   success
 * @retval -INVALID_ARGUMENT
 *   e.g.) `child_bsid_list` is null
 *
 * @details
 *   Free all `(*child_bsid_list)[*]` and `*child_bsid_list`.
 */
int fpga_db_free_child_bitstream_ids(
        char *child_bsid_list[]);


/**
 * @brief Set the Bitstream IDs handled by the APIs of this library to dummy values
 * @param[in] device_name
 *   FPGA's device path or serial id
 * @param[in] parent_dummy_bsid
 *   Dummy parent bitstream id
 * @param[in] child_dummy_bsid
 *   Dummy child bitstream id
 * @retval 0
 *   success
 * @retval -INVALID_ARGUMENT
 *   e.g.) `device_name` is null
 *
 * @details
 *   Set the dummy bitstream IDs for the target device.@n
 *   When you give `NULL` for `parent_dummy_bsid` or `child_dummy_bsid`,
 *    the APIs of this library will use the real value from the device,
 *    otherwise, they will use the given values for bitstream IDs.@n
 *   When both of `parent_dummy_bsid` and `child_dummy_bsid` are `NULL`,
 *    return error.
 */
int fpga_db_enable_dummy_bitstream(
        const char *device_name,
        const uint32_t *parent_dummy_bsid,
        const uint32_t *child_dummy_bsid);


/**
 * @brief Set the Bitstream IDs handled by the APIs of this library to dummy values
 * @param[in] dev_id
 *   device id obtained by fpga_dev_init()
 * @param[in] parent_dummy_bsid
 *   Dummy parent bitstream id
 * @param[in] child_dummy_bsid
 *   Dummy child bitstream id
 * @retval 0
 *   success
 * @retval -INVALID_ARGUMENT
 *   e.g.) both of `parent_dummy_bsid` and `child_dummy_bsid` are null
 *
 * @details
 *   Set the dummy bitstream IDs for the target device.@n
 *   When you give `NULL` for `parent_dummy_bsid` or `child_dummy_bsid`,
 *    the APIs of this library will use the real value from the device,
 *    otherwise, they will use the given values for bitstream IDs.@n
 *   When both of `parent_dummy_bsid` and `child_dummy_bsid` are `NULL`,
 *    return error.
 */
int fpga_db_enable_dummy_bitstream_by_dev_id(
        uint32_t dev_id,
        const uint32_t *parent_dummy_bsid,
        const uint32_t *child_dummy_bsid);


/**
 * @brief Unset the dummy values of fpga_db_enable_dummy_bitstream()
 * @param[in] device_name
 *   FPGA's device path or serial id
 * @retval 0
 *   success
 * @retval -INVALID_ARGUMENT
 *   e.g.) `device_name` is null
 *
 * @details
 *   Unset the dummy values for the target device
 *    set by fpga_db_enable_dummy_bitstream()
 *    or fpga_db_enable_dummy_bitstream_by_dev_id().@n
 *   This API will return success even if dummy values are not set.
 *    for the device.
 */
int fpga_db_disable_dummy_bitstream(
        const char *device_name);


/**
 * @brief Unset the dummy values of fpga_db_enable_dummy_bitstream()
 * @param[in] dev_id
 *   device id obtained by fpga_dev_init()
 * @retval 0
 *   success
 * @retval -INVALID_ARGUMENT
 *   e.g.) `device_name` is null
 *
 * @details
 *   Unset the dummy values for the target device
 *    set by fpga_db_enable_dummy_bitstream()
 *    or fpga_db_enable_dummy_bitstream_by_dev_id().@n
 *   This API will return success even if dummy values are not set.
 *    for the device.
 */
int fpga_db_disable_dummy_bitstream_by_dev_id(
        uint32_t dev_id);

#ifdef __cplusplus
}
#endif

#endif  // LIBFPGADB_INCLUDE_LIBFPGADB_H_
