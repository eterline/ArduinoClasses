package qweweproto

type QweweMeanByte byte

const (
	QW_PKT_HEADER_SIZE int = 5

	PKT_QW_START   QweweMeanByte = '\n'
	PKT_QW_WIDE_ID QweweMeanByte = 0xFF

	TYPE_QW_ACK     QweweMeanByte = 0x01
	TYPE_QW_DATA    QweweMeanByte = 0x02
	TYPE_QW_COMMAND QweweMeanByte = 0x03
)
