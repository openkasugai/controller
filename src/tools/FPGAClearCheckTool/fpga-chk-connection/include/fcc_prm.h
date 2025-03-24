/*************************************************
* Copyright 2025 NTT Corporation , FUJITSU LIMITED
*************************************************/
#ifndef FCC_PRM_H_
#define FCC_PRM_H_

#include <stdint.h>
#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif


#define FCC_PRM_EXTIF_LLDMA 0
#define FCC_PRM_EXTIF_PTU   1
#define FCC_PRM_EXTIF_TO_STR(id) ((id) == FCC_PRM_EXTIF_LLDMA \
                                              ? "LLDMA" \
                                              : "PTU")

#define FCC_PRM_DIR_INGR    0b01
#define FCC_PRM_DIR_EGR     0b10
#define FCC_PRM_DIR_BOTH    0b11
#define FCC_PRM_DIR_TO_STR(id) \
           ( (id) == FCC_PRM_DIR_INGR ? "Ingress" \
           : (id) == FCC_PRM_DIR_EGR  ? "Egress"  \
                                         : "Both")

#define FCC_PRM_OUTPUT_FILE_PATH_DEFAULT "df-found.json"

#define FCC_PRM_ERRNO_DUMP 254
#define FCC_PRM_ERRNO_HELP 255

struct fcc_prm_ptu {
  char *device;
  uint32_t lane;
  uint32_t extif_id;
  uint32_t cid;
};

struct fcc_prm_chain {
  char *device;
  uint32_t lane;
  uint32_t dir;
  uint32_t fchid;
};


typedef char*                 fcc_prm_lldma_t;
typedef struct fcc_prm_ptu*   fcc_prm_ptu_t;
typedef struct fcc_prm_chain* fcc_prm_chain_t;

typedef int*                  fcc_prm_index_t;


const bool fcc_prm_get_is_dump(
        void);
void fcc_prm_set_is_dump(
        bool is_dump);

const char *fcc_prm_get_output_file_path(
        void);
int fcc_prm_set_output_file_path(
        const char *output_file_path);
void fcc_prm_free_output_file_path(
        void);

int fcc_prm_push_lldma_list(
        const char *connector_id);
const fcc_prm_lldma_t *fcc_prm_get_lldma_list(
        void);
void fcc_prm_free_lldma_list(
        void);
int fcc_prm_get_lldma_list_size(
        void);

int fcc_prm_push_ptu_list(
        const char *device,
        uint32_t lane,
        uint32_t extif_id,
        uint32_t cid);
const fcc_prm_ptu_t *fcc_prm_get_ptu_list(
        void);
void fcc_prm_free_ptu_list(
        void);
int fcc_prm_get_ptu_list_size(
        void);

int fcc_prm_push_chain_list(
        const char *device,
        uint32_t lane,
        uint32_t fchid,
        uint32_t dir);
const fcc_prm_chain_t *fcc_prm_get_chain_list(
        void);
void fcc_prm_free_chain_list(
        void);
int fcc_prm_get_chain_list_size(
        void);

int fcc_prm_push_err_lldma_list(
        int input_data);
const fcc_prm_index_t *fcc_prm_get_err_lldma_list(
        void);
void fcc_prm_free_err_lldma_list(
        void);

int fcc_prm_push_err_ptu_list(
        int input_data);
const fcc_prm_index_t *fcc_prm_get_err_ptu_list(
        void);
void fcc_prm_free_err_ptu_list(
        void);

int fcc_prm_push_err_chain_list(
        int input_data);
const fcc_prm_index_t *fcc_prm_get_err_chain_list(
        void);
void fcc_prm_free_err_chain_list(
        void);

void fcc_prm_show(
        void);


#ifdef __cplusplus
}
#endif

#endif  // FCC_PRM_H_
