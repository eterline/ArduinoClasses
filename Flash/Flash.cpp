#include <Arduino.h>
#include <SPI.h>

//Connection scheme MX25L4005A
//
// Bios MX MX25L4005A 25xx to ESP8266 connection
// CS#      (1) to 	D2 (GPIO4)
// SO       (2) to 	D6 (GPIO12)
// WP#      (3) to 	3.3 В (PULL_UP)
// GND      (4) to 	GND
// SI       (5) to 	D7 (GPIO13)
// SCLK     (6) to 	D5 (GPIO14)
// HOLD#    (7) to 	3.3 В (PULL_UP)
// VCC      (8) to 	3.3 В
//
// 8-pin DIP case pins
//
//               +------U------+
//          CS --| 1         8 |-- VCC
//          SO --| 2         7 |-- HOLD# PULL_UP
// PULL_UP WP# --| 3         6 |-- SCK
//         GND --| 4         5 |-- SI
//               +-------------+

class Flash
{
private:
    byte __manuf, __mem_type, __cap;
    uint32_t fails  = 0;
    size_t wrote    = 0;
    size_t read     = 0;
    uint8_t __cs_pin;

    void gotFail() { fails++; }
    void addWriteBytes(size_t v) { wrote += v; }
    void addReadBytes(size_t v) { read += v; }

    unsigned long getFlashSizeBytes() {
        return 1UL << __cap;
    }

    bool waitForWriteEnd() {
        byte status;
        do {
            digitalWrite(__cs_pin, LOW);
            SPI.transfer(0x05);
            status = SPI.transfer(0x00);
            digitalWrite(__cs_pin, HIGH);
        } while (status & 0x01);

        return true;
    }

public:
    Flash(uint8_t cs_pin) {
        __cs_pin = cs_pin;
        pinMode(__cs_pin, OUTPUT);
        digitalWrite(__cs_pin, HIGH);
    };

    uint32_t Fails() { return fails; };

    bool run() {
        SPI.begin();

        digitalWrite(__cs_pin, LOW);
        SPI.transfer(0x9F);
        byte __manuf = SPI.transfer(0);
        byte __mem_type = SPI.transfer(0);
        byte __cap = SPI.transfer(0);
        digitalWrite(__cs_pin, HIGH);

        return __manuf != 0 && __mem_type != 0 && __cap != 0;
    }

    size_t ReadFlashTo(uint32_t startAt, byte* buff, size_t len) {

        if (startAt >= getFlashSizeBytes()) {
            gotFail();
            return 0;
        }

        if ((startAt + len) > getFlashSizeBytes()) {
            len = getFlashSizeBytes() - startAt;
        }

        digitalWrite(__cs_pin, LOW);

        SPI.transfer(0x03);
        SPI.transfer((startAt >> 16) & 0xFF);
        SPI.transfer((startAt >> 8) & 0xFF);
        SPI.transfer(startAt & 0xFF);

        for (size_t i = 0; i < len; i++) buff[i] = SPI.transfer(0x00); 

        digitalWrite(__cs_pin, HIGH);

        return len;
    }

    void EraseSector(uint32_t address) {
        digitalWrite(__cs_pin, LOW);
        SPI.transfer(0x06);
        digitalWrite(__cs_pin, HIGH);
        delayMicroseconds(1);

        digitalWrite(__cs_pin, LOW);
        SPI.transfer(0x20);
        SPI.transfer((address >> 16) & 0xFF);
        SPI.transfer((address >> 8) & 0xFF);
        SPI.transfer(address & 0xFF);
        digitalWrite(__cs_pin, HIGH);

        delay(50);
    }

    bool WritePage(uint32_t address, size_t length) {
        if (length == 0 || length > 256) return false;
        if ((address & 0xFF) + length > 256) return false;

        digitalWrite(__cs_pin, LOW);
        SPI.transfer(0x06);
        digitalWrite(__cs_pin, HIGH);
        delayMicroseconds(1);

        digitalWrite(__cs_pin, LOW);
        SPI.transfer(0x02);  // Page Program
        SPI.transfer((address >> 16) & 0xFF);
        SPI.transfer((address >> 8) & 0xFF);
        SPI.transfer(address & 0xFF);

        for (size_t i = 0; i < length; i++) {
            SPI.transfer(0x10);
        }

        digitalWrite(__cs_pin, HIGH);

        return waitForWriteEnd();
    }

    void WriteAt(uint32_t startAt, byte value) {

        digitalWrite(__cs_pin, LOW);
        SPI.transfer(0x06);
        digitalWrite(__cs_pin, HIGH);
        delayMicroseconds(1);

        digitalWrite(__cs_pin, LOW);

        SPI.transfer(0x02);
        SPI.transfer((startAt >> 16) & 0xFF);
        SPI.transfer((startAt >> 8) & 0xFF);
        SPI.transfer(startAt & 0xFF);
        SPI.transfer(value);

        digitalWrite(__cs_pin, HIGH);
        delay(5);
    }

    ~Flash() {};
};
