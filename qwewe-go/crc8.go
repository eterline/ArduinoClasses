package qweweproto

import "math/bits"

// Params represents parameters of a CRC-8 algorithm including polynomial and initial value.
// More information about algorithms parametrization and parameter descriptions
// can be found here - http://www.zlib.net/crc_v3.txt
type Params struct {
	Poly   uint8
	Init   uint8
	RefIn  bool
	RefOut bool
	XorOut uint8
	Check  uint8
	Name   string
}

// Predefined CRC-8 algorithms.
// List of algorithms with their parameters borrowed from here - http://reveng.sourceforge.net/crc-catalogue/1-15.htm#crc.cat-bits.8
//
// The variables can be used to create Table for the selected algorithm.
var (
	CRC8          = Params{0x07, 0x00, false, false, 0x00, 0xF4, "CRC-8"}
	CRC8_CDMA2000 = Params{0x9B, 0xFF, false, false, 0x00, 0xDA, "CRC-8/CDMA2000"}
	CRC8_DARC     = Params{0x39, 0x00, true, true, 0x00, 0x15, "CRC-8/DARC"}
	CRC8_DVB_S2   = Params{0xD5, 0x00, false, false, 0x00, 0xBC, "CRC-8/DVB-S2"}
	CRC8_EBU      = Params{0x1D, 0xFF, true, true, 0x00, 0x97, "CRC-8/EBU"}
	CRC8_I_CODE   = Params{0x1D, 0xFD, false, false, 0x00, 0x7E, "CRC-8/I-CODE"}
	CRC8_ITU      = Params{0x07, 0x00, false, false, 0x55, 0xA1, "CRC-8/ITU"}
	CRC8_MAXIM    = Params{0x31, 0x00, true, true, 0x00, 0xA1, "CRC-8/MAXIM"}
	CRC8_ROHC     = Params{0x07, 0xFF, true, true, 0x00, 0xD0, "CRC-8/ROHC"}
	CRC8_WCDMA    = Params{0x9B, 0x00, true, true, 0x00, 0x25, "CRC-8/WCDMA"}
)

// Table is a 256-byte table representing polynomial and algorithm settings for efficient processing.
type Table struct {
	params Params
	data   [256]uint8
}

// MakeTable returns the Table constructed from the specified algorithm.
func newCRCTable(params Params) *Table {
	table := new(Table)
	table.params = params
	for n := 0; n < 256; n++ {
		crc := uint8(n)
		for i := 0; i < 8; i++ {
			bit := (crc & 0x80) != 0
			crc <<= 1
			if bit {
				crc ^= params.Poly
			}
		}
		table.data[n] = crc
	}
	return table
}

// Init returns the initial value for CRC register corresponding to the specified algorithm.
func InitCRCByte(params Params) *DataCRC {
	table := newCRCTable(params)
	return &DataCRC{
		Register: table.params.Init,
		Table:    table,
	}
}

type DataCRC struct {
	Register uint8
	Table    *Table
}

// Update returns the result of adding the bytes in data to the crc.
func (crc *DataCRC) Update(data []byte) {
	if crc.Table.params.RefIn {
		for _, d := range data {
			d = bits.Reverse8(d)
			crc.Register = crc.Table.data[crc.Register^d]
		}
	} else {
		for _, d := range data {
			crc.Register = crc.Table.data[crc.Register^d]
		}
	}
}

// Complete returns the result of CRC calculation and post-calculation processing of the crc.
func (crc *DataCRC) Complete() uint8 {
	if crc.Table.params.RefOut {
		crc.Register = bits.Reverse8(crc.Register)
	}

	return crc.Register ^ crc.Table.params.XorOut
}

// Checksum returns CRC checksum of data using specified algorithm represented by the Table.
func ChecksumCRC(data []byte, params Params) uint8 {
	crc := InitCRCByte(params)
	crc.Update(data)
	return crc.Complete()
}
