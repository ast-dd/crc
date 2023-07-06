// Copyright 2016, S&K Software Development Ltd.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package crc implements generic CRC calculations up to 64 bits wide.
// It aims to be fairly complete, allowing users to match pretty much
// any CRC algorithm used in the wild by choosing appropriate Parameters.
// And it's also fairly fast for everyday use.
//
// This package has been largely inspired by Ross Williams' 1993 paper "A Painless Guide to CRC Error Detection Algorithms".
// A good list of parameter sets for various CRC algorithms can be found at http://reveng.sourceforge.net/crc-catalogue/.
package crc

// Parameters represents set of parameters defining a particular CRC algorithm.
type Parameters struct {
	Width      uint   // Width of the CRC expressed in bits
	Polynomial uint64 // Polynomial used in this CRC calculation
	ReflectIn  bool   // ReflectIn indicates whether input bytes should be reflected
	ReflectOut bool   // ReflectOut indicates whether input bytes should be reflected
	Init       uint64 // Init is initial value for CRC calculation
	FinalXor   uint64 // Xor is a value for final xor to be applied before returning result
}

var (
	// CRC-8, CRC-8/SMBUS
	CRC8 = &Parameters{Width: 8, Polynomial: 0x07, Init: 0x00, ReflectIn: false, ReflectOut: false, FinalXor: 0x00}
	// CRC-8/CDMA2000
	CRC8CDMA2000 = &Parameters{Width: 8, Polynomial: 0x9B, Init: 0xFF, ReflectIn: false, ReflectOut: false, FinalXor: 0x00}
	// CRC-8/DARC
	CRC8DARC = &Parameters{Width: 8, Polynomial: 0x39, Init: 0x00, ReflectIn: true, ReflectOut: true, FinalXor: 0x00}
	// CRC-8/DVB-S2
	CRC8DVBS2 = &Parameters{Width: 8, Polynomial: 0xD5, Init: 0x00, ReflectIn: false, ReflectOut: false, FinalXor: 0x00}
	// CRC-8/TECH-3250, CRC-8/AES, CRC-8/EBU
	CRC8EBU = &Parameters{Width: 8, Polynomial: 0x1D, Init: 0xFF, ReflectIn: true, ReflectOut: true, FinalXor: 0x00}
	// CRC-8/ICODE
	CRC8ICODE = &Parameters{Width: 8, Polynomial: 0x1D, Init: 0xFD, ReflectIn: false, ReflectOut: false, FinalXor: 0x00}
	// CRC-8/ITU, CRC-8/I-432-1
	CRC8ITU = &Parameters{Width: 8, Polynomial: 0x07, Init: 0x00, ReflectIn: false, ReflectOut: false, FinalXor: 0x55}
	// CRC-8/MAXIM, CRC-8/MAXIM-DOW, DOW-CRC
	CRC8MAXIM = &Parameters{Width: 8, Polynomial: 0x31, Init: 0x00, ReflectIn: true, ReflectOut: true, FinalXor: 0x00}
	// CRC-8/ROHC
	CRC8ROHC = &Parameters{Width: 8, Polynomial: 0x07, Init: 0xFF, ReflectIn: true, ReflectOut: true, FinalXor: 0x00}
	// CRC-8/WCDMA
	CRC8WCDMA = &Parameters{Width: 8, Polynomial: 0x9B, Init: 0x00, ReflectIn: true, ReflectOut: true, FinalXor: 0x00}

	// Missing
	// CRC-8/AUTOSAR
	CRC8AUTOSAR = &Parameters{Width: 8, Polynomial: 0x2F, Init: 0xFF, ReflectIn: false, ReflectOut: false, FinalXor: 0xFF}
	// CRC-8/BLUETOOTH
	CRC8BLUETOOTH = &Parameters{Width: 8, Polynomial: 0xA7, Init: 0x00, ReflectIn: true, ReflectOut: true, FinalXor: 0x00}
	// CRC-8/GSM-A
	CRC8GSMA = &Parameters{Width: 8, Polynomial: 0x1D, Init: 0x00, ReflectIn: false, ReflectOut: false, FinalXor: 0x00}
	// CRC-8/GSM-B
	CRC8GSMB = &Parameters{Width: 8, Polynomial: 0x49, Init: 0x00, ReflectIn: false, ReflectOut: false, FinalXor: 0xFF}
	// CRC-8/HITAG
	CRC8HITAG = &Parameters{Width: 8, Polynomial: 0x1D, Init: 0xFF, ReflectIn: false, ReflectOut: false, FinalXor: 0x00}
	// CRC-8/LTE
	CRC8LTE = &Parameters{Width: 8, Polynomial: 0x9b, Init: 0x00, ReflectIn: false, ReflectOut: false, FinalXor: 0x00}
	// CRC-8/MIFARE-MAD
	CRC8MIFAREMAD = &Parameters{Width: 8, Polynomial: 0x1D, Init: 0xC7, ReflectIn: false, ReflectOut: false, FinalXor: 0x00}
	// CRC-8/NRSC-5
	CRC8NRSC5 = &Parameters{Width: 8, Polynomial: 0x31, Init: 0xFF, ReflectIn: false, ReflectOut: false, FinalXor: 0x00}
	// CRC-8/OPENSAFETY
	CRC8OPENSAFETY = &Parameters{Width: 8, Polynomial: 0x2F, Init: 0x00, ReflectIn: false, ReflectOut: false, FinalXor: 0x00}
	// CRC-8/SAE-J1850
	CRC8SAEJ1850 = &Parameters{Width: 8, Polynomial: 0x1D, Init: 0xFF, ReflectIn: false, ReflectOut: false, FinalXor: 0xFF}

	// CRC-16/ARC, ARC, CRC-16, CRC-16/LHA, CRC-IBM
	CRC16ARC = &Parameters{Width: 16, Polynomial: 0x8005, Init: 0x0000, ReflectIn: true, ReflectOut: true, FinalXor: 0x0000}
	// CRC-16/SPI-FUJITSU, CRC-16/AUG-CCITT
	CRC16AUGCCITT = &Parameters{Width: 16, Polynomial: 0x1021, Init: 0x1D0F, ReflectIn: false, ReflectOut: false, FinalXor: 0x0000}
	// CRC-16/UMTS, CRC-16/BUYPASS, CRC-16/VERIFONE
	CRC16BUYPASS = &Parameters{Width: 16, Polynomial: 0x8005, Init: 0x0000, ReflectIn: false, ReflectOut: false, FinalXor: 0x0000}
	// CCITT CRC parameters, CRC-16/IBM-3740, CRC-16/AUTOSAR
	CRC16CCITTFALSE = &Parameters{Width: 16, Polynomial: 0x1021, Init: 0xFFFF, ReflectIn: false, ReflectOut: false, FinalXor: 0x0000}
	CCITT           = CRC16CCITTFALSE
	// CRC-16/CDMA2000
	CRC16CDMA2000 = &Parameters{Width: 16, Polynomial: 0xC867, Init: 0xFFFF, ReflectIn: false, ReflectOut: false, FinalXor: 0x0000}
	// CRC-16/DDS-110
	CRC16DDS110 = &Parameters{Width: 16, Polynomial: 0x8005, Init: 0x800D, ReflectIn: false, ReflectOut: false, FinalXor: 0x0000}
	// CRC-16/DECT-R, R-CRC-16
	CRC16DECTR = &Parameters{Width: 16, Polynomial: 0x0589, Init: 0x0000, ReflectIn: false, ReflectOut: false, FinalXor: 0x0001}
	// CRC-16/DECT-X, X-CRC-16
	CRC16DECTX = &Parameters{Width: 16, Polynomial: 0x0589, Init: 0x0000, ReflectIn: false, ReflectOut: false, FinalXor: 0x0000}
	// CRC-16/DNP
	CRC16DNP = &Parameters{Width: 16, Polynomial: 0x3D65, Init: 0x0000, ReflectIn: true, ReflectOut: true, FinalXor: 0xFFFF}
	// CRC-16/EN-13757
	CRC16EN13757 = &Parameters{Width: 16, Polynomial: 0x3D65, Init: 0x0000, ReflectIn: false, ReflectOut: false, FinalXor: 0xFFFF}
	// CRC-16/GENIBUS, CRC-16/DARC, CRC-16/EPC, CRC-16/EPC-C1G2, CRC-16/I-CODE
	CRC16GENIBUS = &Parameters{Width: 16, Polynomial: 0x1021, Init: 0xFFFF, ReflectIn: false, ReflectOut: false, FinalXor: 0xFFFF}
	// CRC-16/KERMIT, CRC-16/BLUETOOTH, CRC-16/CCITT, CRC-16/CCITT-TRUE, CRC-16/V-41-LSB, CRC-CCITT, KERMIT
	CRC16KERMIT = &Parameters{Width: 16, Polynomial: 0x1021, Init: 0x0000, ReflectIn: true, ReflectOut: true, FinalXor: 0x0000}
	// CRC-16/MAXIM-DOW, CRC-16/MAXIM
	CRC16MAXIM = &Parameters{Width: 16, Polynomial: 0x8005, Init: 0x0000, ReflectIn: true, ReflectOut: true, FinalXor: 0xFFFF}
	// CRC-16/MCRF4XX
	CRC16MCRF4XX = &Parameters{Width: 16, Polynomial: 0x1021, Init: 0xFFFF, ReflectIn: true, ReflectOut: true, FinalXor: 0x0000}
	// CRC-16/MODBUS, MODBUS
	CRC16MODBUS = &Parameters{Width: 16, Polynomial: 0x8005, Init: 0xFFFF, ReflectIn: true, ReflectOut: true, FinalXor: 0x0000}
	// CRC-16/RIELLO
	CRC16RIELLO = &Parameters{Width: 16, Polynomial: 0x1021, Init: 0xB2AA, ReflectIn: true, ReflectOut: true, FinalXor: 0x0000}
	// CRC-16/T10-DIF
	CRC16T10DIF = &Parameters{Width: 16, Polynomial: 0x8BB7, Init: 0x0000, ReflectIn: false, ReflectOut: false, FinalXor: 0x0000}
	// CRC-16/TELEDISK
	CRC16TELEDISK = &Parameters{Width: 16, Polynomial: 0xA097, Init: 0x0000, ReflectIn: false, ReflectOut: false, FinalXor: 0x0000}
	// CRC-16/TMS37157
	CRC16TMS37157 = &Parameters{Width: 16, Polynomial: 0x1021, Init: 0x89EC, ReflectIn: true, ReflectOut: true, FinalXor: 0x0000}
	// CRC-16/USB
	CRC16USB = &Parameters{Width: 16, Polynomial: 0x8005, Init: 0xFFFF, ReflectIn: true, ReflectOut: true, FinalXor: 0xFFFF}
	// CRC-16/IBM-SDLC, CRC-16/ISO-HDLC, CRC-16/ISO-IEC-14443-3-B, CRC-16/X-25, CRC-B, X-25
	CRC16X25 = &Parameters{Width: 16, Polynomial: 0x1021, Init: 0xFFFF, ReflectIn: true, ReflectOut: true, FinalXor: 0xFFFF}
	X25      = CRC16X25
	// CRC-16/XMODEM, CRC-16/ACORN, CRC-16/LTE, CRC-16/V-41-MSB, XMODEM, ZMODEM
	CRC16XMODEM = &Parameters{Width: 16, Polynomial: 0x1021, Init: 0x0000, ReflectIn: false, ReflectOut: false, FinalXor: 0x0000}
	// XMODEM2 is another set of CRC parameters commonly referred as "XMODEM"
	XMODEM2 = &Parameters{Width: 16, Polynomial: 0x8408, Init: 0x0000, ReflectIn: true, ReflectOut: true, FinalXor: 0x0}
	// CRC-16/ISO-IEC-14443-3-A, CRC-A
	CRCA = &Parameters{Width: 16, Polynomial: 0x1021, Init: 0xC6C6, ReflectIn: true, ReflectOut: true, FinalXor: 0x0000}

	// MISSING
	// CRC-16/CMS
	CRC16CMS = &Parameters{Width: 16, Polynomial: 0x8005, Init: 0xFFFF, ReflectIn: false, ReflectOut: false, FinalXor: 0x0000}
	// CRC-16/GSM
	CRC16GSM = &Parameters{Width: 16, Polynomial: 0x1021, Init: 0x0000, ReflectIn: false, ReflectOut: false, FinalXor: 0xFFFF}
	// CRC-16/LJ1200
	CRC16LJ1200 = &Parameters{Width: 16, Polynomial: 0x6F63, Init: 0x0000, ReflectIn: false, ReflectOut: false, FinalXor: 0x0000}
	// CRC-16/M17
	CRC16M17 = &Parameters{Width: 16, Polynomial: 0x5935, Init: 0xFFFF, ReflectIn: false, ReflectOut: false, FinalXor: 0x0000}
	// CRC-16/NRSC-5
	CRC16NRSC5 = &Parameters{Width: 16, Polynomial: 0x080B, Init: 0xFFFF, ReflectIn: true, ReflectOut: true, FinalXor: 0x0000}
	// CRC-16/OPENSAFETY-A
	CRC16OPENSAFETYA = &Parameters{Width: 16, Polynomial: 0x5935, Init: 0x0000, ReflectIn: false, ReflectOut: false, FinalXor: 0x0000}
	// CRC-16/OPENSAFETY-B
	CRC16OPENSAFETYB = &Parameters{Width: 16, Polynomial: 0x755B, Init: 0x0000, ReflectIn: false, ReflectOut: false, FinalXor: 0x0000}
	// CRC-16/PROFIBUS
	CRC16PROFIBUS = &Parameters{Width: 16, Polynomial: 0x1DCF, Init: 0xFFFF, ReflectIn: false, ReflectOut: false, FinalXor: 0xFFFF}

	// CRC32 is by far the the most commonly used CRC-32 polynom and set of parameters
	// CRC-32, CRC-32/ISO-HDLC, CRC-32/ADCCP, CRC-32/V-42, CRC-32/XZ, PKZIP
	CRC32 = &Parameters{Width: 32, Polynomial: 0x04C11DB7, Init: 0xFFFFFFFF, ReflectIn: true, ReflectOut: true, FinalXor: 0xFFFFFFFF}
	// IEEE is an alias to CRC32
	IEEE = CRC32
	// CRC-32/BZIP2, CRC-32/AAL5, CRC-32/DECT-B, B-CRC-32
	CRC32BZIP2 = &Parameters{Width: 32, Polynomial: 0x04C11DB7, Init: 0xFFFFFFFF, ReflectIn: false, ReflectOut: false, FinalXor: 0xFFFFFFFF}
	// CRC-32/JAMCRC, JAMCRC
	CRC32JAMCRC = &Parameters{Width: 32, Polynomial: 0x04C11DB7, Init: 0xFFFFFFFF, ReflectIn: true, ReflectOut: true, FinalXor: 0x00000000}
	// CRC-32/MPEG-2
	CRC32MPEG2 = &Parameters{Width: 32, Polynomial: 0x04C11DB7, Init: 0xFFFFFFFF, ReflectIn: false, ReflectOut: false, FinalXor: 0x00000000}
	// CRC-32/POSIX, CRC-32/CKSUM, CKSUM
	CRC32POSIX = &Parameters{Width: 32, Polynomial: 0x04C11DB7, Init: 0x00000000, ReflectIn: false, ReflectOut: false, FinalXor: 0xFFFFFFFF}
	// CRC-32/SATA
	CRC32SATA = &Parameters{Width: 32, Polynomial: 0x04C11DB7, Init: 0x52325032, ReflectIn: false, ReflectOut: false, FinalXor: 0x00000000}
	// CRC-32/XFER
	CRC32XFER = &Parameters{Width: 32, Polynomial: 0x000000AF, Init: 0x00000000, ReflectIn: false, ReflectOut: false, FinalXor: 0x00000000}
	// CRC-32C, CRC-32/BASE91-C, CRC-32/CASTAGNOLI, CRC-32/INTERLAKEN, CRC-32/ISCSI
	CRC32C     = &Parameters{Width: 32, Polynomial: 0x1EDC6F41, Init: 0xFFFFFFFF, ReflectIn: true, ReflectOut: true, FinalXor: 0xFFFFFFFF}
	Castagnoli = CRC32C
	// CRC-32D, CRC-32/BASE91-D
	CRC32D = &Parameters{Width: 32, Polynomial: 0xA833982B, Init: 0xFFFFFFFF, ReflectIn: true, ReflectOut: true, FinalXor: 0xFFFFFFFF}
	// CRC-32Q, CRC-32/AIXM
	CRC32Q = &Parameters{Width: 32, Polynomial: 0x814141AB, Init: 0x00000000, ReflectIn: false, ReflectOut: false, FinalXor: 0x00000000}

	// MISSING
	// CRC-32/AUTOSAR
	CRC32AUTOSAR = &Parameters{Width: 32, Polynomial: 0xF4ACFB13, Init: 0xFFFFFFFF, ReflectIn: true, ReflectOut: true, FinalXor: 0xFFFFFFFF}
	// CRC-32/CD-ROM-EDC
	CRC32CDROMEDC = &Parameters{Width: 32, Polynomial: 0x8001801B, Init: 0x00000000, ReflectIn: true, ReflectOut: true, FinalXor: 0x00000000}
	// CRC-32/MEF
	CRC32MEF = &Parameters{Width: 32, Polynomial: 0x741B8CD7, Init: 0xFFFFFFFF, ReflectIn: true, ReflectOut: true, FinalXor: 0x00000000}

	// Koopman polynomial - is this CRC-32/MEF but finalxor is wrong?
	Koopman = &Parameters{Width: 32, Polynomial: 0x741B8CD7, Init: 0xFFFFFFFF, ReflectIn: true, ReflectOut: true, FinalXor: 0xFFFFFFFF}

	// CRC64ISO is set of parameters commonly known as CRC64-ISO
	CRC64ISO = &Parameters{Width: 64, Polynomial: 0x000000000000001B, Init: 0xFFFFFFFFFFFFFFFF, ReflectIn: true, ReflectOut: true, FinalXor: 0xFFFFFFFFFFFFFFFF}
	// CRC64ECMA is set of parameters commonly known as CRC64-ECMA
	CRC64ECMA = &Parameters{Width: 64, Polynomial: 0x42F0E1EBA9EA3693, Init: 0xFFFFFFFFFFFFFFFF, ReflectIn: true, ReflectOut: true, FinalXor: 0xFFFFFFFFFFFFFFFF}
)

// reflect reverses order of last count bits
func reflect(in uint64, count uint) uint64 {
	ret := in
	for idx := uint(0); idx < count; idx++ {
		srcbit := uint64(1) << idx
		dstbit := uint64(1) << (count - idx - 1)
		if (in & srcbit) != 0 {
			ret |= dstbit
		} else {
			ret = ret & (^dstbit)
		}
	}
	return ret
}

// CalculateCRC implements simple straight forward bit by bit calculation.
// It is relatively slow for large amounts of data, but does not require
// any preparation steps. As a result, it might be faster in some cases
// then building a table required for faster calculation.
//
// Note: this implementation follows section 8 ("A Straightforward CRC Implementation")
// of Ross N. Williams paper as even though final/sample implementation of this algorithm
// provided near the end of that paper (and followed by most other implementations)
// is a bit faster, it does not work for polynomials shorter then 8 bits. And if you need
// speed, you shoud probably be using table based implementation anyway.
func CalculateCRC(crcParams *Parameters, data []byte) uint64 {

	curValue := crcParams.Init
	topBit := uint64(1) << (crcParams.Width - 1)
	mask := (topBit << 1) - 1

	for i := 0; i < len(data); i++ {
		var curByte = uint64(data[i]) & 0x00FF
		if crcParams.ReflectIn {
			curByte = reflect(curByte, 8)
		}
		for j := uint64(0x0080); j != 0; j >>= 1 {
			bit := curValue & topBit
			curValue <<= 1
			if (curByte & j) != 0 {
				bit = bit ^ topBit
			}
			if bit != 0 {
				curValue = curValue ^ crcParams.Polynomial
			}
		}
	}
	if crcParams.ReflectOut {
		curValue = reflect(curValue, crcParams.Width)
	}

	curValue = curValue ^ crcParams.FinalXor

	return curValue & mask
}

// Table represents the partial evaluation of a checksum using table-driven
// implementation. It is essentially immutable once initialized and thread safe as a result.
type Table struct {
	crcParams Parameters
	crctable  []uint64
	mask      uint64
	initValue uint64
}

// NewTable creates and initializes a new Table for the CRC algorithm specified by the crcParams.
func NewTable(crcParams *Parameters) *Table {
	ret := &Table{crcParams: *crcParams}
	ret.mask = (uint64(1) << crcParams.Width) - 1
	ret.crctable = make([]uint64, 256, 256)
	ret.initValue = crcParams.Init
	if crcParams.ReflectIn {
		ret.initValue = reflect(crcParams.Init, crcParams.Width)
	}

	tmp := make([]byte, 1, 1)
	tableParams := *crcParams
	tableParams.Init = 0
	tableParams.ReflectOut = tableParams.ReflectIn
	tableParams.FinalXor = 0
	for i := 0; i < 256; i++ {
		tmp[0] = byte(i)
		ret.crctable[i] = CalculateCRC(&tableParams, tmp)
	}
	return ret
}

// InitCrc returns a stating value for a new CRC calculation
func (t *Table) InitCrc() uint64 {
	return t.initValue
}

// UpdateCrc process supplied bytes and updates current (partial) CRC accordingly.
// It can be called repetitively to process larger data in chunks.
func (t *Table) UpdateCrc(curValue uint64, p []byte) uint64 {
	if t.crcParams.ReflectIn {
		for _, v := range p {
			curValue = t.crctable[(byte(curValue)^v)&0xFF] ^ (curValue >> 8)
		}
	} else if t.crcParams.Width < 8 {
		for _, v := range p {
			curValue = t.crctable[((((byte)(curValue<<(8-t.crcParams.Width)))^v)&0xFF)] ^ (curValue << 8)
		}
	} else {
		for _, v := range p {
			curValue = t.crctable[((byte(curValue>>(t.crcParams.Width-8))^v)&0xFF)] ^ (curValue << 8)
		}
	}
	return curValue
}

// CRC returns CRC value for the data processed so far.
func (t *Table) CRC(curValue uint64) uint64 {
	ret := curValue

	if t.crcParams.ReflectOut != t.crcParams.ReflectIn {
		ret = reflect(ret, t.crcParams.Width)
	}
	return (ret ^ t.crcParams.FinalXor) & t.mask
}

// CRC8 is a convenience method to spare end users from explicit type conversion every time this package is used.
// Underneath, it just calls CRC() method.
func (t *Table) CRC8(curValue uint64) uint8 {
	return uint8(t.CRC(curValue))
}

// CRC16 is a convenience method to spare end users from explicit type conversion every time this package is used.
// Underneath, it just calls CRC() method.
func (t *Table) CRC16(curValue uint64) uint16 {
	return uint16(t.CRC(curValue))
}

// CRC32 is a convenience method to spare end users from explicit type conversion every time this package is used.
// Underneath, it just calls CRC() method.
func (t *Table) CRC32(curValue uint64) uint32 {
	return uint32(t.CRC(curValue))
}

// CalculateCRC is a convenience function allowing to calculate CRC in one call.
func (t *Table) CalculateCRC(data []byte) uint64 {
	crc := t.InitCrc()
	crc = t.UpdateCrc(crc, data)
	return t.CRC(crc)
}

// Hash represents the partial evaluation of a checksum using table-driven
// implementation. It also implements hash.Hash interface.
type Hash struct {
	table    *Table
	curValue uint64
	size     uint
}

// Size returns the number of bytes Sum will return.
// See hash.Hash interface.
func (h *Hash) Size() int { return int(h.size) }

// BlockSize returns the hash's underlying block size.
// The Write method must be able to accept any amount
// of data, but it may operate more efficiently if all writes
// are a multiple of the block size.
// See hash.Hash interface.
func (h *Hash) BlockSize() int { return 1 }

// Reset resets the Hash to its initial state.
// See hash.Hash interface.
func (h *Hash) Reset() {
	h.curValue = h.table.InitCrc()
}

// Sum appends the current hash to b and returns the resulting slice.
// It does not change the underlying hash state.
// See hash.Hash interface.
func (h *Hash) Sum(in []byte) []byte {
	s := h.CRC()
	for i := h.size; i > 0; {
		i--
		in = append(in, byte(s>>(8*i)))
	}
	return in
}

// Write implements io.Writer interface which is part of hash.Hash interface.
func (h *Hash) Write(p []byte) (n int, err error) {
	h.Update(p)
	return len(p), nil
}

// Update updates process supplied bytes and updates current (partial) CRC accordingly.
func (h *Hash) Update(p []byte) {
	h.curValue = h.table.UpdateCrc(h.curValue, p)
}

// CRC returns current CRC value for the data processed so far.
func (h *Hash) CRC() uint64 {
	return h.table.CRC(h.curValue)
}

// CalculateCRC is a convenience function allowing to calculate CRC in one call.
func (h *Hash) CalculateCRC(data []byte) uint64 {
	return h.table.CalculateCRC(data)
}

// NewHashWithTable creates a new Hash instance configured for table driven
// CRC calculation using a Table instance created elsewhere.
func NewHashWithTable(table *Table) *Hash {
	ret := &Hash{table: table}
	ret.size = (table.crcParams.Width + 7) / 8 // smalest number of bytes enough to store produced crc
	ret.Reset()
	return ret
}

// NewHash creates a new Hash instance configured for table driven
// CRC calculation according to parameters specified.
func NewHash(crcParams *Parameters) *Hash {
	return NewHashWithTable(NewTable(crcParams))
}

// CRC8 is a convenience method to spare end users from explicit type conversion every time this package is used.
// Underneath, it just calls CRC() method.
func (h *Hash) CRC8() uint8 {
	return h.table.CRC8(h.curValue)
}

// CRC16 is a convenience method to spare end users from explicit type conversion every time this package is used.
// Underneath, it just calls CRC() method.
func (h *Hash) CRC16() uint16 {
	return h.table.CRC16(h.curValue)
}

// CRC32 is a convenience method to spare end users from explicit type conversion every time this package is used.
// Underneath, it just calls CRC() method.
func (h *Hash) CRC32() uint32 {
	return h.table.CRC32(h.curValue)
}

// Table used by this Hash under the hood
func (h *Hash) Table() *Table {
	return h.table
}
