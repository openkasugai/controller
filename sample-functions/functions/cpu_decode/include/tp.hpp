/*****************************************************************
 * Copyright 2024 NTT Corporation, FUJITSU LIMITED
 *****************************************************************/
#ifndef TP_HPP
#define TP_HPP

extern int32_t tp_function_filter_resize_init(uint32_t dev_id);
extern void tp_function_finish(uint32_t dev_id);
extern int32_t tp_enqueue_lldma_init(uint32_t dev_id);
extern int32_t tp_dequeue_lldma_init(uint32_t dev_id);
extern void tp_enqueue_lldma_finish(void);
extern void tp_dequeue_lldma_finish(void);
extern int32_t tp_enqueue_lldma_queue_setup(void);
extern int32_t tp_dequeue_lldma_queue_setup(void);
extern void tp_enqueue_lldma_queue_finish(void);
extern void tp_dequeue_lldma_queue_finish(void);
extern int32_t tp_chain_connect(uint32_t dev_id);
extern int32_t tp_enqueue_set_dma_cmd(uint32_t enq_id, const mngque_t &mngq, dmacmd_info_t &enqdmacmdinfo);
extern int32_t tp_dequeue_set_dma_cmd(uint32_t enq_id, const mngque_t &mngq, dmacmd_info_t &deqdmacmdinfo);
extern int32_t wait_dma_fpga_enqueue(dma_info_t &dmainfo, dmacmd_info_t &dmacmdinfo, uint32_t enq_id, const uint32_t msec);
extern int32_t wait_dma_fpga_dequeue(dma_info_t &dmainfo, dmacmd_info_t &dmacmdinfo, uint32_t enq_id, const uint32_t msec);

#endif //TP_HPP
