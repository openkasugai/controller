/*************************************************
* Copyright 2025 NTT Corporation , FUJITSU LIMITED
*************************************************/
#ifndef FCC_LOG_H_
#define FCC_LOG_H_

#include <stdio.h>
#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif


#define __fcc_log(fp,...)\
  do {\
    fprintf((fp), __VA_ARGS__);\
    fflush((fp));\
  } while (0)

#define fcc_log_printf(...)\
  __fcc_log(stdout, __VA_ARGS__)

/**
 * @remarks
 *  This cannot use argument which calculate
 *   because calculation will be executed twice.
 */
#define fcc_log_errorf(...)\
  do {\
    __fcc_log(stderr, __VA_ARGS__);\
    __fcc_log(stdout, __VA_ARGS__);\
  } while (0)


uint32_t fcc_log_get_app_version(
        void);
uint32_t fcc_log_get_major_version(
        void);
uint32_t fcc_log_get_minor_version(
        void);
uint32_t fcc_log_get_revision(
        void);
uint32_t fcc_log_get_patch(
        void);


#ifdef __cplusplus
}
#endif

#endif  // FCC_LOG_H_
