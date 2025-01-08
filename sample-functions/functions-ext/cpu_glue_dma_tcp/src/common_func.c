/*****************************************************************
 * Copyright 2024 NTT Corporation, FUJITSU LIMITED
 *****************************************************************/
#include <stdint.h>
#include <unistd.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>
#include <stdbool.h>
#include <stdarg.h>
#include <sys/stat.h>
#include <errno.h>
#include "common.h"


//-----------------------------
//  common func
//-----------------------------
//--- logfile ---
#define LOGFILE "applog"
int logfile(int level, const char * format, ...)
{
	int loglevel;
	char date[32];
	FILE *fp;
	time_t t;
	va_list ap;
	static char file[64];
	static int sta=0;

	loglevel = getopt_loglevel();
	if (level < loglevel) {
		return 0;
	}
	if (sta == 0) {
		time(&t);
		strftime(date, sizeof(date), "%m%d-%H%M%S", localtime(&t) );
		sprintf(file, "%s%s_L%d.log", LOGFILE, date, loglevel );
		printf("logfile= %s\n", file);
		fp = fopen(file, "a");
		if (fp == NULL) {
			return -1;
		}
		fprintf(fp, "log start...%s, (%s)\n", date, file);
		sta = 1;
	} else {
		fp = fopen(file, "a");
		if (fp == NULL) {
			return -1;
		}
	}
	switch( level ){
		case LOG_FORCE:
		case LOG_TRACE:
			break;
		case LOG_ERROR:
			fprintf(fp, "[error] ");
			break;
		case LOG_WARN:
			fprintf(fp, "[warn] ");
			break;
		case LOG_INFO:
			fprintf(fp, "[info] ");
			break;
		case LOG_DEBUG:
			fprintf(fp, "[debug] ");
			break;
		default:
			fprintf(fp, "[?????] ");
			break;
	}
	va_start(ap, format);
	vfprintf(fp, format, ap);
	va_end( ap );
	if (level == LOG_FORCE) {
		va_start(ap, format);
		vprintf(format, ap);
		va_end( ap );
	}
	fclose(fp);
	return 0;
}

//-----------------------------
// Result Output
//-----------------------------
#define RESULTFILE "result.log"
static int rsltstart=0;
int rslt2file(const char * format, ...)
{
        FILE* fp;
        time_t t;
        va_list ap;
        char date[32];

        time(&t);
        strftime(date, sizeof(date), "%m%d-%H%M%S: ", localtime(&t) );
        if (rsltstart == 0){
                fp = fopen(RESULTFILE, "w");
        } else {
                fp = fopen(RESULTFILE, "a");
        }
        if (fp == NULL) {
                return -1;
        }
        //to file
        if (rsltstart == 0){
                fprintf(fp, "result log start...%s\n", date);
                rsltstart = 1;
        }
        va_start(ap, format);
        vfprintf(fp, format, ap);
        va_end( ap );
        fclose(fp);
        //to stdin
        va_start(ap, format);
        vprintf(format, ap);
        va_end( ap );

        return 0;
}
int rslt2fonly(const char * format, ...)
{
        FILE* fp;
        time_t t;
        va_list ap;
        char date[32];

        time(&t);
        strftime(date, sizeof(date), "%m%d-%H%M%S: ", localtime(&t) );
        if (rsltstart == 0){
                fp = fopen(RESULTFILE, "w");
        } else {
                fp = fopen(RESULTFILE, "a");
        }
        if (fp == NULL) {
                return -1;
        }
        //to file
        if (rsltstart == 0){
                fprintf(fp, "result log start...%s\n", date);
                rsltstart = 1;
        }
        va_start(ap, format);
        vfprintf(fp, format, ap);
        va_end( ap );
        fclose(fp);

        return 0;
}
 
//------------------------------
// initialize data
//------------------------------
// type 0:all 0
// type 1:incriment
// type 2:decriment
// type 3:all 1
int init_data(uint8_t* p, uint32_t dsize, int type)
{
        int i;

        logfile(LOG_DEBUG, "init_data..(%p,%d,%d)\n", p, dsize, type);
        switch(type){
                case 0: //0x00
                        for( i=0; i<dsize; i++){
                                *((uint8_t*)p+i) = 0x00;
                        }
                        break;
                case 1: //0xff
                        for( i=0; i<dsize; i++){
                                *(p+i) = 0xff;
                        }
                        break;
                case 2: //inc
                        for( i=0; i<dsize; i++){
                                *(p+i) = (uint8_t)(i+1);
                        }
                        break;
                case 3: //dec
                        for( i=0; i<dsize; i++){
                                *(p+i) = 0xff - (uint8_t)i;
                        }
                        break;
                case 4: //0x55
                        for( i=0; i<dsize; i++){
                                *(p+i) = 0x55;
                        }
                        break;
                case 5: //0xaa
                        for( i=0; i<dsize; i++){
                                *(p+i) = 0xaa;
                        }
                        break;
                default: //0x11
                        for( i=0; i<dsize; i++){
                                *(p+i) = 0x11;
                        }
                        break;
        }
        //initialize data
        logfile(LOG_TRACE, " %p, %d  ", p, dsize);
        for( i=0; i<1024; i++){
                logfile(LOG_TRACE, "  %02x", *(p+i));
        }
        logfile(LOG_TRACE, "\n");
        return 0;
}

//---------------------------------------
// dump memory
//---------------------------------------
int dump_mem(bool tofile, const char* addr, const int size)
{
	int i;
	int ss = size;
	char* buff1;
	char buff2[1024];

	buff1 = (char*)malloc(size*4);
	if (buff1 == NULL) {
		printf("dump_mem malloc error!(%d)\n",size);
		return -1;
	}
	memset(buff1, 0, size);
	memset(buff2, 0, sizeof(buff2));
	if (addr == NULL) {
		free(buff1);
		return -1;
	}
	if (size > 1024) {
		ss = size;
	}
	for (i=0; i < ss; i++) {
		if ( i!=0 && (i%16)==0 ){
			sprintf(buff2, "\n");
			strcat(buff1, buff2);
		}
		sprintf(buff2, "%02x ", *(uint8_t*)(addr+i));
		strcat(buff1, buff2);

	}
	sprintf(buff2, "\n");
	strcat(buff1, buff2);
	if( tofile == true ) {
		rslt2file("%s", buff1);
		printf("%s", buff1);
	} else {
		logfile(LOG_DEBUG, "%s", buff1);
	}
	free(buff1);
	return 0;
}

//---------------------------------------
// make dir
//---------------------------------------
int make_dir(const char* dir)
{
	int32_t ret = 0;
	struct stat st;

	mode_t mode = S_IRUSR|S_IWUSR|S_IXUSR|S_IRGRP|S_IWGRP|S_IXGRP|S_IROTH|S_IXOTH ;
	if (stat(dir, &st) != 0) {
        	logfile(LOG_DEBUG, "make_dir..(%s)\n", dir);
		ret = mkdir(dir, mode);
		if (ret != 0) {
			rslt2file("make_dir error: can't make dir \"%s\"\n", dir);
        		logfile(LOG_ERROR, "make_dir error: can't make dir \"%s\"\n", dir);
		}
	}

	return ret;
}

//---------------------------------------
// check file exist
//---------------------------------------
bool check_file_exist(const char* file)
{
	bool ret = false;
	struct stat st;

	if (stat(file, &st) == 0 ) {
		ret = true;
	}

	return ret;
}

//---------------------------------------
// get file num
//---------------------------------------
int get_file_num(const char* files_path)
{
	int ret = 0;

	FILE *fp;
	char cmd[256];
	sprintf(cmd, "ls -U1 %s 2> /dev/null |wc -l ", files_path);
	fp = popen(cmd, "r");
	if (fp == NULL) {
		rslt2file("get_file_num error.");
		logfile(LOG_ERROR, "get_file_num error");
		return -1;
	}

	char buf[16], *str;
	if ((str = fgets(buf, sizeof(buf), fp)) != NULL) {
		ret = atoi(str);
	}

	(void)pclose(fp);

	return ret;
}

//---------------------------------------
// remove file
//---------------------------------------
int remove_file(const char* file)
{
	int32_t ret = 0;

	logfile(LOG_DEBUG, "remove_file..(%s)\n", file);

	ret = remove(file);
	if (ret != 0) {
		rslt2file("remove_file error: can't remove \"%s\"\n", file);
		logfile(LOG_ERROR, "remove_file error: can't remove \"%s\"\n", file);
	}

	return ret;
}

//---------------------------------------
// remove files
//---------------------------------------
int remove_files_path(const char* files_path)
{
	int32_t before_num = get_file_num(files_path);
	if (before_num > 0) {
		char cmd[256];
		sprintf(cmd, "rm %s", files_path);
		int32_t rs = system(cmd);
		if (WIFEXITED(rs) == 0) {
			rslt2file("remove_files_path error: failed external command \"%s\"\n", cmd);
			logfile(LOG_ERROR, "remove_files_path error: failed external command \"%s\"\n", cmd);
			return -1;
		}

		int32_t after_num = get_file_num(files_path);
		if (after_num != 0) {
			rslt2file("remove_files_path error: can't remove \"%s\"\n", files_path);
			logfile(LOG_ERROR, "remove_files_path error: can't remove \"%s\"\n", files_path);
			return -1;
		}
	}

	return 0;
}

//---------------------------------------
// convert bool to string
//---------------------------------------
char* bool2string (bool b)
{
        return b == true ? "true" : "false";
}

//---------------------------------------
// convert string to an integer
//---------------------------------------
#if 0
int32_t stoi(const char *str, int64_t *result)
{
	char *end;
	errno = 0;
	int64_t num = strtol(str, &end, 10);
	if (errno == ERANGE) {
		rslt2file("stoi error: errno == ERANGE \"%s\"\n", str);
		logfile(LOG_ERROR, "stoi error: errno == ERAGE \"%s\"\n", str);
		return -1;
	} else if (str == end) {
		rslt2file("stoi error: can't convert \"%s\" to an integer.\n", str);
		logfile(LOG_ERROR, "stoi error: can't convert \"%s\" to an integer.\n", str);
		return -1;
	}

	*result = num;

	return 0;
}
#else
int32_t stoi(const char *str, int64_t *result)
{
	char *end;
	errno = 0;

	bool first = true;
	// Check whether integer conversion is possible one character at a time
	for (size_t i = 0; i < strlen(str); i++) {
		if (str[i] == '-' && first) {
			// Leading hyphen is sign
			first = false;
			continue;
		}
		char tmp[2];
		sprintf(tmp, "%c", str[i]);
		strtol(tmp, &end, 10);
		if (errno == ERANGE) {
			rslt2file("stoi error: errno == ERANGE \"%s\"\n", str);
			logfile(LOG_ERROR, "stoi error: errno == ERAGE \"%s\"\n", str);
			return -1;
		} else if (tmp == end) {
			rslt2file("stoi error: Cannot to convert string \"%s\" to an integer.\n", str);
			logfile(LOG_ERROR, "stoi error: Cannot to convert string \"%s\" to an integer.\n", str);
			return -1;
		}
	}

	// Convert string to integer
	int64_t num = strtol(str, &end, 10);
	if (errno == ERANGE) {
		rslt2file("stoi error: errno == ERANGE \"%s\"\n", str);
		logfile(LOG_ERROR, "stoi error: errno == ERAGE \"%s\"\n", str);
		return -1;
	} else if (str == end) {
		rslt2file("stoi error: Cannot to convert string \"%s\" to an integer.\n", str);
		logfile(LOG_ERROR, "stoi error: Cannot to convert string \"%s\" to an integer.\n", str);
		return -1;
	}

	*result = num;

	return 0;
}
#endif

//---------------------------------------
// time duration
//---------------------------------------
uint64_t time_duration(const struct timespec *t1, const struct timespec *t2)
{
	// Get time difference between t2 and t1
	int64_t sec = t2->tv_sec - t1->tv_sec;
	int64_t nsec = t2->tv_nsec - t1->tv_nsec;

	if (nsec < 0) {
		sec--;
		nsec += 1000000000L;
	}
	nsec = nsec + sec*1000000000L;

	return nsec;
}

//---------------------------------------
// Rounds a number up to a power of 2
//---------------------------------------
uint32_t nextPow2(int32_t n)
{
	// Returns 0 if n is less than or equal to 0
	if (n <= 0)
		return 0;

	// Return as is if n is a power of 2
	if ((n & (n - 1)) == 0)
		return (uint32_t)n;

	// Find the power of 2 next to n
	uint32_t ret = 1;
	while (n > 0) {
		ret <<= 1;
		n >>= 1;
	}

	return ret;
}
