/*************************************************
* Copyright 2025 NTT Corporation , FUJITSU LIMITED
*************************************************/

#ifndef LOGGER_HPP__
#define LOGGER_HPP__


#include <stdio.h>


class logger {
 public:
  logger(const char *filename = NULL) {
    if (open(filename))
      fp_ = stdout;
  };
  ~logger() {
    close();
  };

  int open(
          const char *filename = NULL);
  int close(
          void);
  int print(
          const char *fmt,
          ...);
  void set_timestamp(
          bool is_use);
  bool get_timestamp(
          void);

 private:
  FILE *fp_ = NULL;
  const char *default_filename = "logger.log";
  bool is_use_timestamp_ = true;
};

#endif  // LOGGER_HPP__
