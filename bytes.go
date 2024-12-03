package crc

import (
	"encoding/binary"
)

// CalculateCRCBytes works according to CalculateCRC, but returns a byte slice
func CalculateCRCBytes(crcParams *Parameters, data []byte) []byte {
	checksum := CalculateCRC(crcParams, data)
	return checksumToBytes(checksum, int(crcParams.Width/8))
}

// CalculateCRCBytes works according to CalculateCRC, but returns a byte slice
func (h *Hash) CalculateCRCBytes(data []byte) []byte {
	checksum := h.CalculateCRC(data)
	return checksumToBytes(checksum, int(h.table.crcParams.Width/8))
}

func checksumToBytes(checksum uint64, l int) []byte {
	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, checksum)
	return bs[:l]
}
