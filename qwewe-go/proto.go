package qweweproto

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type QwewePacket struct {
	Length   uint16
	SourceID byte
	TargetID byte
	Data     []byte
	CRC8     byte
}

// =============================================================================

func NewPacket() QwewePacket {
	return QwewePacket{
		Length:   3,
		SourceID: 0xFF,
		TargetID: 0xFF,
		Data:     make([]byte, 0),
		CRC8:     0x00,
	}
}

func QweeOfBytes(buf []byte, testCRC bool) (packet QwewePacket, ok bool) {
	if buf == nil || len(buf) < 5 {
		return packet, false
	}

	packet.CRC8 = buf[len(buf)-1]

	if testCRC {
		exceptedCRC8 := ChecksumCRC(buf[2:(len(buf)-1)], CRC8)

		if exceptedCRC8 != packet.CRC8 {
			return packet, false
		}
	}

	packet.Length = uint16(len(buf) - 1)
	packet.SourceID = buf[2]
	packet.TargetID = buf[3]
	packet.Data = buf[4:(len(buf) - 1)]

	return packet, true
}

func ReadQwewePacket(r io.Reader, testCRC bool) (QwewePacket, error) {
	reader := bufio.NewReader(r)

	for {
		packetBytes, err := reader.ReadBytes('\n')
		if err != nil {
			return QwewePacket{}, fmt.Errorf("read error: %w", err)
		}

		if len(packetBytes) < 2 {
			continue
		}

		packetBytes = packetBytes[:len(packetBytes)-1]

		packet, ok := QweeOfBytes(packetBytes, testCRC)
		if !ok {
			continue
		}

		return packet, nil
	}
}

// =============================================================================

func (pkt *QwewePacket) selfCheckSum() {
	data := []byte{pkt.SourceID, pkt.TargetID}
	data = append(data, pkt.Data...)
	pkt.CRC8 = ChecksumCRC(data, CRC8)
}

// =============================================================================

func (pkt *QwewePacket) Write(p []byte) (n int, err error) {
	dataLen := len(p)
	pkt.Data = append(pkt.Data, p...)
	pkt.Length += uint16(dataLen)
	return dataLen, nil
}

func (pkt *QwewePacket) Read(p []byte) (n int, err error) {
	if pkt == nil || pkt.Data == nil {
		return 0, ErrInvalidPacket
	}

	return bytes.NewReader(pkt.Data).Read(p)
}

func (pkt *QwewePacket) Parse(value interface{}) error {
	var buf bytes.Buffer

	err := binary.Write(&buf, binary.LittleEndian, value)
	if err != nil {
		return err
	}
	_, err = pkt.Write(buf.Bytes())
	return err
}

func (pkt *QwewePacket) Marshal() ([]byte, error) {
	buf := make([]byte, pkt.Length+1+2)
	pkt.selfCheckSum()

	// write header bytes
	buf[0] = byte(PKT_QW_START)
	buf[1] = byte(pkt.Length >> 8)
	buf[2] = byte(pkt.Length)
	buf[3] = pkt.SourceID
	buf[4] = pkt.TargetID

	// write data
	copy(buf[5:], pkt.Data)

	// write checksum
	buf[pkt.Length+2] = pkt.CRC8

	return buf, nil
}

func (pkt *QwewePacket) Unmarshal(value interface{}) error {

	if pkt == nil || pkt.Data == nil {
		return ErrInvalidPayload
	}

	return binary.Read(pkt, binary.LittleEndian, value)
}

// =============================================================================

func (pkt QwewePacket) WideSource() bool {
	return pkt.SourceID == 0xFF
}

func (pkt *QwewePacket) MustWideSource() {
	pkt.SourceID = 0xFF
}

func (pkt QwewePacket) WideTarget() bool {
	return pkt.TargetID == 0xFF
}

func (pkt *QwewePacket) MustWideTarget() {
	pkt.TargetID = 0xFF
}

// =============================================================================
