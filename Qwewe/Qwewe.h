#pragma once
#include <Arduino.h>
#include <stdint.h>
#include <string.h>

#define PKT_START_BYTE      '\n'
#define MAX_PACKET_SIZE     1024
#define PACKET_BASE_SIZE    3

class QwewePacket {
private:
    uint16_t packet_ln  = 0;
    uint8_t  target_id  = 0;
    uint8_t  source_id  = 0;
    uint8_t  crc8_hash  = 0;

public:
    QwewePacket() = default;

    void create(uint8_t src_id, uint8_t tgt_id) {
        source_id = src_id;
        target_id = tgt_id;
    }

    void create() {
        packet_ln = PACKET_BASE_SIZE;
        source_id = 0xFF;
        target_id = 0xFF;
    }

    // CRC8 calculation (stub — replace with real implementation)
    uint8_t calc_crc8(const uint8_t* data, size_t len) const {
        uint8_t crc = 0;
        for (size_t i = 0; i < len; ++i) {
            crc ^= data[i];
        }
        return crc;
    }

    // send to buffer from raw bytes
    void put(uint8_t* out_buf, const uint8_t* data, size_t data_len) {
        packet_ln = PACKET_BASE_SIZE + data_len + 1;
        crc8_hash = calc_crc8(data, data_len);

        out_buf[0] = (uint8_t)PKT_START_BYTE;
        out_buf[1] = highByte(packet_ln);
        out_buf[2] = lowByte(packet_ln);
        out_buf[3] = target_id;
        out_buf[4] = source_id;
        memcpy(&out_buf[5], data, data_len);
        out_buf[5 + data_len] = crc8_hash;
    }

    // send to stream from raw bytes
    void put(Stream& stream, const uint8_t* data, size_t data_len) {
        packet_ln = PACKET_BASE_SIZE + data_len + 1;
        crc8_hash = calc_crc8(data, data_len);

        stream.write((uint8_t)PKT_START_BYTE);
        stream.write(highByte(packet_ln));
        stream.write(lowByte(packet_ln));
        stream.write(target_id);
        stream.write(source_id);
        stream.write(data, data_len);
        stream.write(12);
    }

    // Template: send to buffer from structure
    template<typename T>
    void put(uint8_t* out_buf, const T& structure) {
        const uint8_t* data = reinterpret_cast<const uint8_t*>(&structure);
        put(out_buf, data, sizeof(T));
    }

    // Template: send to stream from structure
    template<typename T>
    void put(Stream& stream, const T& structure) {
        const uint8_t* data = reinterpret_cast<const uint8_t*>(&structure);
        put(stream, data, sizeof(T));
    }
};