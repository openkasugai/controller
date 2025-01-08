/*****************************************************************
 * Copyright 2024 NTT Corporation, FUJITSU LIMITED
 *****************************************************************/
#ifndef __GLUE_FUNC_H__
#define __GLUE_FUNC_H__

#include <stdint.h>
#include "glue.h"


//-----------------------------------------------------
// variable
//-----------------------------------------------------

//-----------------------------------------------------
// function
//-----------------------------------------------------
// main
extern int32_t glue(tcp_client_info_t tcp_client_info, char* connector_id, uint32_t width, uint32_t height);

// modules
extern int32_t glue_shmem_allocate(shmem_mode_t shmem_mode, mngque_t *pque, uint32_t width, uint32_t height);
extern int32_t glue_shmem_free(mngque_t *pque);
extern int32_t glue_allocate_buffer(void);
extern void glue_free_buffer(void);
extern int32_t glue_dequeue_lldma_queue_setup(uint32_t dev_id, char* connect_id);
extern void glue_dequeue_lldma_queue_finish(uint32_t dev_id);


#endif /* __GLUE_FUNC_H__ */
