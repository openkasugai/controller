/*************************************************
* Copyright 2025 NTT Corporation , FUJITSU LIMITED
*************************************************/

#include "logger.hpp"

#include <stdio.h>
#include <stdlib.h>
#include <stdarg.h>
#include <string.h>
#include <time.h>
#include <errno.h>


int logger::open(
  const char *filename)
{
  // Create log file
  if (!filename)
    filename = default_filename;
  if (!fp_) {
    fp_ = fopen(filename, "a");
    if (fp_ == NULL) {
      int err = errno;
      fprintf(stdout, " ! Failed to open %s: %s\n", filename, strerror(err));
      fflush(stdout);
      return -1;
    }
    fprintf(fp_, " * Create logfile : %s\n", filename);
    fflush(fp_);
  }
  return 0;
}


int logger::close(
  void)
{
  // Create log file
  if (fp_ && (fp_ != stdout) && (fp_ != stderr)) {
    fclose(fp_);
    fp_ = NULL;
    return 0;
  }
  return -1;
}


int logger::print(
  const char *format,
  ...)
{
  time_t t;
  char date[32];

  va_list args;

  if (is_use_timestamp_) {
    // print timestamp
    time(&t);
    strftime(date, sizeof(date), "%H:%M:%S", localtime(&t));
    fprintf(fp_, "%s", date);
  }

  // Print log into logfile
  va_start(args, format);
  vfprintf(fp_, format, args);
  va_end(args);
  fflush(fp_);

  return 0;
}


void logger::set_timestamp(
  bool is_use)
{
  is_use_timestamp_ = is_use;
}


bool logger::get_timestamp(
  void)
{
  return is_use_timestamp_;
}
