package crc

import (
	"encoding/binary"
)

// CalculateCRCBytes works according to CalculateCRC, but returns a byte slice
func CalculateCRCBytes(crcParams *Parameters, data []byte) (bs []byte) {
	checksum := CalculateCRC(crcParams, data)
	bs = make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, checksum)
	bs = bs[:crcParams.Width/8]

	return
}
