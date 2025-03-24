/*************************************************
* Copyright 2025 NTT Corporation , FUJITSU LIMITED
*************************************************/

#include <fcc_log.h>


uint32_t fcc_log_get_app_version(
  void)
{
  uint32_t version = 0;
  version |= APP_MAJOR_VER;
  version <<= 8;
  version |= APP_MINOR_VER;
  version <<= 8;
  version |= APP_REVISION;
  version <<= 8;
  version |= APP_PATCH;
  return version;
}

uint32_t fcc_log_get_major_version(
  void)
{
  return APP_MAJOR_VER;
}

uint32_t fcc_log_get_minor_version(
  void)
{
  return APP_MINOR_VER;
}

uint32_t fcc_log_get_revision(
  void)
{
  return APP_REVISION;
}

uint32_t fcc_log_get_patch(
  void)
{
  return APP_PATCH;
}
