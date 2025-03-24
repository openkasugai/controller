/*************************************************
* Copyright 2025 NTT Corporation , FUJITSU LIMITED
*************************************************/

#include <fcc_json.h>
#include <fcc_log.h>
#include <fcc_prm.h>

#include <parson.h>

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <errno.h>


static int __fcc_json_parse_params(
  JSON_Object *dataflow_parameter_object)
{
  int ret = 0;

  // param var
  const char *connector_id  = NULL;
  const char *device        = NULL;
  uint32_t lane             = -1;
  const char *direction_str = NULL;
  uint32_t direction        = FCC_PRM_DIR_BOTH;
  uint32_t fchid            = -1;
  const char *extif_id_str  = NULL;
  uint32_t extif_id         = -1;
  uint32_t cid              = -1;

  // Get connector_id
  if (json_object_has_value(dataflow_parameter_object, "connector_id")) {
    connector_id = json_object_get_string(dataflow_parameter_object, "connector_id");
  } else if (json_object_has_value(dataflow_parameter_object,  "connector-id")) {
    connector_id = json_object_get_string(dataflow_parameter_object, "connector-id");
  } else if (json_object_has_value(dataflow_parameter_object,  "matching-key")) {
    connector_id = json_object_get_string(dataflow_parameter_object, "matching-key");
  } else if (json_object_has_value(dataflow_parameter_object,  "matching_key")) {
    connector_id = json_object_get_string(dataflow_parameter_object, "matching_key");
  }

  // Get device
  if (json_object_has_value(dataflow_parameter_object, "device")) {
    device = json_object_get_string(dataflow_parameter_object, "device");
  }

  // Get lane
  if (json_object_has_value(dataflow_parameter_object, "lane")) {
    lane = json_object_get_number(dataflow_parameter_object, "lane");
  }

  // Get direction
  if (json_object_has_value(dataflow_parameter_object, "dir")) {
    direction_str = json_object_get_string(dataflow_parameter_object, "dir");
  } else if (json_object_has_value(dataflow_parameter_object,  "direction")) {
    direction_str = json_object_get_string(dataflow_parameter_object, "direction");
  }
  if (direction_str) {
    if (strcmp(direction_str, "ingress")      == 0) direction = FCC_PRM_DIR_INGR;
    else if (strcmp(direction_str, "ingr")    == 0) direction = FCC_PRM_DIR_INGR;
    else if (strcmp(direction_str, "Ingress") == 0) direction = FCC_PRM_DIR_INGR;
    else if (strcmp(direction_str, "Ingr")    == 0) direction = FCC_PRM_DIR_INGR;
    else if (strcmp(direction_str, "egress")  == 0) direction = FCC_PRM_DIR_EGR;
    else if (strcmp(direction_str, "egr")     == 0) direction = FCC_PRM_DIR_EGR;
    else if (strcmp(direction_str, "Egress")  == 0) direction = FCC_PRM_DIR_EGR;
    else if (strcmp(direction_str, "Egr")     == 0) direction = FCC_PRM_DIR_EGR;
    else if (strcmp(direction_str, "both")    == 0) direction = FCC_PRM_DIR_BOTH;
    else if (strcmp(direction_str, "BOTH")    == 0) direction = FCC_PRM_DIR_BOTH;
    else direction = atoi(direction_str);
  }

  // Get fchid
  if (json_object_has_value(dataflow_parameter_object, "fchid")) {
    fchid = json_object_get_number(dataflow_parameter_object, "fchid");
  } else if (json_object_has_value(dataflow_parameter_object, "function_channel_id")) {
    fchid = json_object_get_number(dataflow_parameter_object, "function_channel_id");
  } else if (json_object_has_value(dataflow_parameter_object, "function-channel-id")) {
    fchid = json_object_get_number(dataflow_parameter_object, "function-channel-id");
  } else if (json_object_has_value(dataflow_parameter_object, "matching_key")) {
    fchid = json_object_get_number(dataflow_parameter_object, "matching_key");
  }

  // Get extif_id
  if (json_object_has_value(dataflow_parameter_object, "extif_id")) {
    extif_id_str = json_object_get_string(dataflow_parameter_object, "extif_id");
  } else if (json_object_has_value(dataflow_parameter_object,  "extif-id")) {
    extif_id_str = json_object_get_string(dataflow_parameter_object, "extif-id");
  }
  if (extif_id_str) {
    if (strcmp(extif_id_str, "lldma")      == 0) extif_id = FCC_PRM_EXTIF_LLDMA;
    else if (strcmp(extif_id_str, "LLDMA") == 0) extif_id = FCC_PRM_EXTIF_LLDMA;
    else if (strcmp(extif_id_str, "ptu")   == 0) extif_id = FCC_PRM_EXTIF_PTU;
    else if (strcmp(extif_id_str, "PTU")   == 0) extif_id = FCC_PRM_EXTIF_PTU;
    else extif_id = atoi(extif_id_str);
  }

  // Get cid
  if (json_object_has_value(dataflow_parameter_object, "cid")) {
    cid = json_object_get_number(dataflow_parameter_object, "cid");
  } else if (json_object_has_value(dataflow_parameter_object, "connection_id")) {
    cid = json_object_get_number(dataflow_parameter_object, "connection_id");
  } else if (json_object_has_value(dataflow_parameter_object, "connection-id")) {
    cid = json_object_get_number(dataflow_parameter_object, "connection-id");
  }


  // Set lldma prm
  if (connector_id) {
    ret = fcc_prm_push_lldma_list(connector_id);
    if (ret) goto out;
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

out:
  return ret;
}


static int __fcc_json_parse_string(
  JSON_Value *root_value)
{
  int ret = -1;

  // JSON param var
  JSON_Object *root_object               = NULL;
  JSON_Array  *dataflow_parameters_array = NULL;

  // Parse root
  root_object = json_object(root_value);
  if (!root_object) {
    fcc_log_errorf(" ! Failed to get root object\n");
    goto out;
  }

  // Get dataflow_parameters[]
  dataflow_parameters_array = json_object_get_array(root_object, "dataflow-parameters");
  if (!dataflow_parameters_array) {
    fcc_log_errorf(" ! Failed to get dataflow-parameter array\n");
    goto out;
  }

  // Access dataflow_parameters[index] and get parameters
  for (int idx = 0; idx < json_array_get_count(dataflow_parameters_array); idx++) {
    ret = __fcc_json_parse_params(json_array_get_object(dataflow_parameters_array, idx));
    if (ret) {
      fcc_log_errorf(" ! Failed to parse dataflow-parameters\n");
      goto out;
    }
  }

  ret = 0;
out:

  return ret;
}


int fcc_json_parse_string(
  const char *json_string)
{
  JSON_Value *root_value = json_parse_string_with_comments(json_string);
  if (!root_value) {
    fcc_log_errorf(" ! Failed to parse json_string:\n%s\n", json_string);
    return -1;
  }

  return __fcc_json_parse_string(root_value);
}


int fcc_json_parse_file(
  const char *json_file_path)
{
  JSON_Value *root_value = json_parse_file_with_comments(json_file_path);
  if (!root_value) {
    fcc_log_errorf(" ! Failed to parse json_file: %s\n", json_file_path);
    return -1;
  }

  return __fcc_json_parse_string(root_value);
}


int fcc_json_create_output_file(
  void)
{
  int ret = -1;

  const char *output_file_path = fcc_prm_get_output_file_path();
  char *output_json_string = NULL;

  // param list
  const fcc_prm_lldma_t *lldma_list = fcc_prm_get_lldma_list();
  const fcc_prm_ptu_t   *ptu_list   = fcc_prm_get_ptu_list();
  const fcc_prm_chain_t *chain_list = fcc_prm_get_chain_list();

  // JSON param var
  JSON_Value  *root_value                = NULL;
  JSON_Object *root_object               = NULL;
  JSON_Value  *dataflow_parameters_value = NULL;
  JSON_Array  *dataflow_parameters_array = NULL;
  JSON_Value  *dataflow_parameter_value  = NULL;
  JSON_Object *dataflow_parameter_object = NULL;
  JSON_Status s;

  // Create a new output_json_file
  FILE *fp = fopen(output_file_path, "w");
  if (!fp) {
    fcc_log_errorf(" ! Failed to open %s: %s\n",
      output_file_path, strerror(errno));
    goto out;
  }

  // Initialize root
  root_value = json_value_init_object();
  if (!root_value) {
    fcc_log_errorf(" ! Failed to init root value\n");
    goto out;
  }
  root_object = json_value_get_object(root_value);
  if (!root_object) {
    fcc_log_errorf(" ! Failed to get root object\n");
    goto out;
  }

  // Initialize dataflow_parameters[]
  dataflow_parameters_value = json_value_init_array();
  if (!dataflow_parameters_value) {
    fcc_log_errorf(" ! Failed to init dataflow-parameters value\n");
    goto out;
  }
  dataflow_parameters_array = json_value_get_array(dataflow_parameters_value);
  if (!dataflow_parameters_array) {
    fcc_log_errorf(" ! Failed to get dataflow-parameter array\n");
    goto out;
  }

  // Append lldma parameters into dataflow_parameters[]
  for (const fcc_prm_index_t *ent = fcc_prm_get_err_lldma_list(); ent && *ent && lldma_list; ent++) {
    dataflow_parameter_value = json_value_init_object();
    if (!dataflow_parameter_value) {
      fcc_log_errorf(" ! Failed to init dataflow-parameter value\n");
      goto out;
    }
    dataflow_parameter_object = json_value_get_object(dataflow_parameter_value);
    if (!dataflow_parameter_object) {
      fcc_log_errorf(" ! Failed to get dataflow-parameter object\n");
      goto out;
    }
    s = json_object_set_string(dataflow_parameter_object, "connector_id", lldma_list[**ent]);
    if (s) {
      fcc_log_errorf(" ! Failed to set dataflow-parameter string(connector_id)\n");
      goto out;
    }
    s = json_array_append_value(dataflow_parameters_array, dataflow_parameter_value);
    if (s) {
      fcc_log_errorf(" ! Failed to append dataflow-parameter value\n");
      goto out;
    }
  }

  // Append ptu parameters into dataflow_parameters[]
  for (const fcc_prm_index_t *ent = fcc_prm_get_err_ptu_list(); ent && *ent && ptu_list; ent++) {
    dataflow_parameter_value = json_value_init_object();
    if (!dataflow_parameter_value) {
      fcc_log_errorf(" ! Failed to init dataflow-parameter value\n");
      goto out;
    }
    dataflow_parameter_object = json_value_get_object(dataflow_parameter_value);
    if (!dataflow_parameter_object) {
      fcc_log_errorf(" ! Failed to get dataflow-parameter object\n");
      goto out;
    }
    s = json_object_set_string(dataflow_parameter_object, "device", ptu_list[**ent]->device);
    if (s) {
      fcc_log_errorf(" ! Failed to set dataflow-parameter string(device)\n");
      goto out;
    }
    s = json_object_set_number(dataflow_parameter_object, "lane", ptu_list[**ent]->lane);
    if (s) {
      fcc_log_errorf(" ! Failed to set dataflow-parameter string(lane)\n");
      goto out;
    }
    s = json_object_set_string(dataflow_parameter_object, "extif_id",
      FCC_PRM_EXTIF_TO_STR(ptu_list[**ent]->extif_id));
    if (s) {
      fcc_log_errorf(" ! Failed to set dataflow-parameter string(extif_id)\n");
      goto out;
    }
    s = json_object_set_number(dataflow_parameter_object, "connection_id", ptu_list[**ent]->cid);
    if (s) {
      fcc_log_errorf(" ! Failed to set dataflow-parameter string(connection_id)\n");
      goto out;
    }
    s = json_array_append_value(dataflow_parameters_array, dataflow_parameter_value);
    if (s) {
      fcc_log_errorf(" ! Failed to append dataflow-parameter value\n");
      goto out;
    }
  }

  // Append chain parameters into dataflow_parameters[]
  for (const fcc_prm_index_t *ent = fcc_prm_get_err_chain_list(); ent && *ent && chain_list; ent++) {
    dataflow_parameter_value = json_value_init_object();
    if (!dataflow_parameter_value) {
      fcc_log_errorf(" ! Failed to init dataflow-parameter value\n");
      goto out;
    }
    dataflow_parameter_object = json_value_get_object(dataflow_parameter_value);
    if (!dataflow_parameter_object) {
      fcc_log_errorf(" ! Failed to get dataflow-parameter object\n");
      goto out;
    }
    s = json_object_set_string(dataflow_parameter_object, "device", chain_list[**ent]->device);
    if (s) {
      fcc_log_errorf(" ! Failed to set dataflow-parameter string(device)\n");
      goto out;
    }
    s = json_object_set_number(dataflow_parameter_object, "lane", chain_list[**ent]->lane);
    if (s) {
      fcc_log_errorf(" ! Failed to set dataflow-parameter string(lane)\n");
      goto out;
    }
    s = json_object_set_string(dataflow_parameter_object, "direction",
      FCC_PRM_DIR_TO_STR(chain_list[**ent]->dir));
    if (s) {
      fcc_log_errorf(" ! Failed to set dataflow-parameter string(direction)\n");
      goto out;
    }
    s = json_object_set_number(dataflow_parameter_object, "function_channel_id", chain_list[**ent]->fchid);
    if (s) {
      fcc_log_errorf(" ! Failed to set dataflow-parameter string(function_channel_id)\n");
      goto out;
    }
    s = json_array_append_value(dataflow_parameters_array, dataflow_parameter_value);
    if (s) {
      fcc_log_errorf(" ! Failed to append dataflow-parameter value\n");
      goto out;
    }
  }

  // Avoid double free
  dataflow_parameter_value = NULL;

  // Set dataflow_parameters[] into root
  s = json_object_set_value(root_object, "dataflow-parameters", dataflow_parameters_value);
  if (s) {
    fcc_log_errorf(" ! Failed to set dataflow-parameters value\n");
    goto out;
  }
  dataflow_parameters_value = NULL;

  // Convert root to string
  json_set_escape_slashes(0);
  output_json_string = json_serialize_to_string_pretty(root_value);

  fprintf(fp, "%s\n", output_json_string);

  ret = 0;

out:
  if (output_json_string)
    json_free_serialized_string(output_json_string);
  if (root_value)
    json_value_free(root_value);
  if (dataflow_parameters_value)
    json_value_free(dataflow_parameters_value);
  if (dataflow_parameter_value)
    json_value_free(dataflow_parameter_value);

  if (fp)
    fclose(fp);

  return ret;
}
