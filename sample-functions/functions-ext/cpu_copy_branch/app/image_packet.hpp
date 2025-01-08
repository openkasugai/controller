//=========================================================================
// image_packet.hpp Copyright 2024 NTT Corporation , FUJITSU LIMITED
//-------------------------------------------------------------------------
#ifndef IMAGE_PACKET_HPP
#define IMAGE_PACKET_HPP

#include <stdint.h>

struct ImagePacketHeader {
    uint32_t marker;
    uint32_t payload_len;
    uint8_t reserve1[4];
    uint32_t sequence_num;
    uint8_t reserve2[8];
    uint64_t timestamp;
    uint32_t data_id;
    uint8_t reserve3[10];
    uint16_t header_checksum;
};

#endif
