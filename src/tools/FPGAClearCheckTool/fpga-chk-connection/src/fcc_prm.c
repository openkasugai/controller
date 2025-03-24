/*************************************************
* Copyright 2025 NTT Corporation , FUJITSU LIMITED
*************************************************/

#include <fcc_prm.h>
#include <fcc_log.h>

#include <stddef.h>
#include <stdlib.h>
#include <string.h>


static fcc_prm_lldma_t *fcc_lldma_list = NULL;
static fcc_prm_ptu_t   *fcc_ptu_list   = NULL;
static fcc_prm_chain_t *fcc_chain_list = NULL;

static fcc_prm_index_t *fcc_err_lldma_list = NULL;
static fcc_prm_index_t *fcc_err_ptu_list   = NULL;
static fcc_prm_index_t *fcc_err_chain_list = NULL;

static bool fcc_is_dump = false;

static char *fcc_output_file_path = NULL;


const bool fcc_prm_get_is_dump(
  void)
{
  return fcc_is_dump;
}

void fcc_prm_set_is_dump(
  bool is_dump)
{
  fcc_is_dump = is_dump;
}


const char *fcc_prm_get_output_file_path(
  void)
{
  return fcc_output_file_path ? fcc_output_file_path
                              : FCC_PRM_OUTPUT_FILE_PATH_DEFAULT;
}

int fcc_prm_set_output_file_path(
  const char *output_file_path)
{
  if (fcc_output_file_path)
    free(fcc_output_file_path);

  fcc_output_file_path = strdup(output_file_path);

  if (fcc_output_file_path)
    return 0;
  return -1;
}

void fcc_prm_free_output_file_path(
  void)
{
  if (fcc_output_file_path)
    free(fcc_output_file_path);
  fcc_output_file_path = NULL;
}


const fcc_prm_lldma_t *fcc_prm_get_lldma_list(
  void)
{
  return fcc_lldma_list;
}

const fcc_prm_ptu_t *fcc_prm_get_ptu_list(
  void)
{
  return fcc_ptu_list;
}

const fcc_prm_chain_t *fcc_prm_get_chain_list(
  void)
{
  return fcc_chain_list;
}

int fcc_prm_get_lldma_list_size(
  void)
{
  int size = 0;
  for (const fcc_prm_lldma_t *ent = fcc_lldma_list; ent && *ent; ent++)
    size++;
  return size;
}

int fcc_prm_get_ptu_list_size(
  void)
{
  int size = 0;
  for (const fcc_prm_ptu_t *ent = fcc_ptu_list; ent && *ent; ent++)
    size++;
  return size;
}

int fcc_prm_get_chain_list_size(
  void)
{
  int size = 0;
  for (const fcc_prm_chain_t *ent = fcc_chain_list; ent && *ent; ent++)
    size++;
  return size;
}

const fcc_prm_index_t *fcc_prm_get_err_lldma_list(
  void)
{
  return fcc_err_lldma_list;
}

const fcc_prm_index_t *fcc_prm_get_err_ptu_list(
  void)
{
  return fcc_err_ptu_list;
}

const fcc_prm_index_t *fcc_prm_get_err_chain_list(
  void)
{
  return fcc_err_chain_list;
}

void fcc_prm_free_lldma_list(
  void)
{
  if (!fcc_lldma_list)
    return;
  for (fcc_prm_lldma_t *e = fcc_lldma_list; *e; e++) {
    free(*e);
  }
  free(fcc_lldma_list);
  fcc_lldma_list = NULL;
}

void fcc_prm_free_ptu_list(
  void)
{
  if (!fcc_ptu_list)
    return;
  for (fcc_prm_ptu_t *e = fcc_ptu_list; *e; e++) {
    free((*e)->device);
    free(*e);
  }
  free(fcc_ptu_list);
  fcc_ptu_list = NULL;
}

void fcc_prm_free_chain_list(
  void)
{
  if (!fcc_chain_list)
    return;
  for (fcc_prm_chain_t *e = fcc_chain_list; *e; e++) {
    free((*e)->device);
    free(*e);
  }
  free(fcc_chain_list);
  fcc_chain_list = NULL;
}

static void __fcc_prm_free_index_list(
  fcc_prm_index_t **plist)
{
  if (!plist) return;
  fcc_prm_index_t *list = *plist;
  if (!list) return;
  for (fcc_prm_index_t *e = list; *e; e++)
    free(*e);
  free(list);
  *plist = NULL;
}

void fcc_prm_free_err_lldma_list(
  void)
{
  __fcc_prm_free_index_list(&fcc_err_lldma_list);
}

void fcc_prm_free_err_ptu_list(
  void)
{
  __fcc_prm_free_index_list(&fcc_err_ptu_list);
}

void fcc_prm_free_err_chain_list(
  void)
{
  __fcc_prm_free_index_list(&fcc_err_chain_list);
}


int fcc_prm_push_lldma_list(
  const char *connector_id)
{
  // Check input
  if (!connector_id) {
    fcc_log_errorf(" ! connector_id is NULL for %s\n", __func__);
    return -1;
  }

  // Dup data
  fcc_prm_lldma_t data = strdup(connector_id);
  if (!data) {
    fcc_log_errorf(" ! Failed to allocate data for %s\n", __func__);
    return -2;
  }

  // Get old list size
  int size_old = 0;
  if (fcc_lldma_list)
    for(; fcc_lldma_list[size_old]; size_old++)
      ; // do nothing

  // Allocate new list
  fcc_prm_lldma_t *list_new = (fcc_prm_lldma_t*)malloc(sizeof(fcc_prm_lldma_t) * (size_old + 2));
  if (!list_new) {
    fcc_log_errorf(" ! Failed to allocate list for %s\n", __func__);
    free(data);
    return -2;
  }

  // Copy into new list
  if (fcc_lldma_list)
    for(int index = 0; index < size_old; index++)
      list_new[index] = fcc_lldma_list[index];

  // Add new data
  list_new[size_old]     = data;
  list_new[size_old + 1] = NULL;

  // Delete old list
  if (fcc_lldma_list)
    free(fcc_lldma_list);

  // Set new list
  fcc_lldma_list = list_new;

  return 0;
}

int fcc_prm_push_ptu_list(
  const char *device,
  uint32_t lane,
  uint32_t extif_id,
  uint32_t cid)
{
  // Check input
  if (!device) {
    fcc_log_errorf(" ! device is NULL for %s\n", __func__);
    return -1;
  }

  // Dup data
  fcc_prm_ptu_t data = NULL;
  data = (fcc_prm_ptu_t)malloc(sizeof(*data));
  if (!data) {
    fcc_log_errorf(" ! Failed to allocate data for %s\n", __func__);
    return -2;
  }
  data->lane     = lane;
  data->extif_id = extif_id;
  data->cid      = cid;
  data->device   = strdup(device);
  if (!data->device) {
    fcc_log_errorf(" ! Failed to allocate data(device) for %s\n", __func__);
    free(data);
    return -2;
  }

  // Get old list size
  int size_old = 0;
  if (fcc_ptu_list)
    for(; fcc_ptu_list[size_old]; size_old++)
      ; // do nothing

  // Allocate new list
  fcc_prm_ptu_t *list_new = (fcc_prm_ptu_t*)malloc(sizeof(fcc_prm_ptu_t) * (size_old + 2));
  if (!list_new) {
    fcc_log_errorf(" ! Failed to allocate list for %s\n", __func__);
    free(data->device);
    free(data);
    return -2;
  }

  // Copy into new list
  if (fcc_ptu_list)
    for(int index = 0; index < size_old; index++)
      list_new[index] = fcc_ptu_list[index];

  // Add new data
  list_new[size_old]     = data;
  list_new[size_old + 1] = NULL;

  // Delete old list
  if (fcc_ptu_list)
    free(fcc_ptu_list);

  // Set new list
  fcc_ptu_list = list_new;

  return 0;
}

int fcc_prm_push_chain_list(
  const char *device,
  uint32_t lane,
  uint32_t fchid,
  uint32_t dir)
{
  // Check input
  if (!device) {
    fcc_log_errorf(" ! device is NULL for %s\n", __func__);
    return -1;
  }

  // Dup data
  fcc_prm_chain_t data = NULL;
  data = (fcc_prm_chain_t)malloc(sizeof(*data));
  if (!data) {
    fcc_log_errorf(" ! Failed to allocate data for %s\n", __func__);
    return -2;
  }
  data->lane   = lane;
  data->fchid  = fchid;
  data->dir    = dir;
  data->device = strdup(device);
  if (!data->device) {
    fcc_log_errorf(" ! Failed to allocate data(device) for %s\n", __func__);
    free(data);
    return -2;
  }

  // Get old list size
  int size_old = 0;
  if (fcc_chain_list)
    for(; fcc_chain_list[size_old]; size_old++)
      ; // do nothing

  // Allocate new list
  fcc_prm_chain_t *list_new = (fcc_prm_chain_t*)malloc(sizeof(fcc_prm_chain_t) * (size_old + 2));
  if (!list_new) {
    fcc_log_errorf(" ! Failed to allocate list for %s\n", __func__);
    free(data->device);
    free(data);
    return -2;
  }

  // Copy into new list
  if (fcc_chain_list)
    for(int index = 0; index < size_old; index++)
      list_new[index] = fcc_chain_list[index];

  // Add new data
  list_new[size_old]     = data;
  list_new[size_old + 1] = NULL;

  // Delete old list
  if (fcc_chain_list)
    free(fcc_chain_list);

  // Set new list
  fcc_chain_list = list_new;

  return 0;
}


static int __fcc_prm_push_index_list(
  fcc_prm_index_t **plist,
  int input_data)
{
  if (!plist) {
    fcc_log_errorf(" ! List is invalid in %s\n", __func__);
    return -2;
  }

  fcc_prm_index_t *list = *plist;

    // Dup data
  fcc_prm_index_t data = NULL;
  data = (fcc_prm_index_t)malloc(sizeof(*data));
  if (!data) {
    fcc_log_errorf(" ! Failed to allocate data for %s\n", __func__);
    return -2;
  }
  *data = input_data;

  // Get old list size
  int size_old = 0;
  if (list)
    for(; list[size_old]; size_old++)
      ; // do nothing

  // Allocate new list
  fcc_prm_index_t *list_new = (fcc_prm_index_t*)malloc(sizeof(fcc_prm_index_t*) * (size_old + 2));
  if (!list_new) {
    fcc_log_errorf(" ! Failed to allocate list for %s\n", __func__);
    free(data);
    return -2;
  }

  // Copy into new list
  if (list)
    for(int index = 0; index < size_old; index++)
      list_new[index] = list[index];

  // Add new data
  list_new[size_old]     = data;
  list_new[size_old + 1] = NULL;

  // Delete old list
  if (list)
    free(list);

  // Set new list
  *plist = list_new;

  return 0;
}

int fcc_prm_push_err_lldma_list(
  int input_data)
{
  return __fcc_prm_push_index_list(&fcc_err_lldma_list, input_data);
}

int fcc_prm_push_err_ptu_list(
  int input_data)
{
  return __fcc_prm_push_index_list(&fcc_err_ptu_list, input_data);
}

int fcc_prm_push_err_chain_list(
  int input_data)
{
  return __fcc_prm_push_index_list(&fcc_err_chain_list, input_data);
}


void fcc_prm_show(
  void)
{
  for (const fcc_prm_lldma_t *head = fcc_prm_get_lldma_list(); head && *head; head++) {
    fcc_log_printf(" LLDMA : connector_id(%s)\n", *head);
  }
  for (const fcc_prm_ptu_t *head = fcc_prm_get_ptu_list(); head && *head; head++) {
    fcc_log_printf(" PTU   : device(%s),lane(%u),extif_id(%u),cid(%u)\n",
      (*head)->device, (*head)->lane, (*head)->extif_id, (*head)->cid);
  }
  for (const fcc_prm_chain_t *head = fcc_prm_get_chain_list(); head && *head; head++) {
    fcc_log_printf(" CHAIN : device(%s),lane(%u),fchid(%u),dir(%u)\n",
      (*head)->device, (*head)->lane, (*head)->fchid, (*head)->dir);
  }
  for (const fcc_prm_index_t *head = fcc_prm_get_err_lldma_list(); head && *head; head++) {
    fcc_log_printf(" ERRSSSN : Index(%u)\n", **head);
  }
  for (const fcc_prm_index_t *head = fcc_prm_get_err_chain_list(); head && *head; head++) {
    fcc_log_printf(" ERRCHAIN: Index(%u)\n", **head);
  }
  fcc_log_printf(" DUMP  : %s\n", fcc_prm_get_is_dump() ? "Yes" : "No");
  fcc_log_printf(" OUTPUT: %s\n", fcc_prm_get_output_file_path());
}
