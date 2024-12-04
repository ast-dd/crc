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

func bytesToChecksum(bytes []byte) uint64 {
	full := make([]byte, 8)
	copy(full, bytes)
	return binary.LittleEndian.Uint64(full)
}

// AppendCRCBytes returns a copy of the data byte slice with the checksum appended
func AppendCRCBytes(crcParams *Parameters, data []byte) []byte {
	checksum := CalculateCRCBytes(crcParams, data)
	appended := make([]byte, len(data), len(data)+len(checksum))
	copy(appended, data)
	appended = append(appended, checksum...)
	return appended
}

func CheckCRCBytes(crcParams *Parameters, data []byte, checksum []byte) bool {
	if len(checksum) != int(crcParams.Width/8) {
		return false
	}
	got := bytesToChecksum(checksum)
	calculated := CalculateCRC(crcParams, data)
	return got == calculated
}
