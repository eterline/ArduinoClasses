package qweweproto

type QweweProtoError string

func (e QweweProtoError) Error() string {
	return "qwewe protocol error: " + string(e)
}

const (
	ErrInvalidPacket  QweweProtoError = "invalid data structure"
	ErrInvalidLength  QweweProtoError = "invalid packet length information"
	ErrInvalidPayload QweweProtoError = "invalid packet data payload"
	ErrInvalidCRC     QweweProtoError = "invalid control sum CRC8"
)
