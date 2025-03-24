/*************************************************
* Copyright 2025 NTT Corporation , FUJITSU LIMITED
*************************************************/
#ifndef FCC_JSON_H_
#define FCC_JSON_H_

#ifdef __cplusplus
extern "C" {
#endif


int fcc_json_parse_string(
        const char *json_string);

int fcc_json_parse_file(
        const char *json_file_path);

int fcc_json_create_output_file(
        void);


#ifdef __cplusplus
}
#endif

#endif  // FCC_JSON_H_
