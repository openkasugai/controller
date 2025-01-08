/*****************************************************************
 * Copyright 2024 NTT Corporation, FUJITSU LIMITED
 *****************************************************************/
#ifndef __COMMON_H__
#define __COMMON_H__

#include <stdint.h>
#include <stdbool.h>

//-----------------------------------------------------
// define
//-----------------------------------------------------
//--LOG LEVEL--
#define LOG_FORCE	7
#define LOG_ERROR	5	//default
#define LOG_WARN	4
#define LOG_INFO	3
#define LOG_DEBUG	2
#define LOG_TRACE	1
//--DATA PATTERN--
#define DATA_00		0
#define DATA_FF		1
#define DATA_INC	2
#define DATA_DEC	3
#define DATA_55		4
#define DATA_AA		5

//-----------------------------------------------------
// function
//-----------------------------------------------------
int getopt_loglevel(void);
int logfile(int level, const char * format, ...);
int rslt2file(const char * format, ...);
int rslt2fonly(const char * format, ...);
int init_data(uint8_t* p, uint32_t dsize, int type);
int dump_mem(bool tofile, const char* addr, const int size);
int make_dir(const char* dir);
bool check_file_exist(const char* file);
int get_file_num(const char* files_path);
int remove_file(const char* file);
int remove_files_path(const char* files_path);
char* bool2string (bool b);
int32_t stoi(const char *str, int64_t *result);
uint64_t time_duration(const struct timespec *t1, const struct timespec *t2);
uint32_t nextPow2(int32_t n);

#endif		/* __COMMON_H__ */
