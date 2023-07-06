package crc

import (
	"hash"
	"testing"
)

func TestCRCAlgorithms(t *testing.T) {

	doTest := func(crcParams *Parameters, data string, crc uint64) {
		calculated := CalculateCRC(crcParams, []byte(data))
		if calculated != crc {
			t.Errorf("Incorrect CRC 0x%04x calculated for %s (should be 0x%04x)", calculated, data, crc)
		}

		// same test using table driven
		tableDriven := NewHash(crcParams)
		calculated = tableDriven.CalculateCRC([]byte(data))
		if calculated != crc {
			t.Errorf("Incorrect CRC 0x%04x calculated for %s (should be 0x%04x)", calculated, data, crc)
		}

		// same test feeding data in chunks of different size
		tableDriven.Reset()
		var start = 0
		var step = 1
		for start < len(data) {
			end := start + step
			if end > len(data) {
				end = len(data)
			}
			tableDriven.Update([]byte(data[start:end]))
			start = end
			step *= 2
		}
		calculated = tableDriven.CRC()
		if calculated != crc {
			t.Errorf("Incorrect CRC 0x%04x calculated for %s (should be 0x%04x)", calculated, data, crc)
		}

		// Test helper methods return correct values as well
		if crcParams.Width == 8 {
			crc8 := tableDriven.CRC8()
			if crc8 != uint8(crc&0x00FF) {
				t.Errorf("Incorrect CRC8 0x%02x retrived %s (should be 0x%02x)", crc8, data, crc)
			}
		} else if crcParams.Width == 16 {
			crc16 := tableDriven.CRC16()
			if crc16 != uint16(crc&0x00FFFF) {
				t.Errorf("Incorrect CRC16 0x%04x retrived %s (should be 0x%04x)", crc16, data, crc)
			}
		} else if crcParams.Width == 32 {
			crc32 := tableDriven.CRC32()
			if crc32 != uint32(crc&0x00FFFFFFFF) {
				t.Errorf("Incorrect CRC8 0x%08x retrived %s (should be 0x%08x)", crc32, data, crc)
			}
		}

		// Test Hash's table directly and see there is no difference
		table := tableDriven.Table()
		calculated = table.CalculateCRC([]byte(data))
		if calculated != crc {
			t.Errorf("Incorrect CRC 0x%04x calculated for %s (should be 0x%04x)", calculated, data, crc)
		}

	}

	testStrings := [4]string{
		"123456789",
		"12345678901234567890",
		"Introduction on CRC calculations",
		"Whenever digital data is stored or interfaced, data corruption might occur. Since the beginning of computer science, people have been thinking of ways to deal with this type of problem. For serial data they came up with the solution to attach a parity bit to each sent byte. This simple detection mechanism works if an odd number of bits in a byte changes, but an even number of false bits in one byte will not be detected by the parity check. To overcome this problem people have searched for mathematical sound mechanisms to detect multiple false bits.",
	}

	type testVectors struct {
		algo *Parameters
		crc  [4]uint64 // the expected CRC for each of the 4 test strings
	}

	crcTests := []testVectors{
		{algo: CRC8, crc: [4]uint64{0xF4, 0x56, 0x58, 0x83}},
		{algo: CRC8CDMA2000, crc: [4]uint64{0xDA, 0x36, 0xAB, 0xD3}},
		{algo: CRC8DARC, crc: [4]uint64{0x15, 0x1B, 0xBF, 0x3D}},
		{algo: CRC8DVBS2, crc: [4]uint64{0xBC, 0x99, 0xA4, 0x86}},
		{algo: CRC8EBU, crc: [4]uint64{0x97, 0x7B, 0x72, 0xE1}},
		{algo: CRC8ICODE, crc: [4]uint64{0x7E, 0xBF, 0x76, 0xFD}},
		{algo: CRC8ITU, crc: [4]uint64{0xA1, 0x03, 0x0D, 0xD6}},
		{algo: CRC8MAXIM, crc: [4]uint64{0xA1, 0x18, 0x45, 0x3D}},
		{algo: CRC8ROHC, crc: [4]uint64{0xD0, 0xA6, 0x1F, 0x7A}},
		{algo: CRC8WCDMA, crc: [4]uint64{0x25, 0x95, 0x61, 0x64}},
		{algo: CRC8AUTOSAR, crc: [4]uint64{0xDF, 0x74, 0xAB, 0x7A}},
		{algo: CRC8BLUETOOTH, crc: [4]uint64{0x26, 0xE4, 0xA8, 0x8C}},
		{algo: CRC8GSMA, crc: [4]uint64{0x37, 0, 0, 0}},
		{algo: CRC8GSMB, crc: [4]uint64{0x94, 0x72, 0xA4, 0xCE}},
		{algo: CRC8HITAG, crc: [4]uint64{0xB4, 0, 0, 0}},
		{algo: CRC8LTE, crc: [4]uint64{0xEA, 0, 0, 0}},
		{algo: CRC8MIFAREMAD, crc: [4]uint64{0x99, 0, 0, 0}},
		{algo: CRC8NRSC5, crc: [4]uint64{0xF7, 0, 0, 0}},
		{algo: CRC8OPENSAFETY, crc: [4]uint64{0x3E, 0, 0, 0}},
		{algo: CRC8SAEJ1850, crc: [4]uint64{0x4B, 0x91, 0x8D, 0x41}},

		{algo: CRC16ARC, crc: [4]uint64{0xBB3D, 0x9B37, 0x728A, 0x1EB5}},
		{algo: CRC16AUGCCITT, crc: [4]uint64{0xE5CC, 0xB33D, 0x908A, 0x0209}},
		{algo: CRC16BUYPASS, crc: [4]uint64{0xFEE8, 0xC2F4, 0xAC80, 0x4688}},
		{algo: CRC16CCITTFALSE, crc: [4]uint64{0x29B1, 0xDA31, 0xC87E, 0xD6ED}},
		{algo: CRC16CDMA2000, crc: [4]uint64{0x4C06, 0xBC63, 0xB4DD, 0x8B13}},
		{algo: CRC16DDS110, crc: [4]uint64{0x9ECF, 0x1824, 0xAC7C, 0x319F}},
		{algo: CRC16DECTR, crc: [4]uint64{0x007E, 0x6F58, 0x5240, 0xE9D5}},
		{algo: CRC16DECTX, crc: [4]uint64{0x007F, 0x6F59, 0x5241, 0xE9D4}},
		{algo: CRC16DNP, crc: [4]uint64{0xEA82, 0x745F, 0x035B, 0x8933}},
		{algo: CRC16EN13757, crc: [4]uint64{0xC2B7, 0xAD9D, 0xCD00, 0xA5FD}},
		{algo: CRC16GENIBUS, crc: [4]uint64{0xD64E, 0x25CE, 0x3781, 0x2912}},
		{algo: CRC16KERMIT, crc: [4]uint64{0x2189, 0x4016, 0x34C6, 0x4157}},
		{algo: CRC16MAXIM, crc: [4]uint64{0x44C2, 0x64C8, 0x8D75, 0xE14A}},
		{algo: CRC16MCRF4XX, crc: [4]uint64{0x6F91, 0x5D79, 0x0649, 0x974E}},
		{algo: CRC16MODBUS, crc: [4]uint64{0x4B37, 0x8013, 0xE68B, 0xFFDF}},
		{algo: CRC16RIELLO, crc: [4]uint64{0x63D0, 0x1C58, 0x3038, 0x44F2}},
		{algo: CRC16T10DIF, crc: [4]uint64{0xD0DB, 0x1A0D, 0xCD63, 0xBF32}},
		{algo: CRC16TELEDISK, crc: [4]uint64{0x0FB3, 0xC503, 0x776E, 0xA357}},
		{algo: CRC16TMS37157, crc: [4]uint64{0x26B1, 0xAB86, 0x7053, 0x2533}},
		{algo: CRC16USB, crc: [4]uint64{0xB4C8, 0x7FEC, 0x1974, 0x0020}},
		{algo: CRC16X25, crc: [4]uint64{0x906E, 0xA286, 0xF9B6, 0x68B1}},
		{algo: CRC16XMODEM, crc: [4]uint64{0x31C3, 0x2C89, 0x3932, 0x4E86}},
		{algo: XMODEM2, crc: [4]uint64{0x0C73, 0x122E, 0x0638, 0x187A}},
		{algo: CRCA, crc: [4]uint64{0xBF05, 0xC1AD, 0xEEE6, 0x7A71}},

		{algo: CRC16CMS, crc: [4]uint64{0xAEE7, 0, 0, 0}},
		{algo: CRC16GSM, crc: [4]uint64{0xCE3C, 0, 0, 0}},
		{algo: CRC16LJ1200, crc: [4]uint64{0xBDF4, 0, 0, 0}},
		{algo: CRC16M17, crc: [4]uint64{0x772B, 0, 0, 0}},
		{algo: CRC16NRSC5, crc: [4]uint64{0xA066, 0, 0, 0}},
		{algo: CRC16OPENSAFETYA, crc: [4]uint64{0x5D38, 0, 0, 0}},
		{algo: CRC16OPENSAFETYB, crc: [4]uint64{0x20FE, 0, 0, 0}},
		{algo: CRC16PROFIBUS, crc: [4]uint64{0xA819, 0x3D0B, 0x47A, 0xFF77}},

		{algo: CRC32, crc: [4]uint64{0xCBF43926, 0x906319F2, 0x814F2B45, 0x8F273817}},
		{algo: CRC32BZIP2, crc: [4]uint64{0xFC891918, 0x47D3AC25, 0x05ACC4EE, 0xBB6BB3F9}},
		{algo: CRC32JAMCRC, crc: [4]uint64{0x340BC6D9, 0x6F9CE60D, 0x7EB0D4BA, 0x70D8C7E8}},
		{algo: CRC32MPEG2, crc: [4]uint64{0x0376E6E7, 0xB82C53DA, 0xFA533B11, 0x44944C06}},
		{algo: CRC32POSIX, crc: [4]uint64{0x765E7680, 0x09F5F82A, 0x4FF96B89, 0x54E6B1BB}},
		{algo: CRC32SATA, crc: [4]uint64{0xCF72AFE8, 0x4426E87D, 0xCA7E8932, 0x7D81662A}},
		{algo: CRC32XFER, crc: [4]uint64{0xBD0BE338, 0x326A653D, 0x5D470C48, 0x9001FA5D}},
		{algo: CRC32C, crc: [4]uint64{0xE3069283, 0xA8B4A6B9, 0x54F98A9E, 0x864FDAFC}},
		{algo: CRC32D, crc: [4]uint64{0x87315576, 0xF2D28B69, 0x231CB16A, 0xBEB35C2D}},
		{algo: CRC32Q, crc: [4]uint64{0x3010BF7F, 0x59F2D11A, 0xDBE0EEB1, 0x4ED076AE}},
		{algo: Koopman, crc: [4]uint64{0x2D3DD0AE, 0xCC53DEAC, 0x1B8101F9, 0xA41634B2}},

		{algo: CRC32AUTOSAR, crc: [4]uint64{0x1697d06a, 0, 0, 0}},
		{algo: CRC32CDROMEDC, crc: [4]uint64{0x6ec2edc4, 0, 0, 0}},
		{algo: CRC32MEF, crc: [4]uint64{0xd2c22f51, 0, 0, 0}},

		{algo: CRC64ISO, crc: [4]uint64{0xB90956C775A41001, 0x8DB93749FB37B446, 0xBAA81A1ED1A9209B, 0x347969424A1A7628}},
		{algo: CRC64ECMA, crc: [4]uint64{0x995DC9BBDF1939FA, 0x0DA1B82EF5085A4A, 0xCF8C40119AE90DCB, 0x31610F76CFB272A5}},
	}

	for _, v := range crcTests {
		for i := 0; i < len(testStrings); i++ {
			if v.crc[i] != 0 {
				doTest(v.algo, testStrings[i], v.crc[i])
			}
		}
	}

	// More tests for various CRC algorithms (copied from java version)
	longText := "Whenever digital data is stored or interfaced, data corruption might occur. Since the beginning of computer science, people have been thinking of ways to deal with this type of problem. For serial data they came up with the solution to attach a parity bit to each sent byte. This simple detection mechanism works if an odd number of bits in a byte changes, but an even number of false bits in one byte will not be detected by the parity check. To overcome this problem people have searched for mathematical sound mechanisms to detect multiple false bits."

	testArrayData := make([]byte, 256)
	for i := 0; i < len(testArrayData); i++ {
		testArrayData[i] = byte(i & 0x0FF)
	}
	testArray := string(testArrayData)
	if len(testArray) != 256 {
		t.Fatalf("Logic error")
	}

	// merely a helper to make copying Spock test sets from java version of this library a bit easier
	doTestWithParameters := func(width uint, polynomial uint64, init uint64, reflectIn bool, reflectOut bool, finalXor uint64, crc uint64, testData string) {
		doTest(&Parameters{Width: width, Polynomial: polynomial, Init: init, ReflectIn: reflectIn, ReflectOut: reflectOut, FinalXor: finalXor}, testData, crc)
	}

	doTestWithParameters(3, 0x03, 0x00, false, false, 0x7, 0x04, "123456789") // CRC-3/GSM
	doTestWithParameters(3, 0x03, 0x00, false, false, 0x7, 0x06, longText)
	doTestWithParameters(3, 0x03, 0x00, false, false, 0x7, 0x02, testArray)
	doTestWithParameters(3, 0x03, 0x07, true, true, 0x0, 0x06, "123456789") // CRC-3/ROHC
	doTestWithParameters(3, 0x03, 0x07, true, true, 0x0, 0x03, longText)
	doTestWithParameters(4, 0x03, 0x00, true, true, 0x0, 0x07, "123456789")   // CRC-4/ITU
	doTestWithParameters(4, 0x03, 0x0f, false, false, 0xf, 0x0b, "123456789") // CRC-4/INTERLAKEN
	doTestWithParameters(4, 0x03, 0x0f, false, false, 0xf, 0x01, longText)    // CRC-4/INTERLAKEN
	doTestWithParameters(4, 0x03, 0x0f, false, false, 0xf, 0x07, testArray)   // CRC-4/INTERLAKEN
	doTestWithParameters(5, 0x09, 0x09, false, false, 0x0, 0x00, "123456789") // CRC-5/EPC
	doTestWithParameters(5, 0x15, 0x00, true, true, 0x0, 0x07, "123456789")   // CRC-5/ITU
	doTestWithParameters(6, 0x27, 0x3f, false, false, 0x0, 0x0d, "123456789") // CRC-6/CDMA2000-A
	doTestWithParameters(6, 0x07, 0x3f, false, false, 0x0, 0x3b, "123456789") // CRC-6/CDMA2000-B
	doTestWithParameters(6, 0x07, 0x3f, false, false, 0x0, 0x24, testArray)   // CRC-6/CDMA2000-B
	doTestWithParameters(7, 0x09, 0x00, false, false, 0x0, 0x75, "123456789") // CRC-7
	doTestWithParameters(7, 0x09, 0x00, false, false, 0x0, 0x78, testArray)   // CRC-7
	doTestWithParameters(7, 0x4f, 0x7f, true, true, 0x0, 0x53, "123456789")   // CRC-7/ROHC

	doTestWithParameters(8, 0x07, 0x00, false, false, 0x00, 0xf4, "123456789") // CRC-8
	doTestWithParameters(8, 0xa7, 0x00, true, true, 0x00, 0x26, "123456789")   // CRC-8/BLUETOOTH
	doTestWithParameters(8, 0x07, 0x00, false, false, 0x55, 0xa1, "123456789") // CRC-8/ITU
	doTestWithParameters(8, 0x9b, 0x00, true, true, 0x00, 0x25, "123456789")   // CRC-8/WCDMA
	doTestWithParameters(8, 0x31, 0x00, true, true, 0x00, 0xa1, "123456789")   // CRC-8/MAXIM

	doTestWithParameters(10, 0x233, 0x000, false, false, 0x000, 0x199, "123456789") // CRC-10

	doTestWithParameters(12, 0xd31, 0x00, false, false, 0xfff, 0x0b34, "123456789")   // CRC-12/GSM
	doTestWithParameters(12, 0x80f, 0x00, false, true, 0x00, 0x0daf, "123456789")     // CRC-12/UMTS
	doTestWithParameters(13, 0x1cf5, 0x00, false, false, 0x00, 0x04fa, "123456789")   // CRC-13/BBC
	doTestWithParameters(14, 0x0805, 0x00, true, true, 0x00, 0x082d, "123456789")     // CRC-14/DARC
	doTestWithParameters(14, 0x202d, 0x00, false, false, 0x3fff, 0x30ae, "123456789") // CRC-14/GSM

	doTestWithParameters(15, 0x4599, 0x00, false, false, 0x00, 0x059e, "123456789") // CRC-15
	doTestWithParameters(15, 0x4599, 0x00, false, false, 0x00, 0x2857, longText)
	doTestWithParameters(15, 0x6815, 0x00, false, false, 0x0001, 0x2566, "123456789") // CRC-15/MPT1327

	doTestWithParameters(21, 0x102899, 0x000000, false, false, 0x000000, 0x0ed841, "123456789") // CRC-21/CAN-FD
	doTestWithParameters(24, 0x864cfb, 0xb704ce, false, false, 0x000000, 0x21cf02, "123456789") // CRC-24
	doTestWithParameters(24, 0x5d6dcb, 0xfedcba, false, false, 0x000000, 0x7979bd, "123456789") // CRC-24/FLEXRAY-A
	doTestWithParameters(24, 0x00065b, 0x555555, true, true, 0x000000, 0xc25a56, "123456789")   // "CRC-24/BLE"

	doTestWithParameters(31, 0x04c11db7, 0x7fffffff, false, false, 0x7fffffff, 0x0ce9e46c, "123456789") // CRC-31/PHILIPS
}

func TestSizeMethods(t *testing.T) {
	testWidth := func(width uint, expectedSize int) {
		h := NewHash(&Parameters{Width: width, Polynomial: 1})
		s := h.Size()
		if s != expectedSize {
			t.Errorf("Incorrect Size calculated for width %d:  %d when should be %d", width, s, expectedSize)
		}
		bs := h.BlockSize()
		if bs != 1 {
			t.Errorf("Incorrect Block Size returned for width %d:  %d when should always be 1", width, bs)
		}
	}

	testWidth(3, 1)
	testWidth(8, 1)
	testWidth(12, 2)
	testWidth(16, 2)
	testWidth(32, 4)
	testWidth(64, 8)

}

func TestHashInterface(t *testing.T) {
	doTest := func(crcParams *Parameters, data string, crc uint64) {
		// same test using table driven
		var h hash.Hash = NewHash(crcParams)

		// same test feeding data in chunks of different size
		h.Reset()
		var start = 0
		var step = 1
		for start < len(data) {
			end := start + step
			if end > len(data) {
				end = len(data)
			}
			h.Write([]byte(data[start:end]))
			start = end
			step *= 2
		}

		buf := make([]byte, 0, 0)
		buf = h.Sum(buf)

		if len(buf) != h.Size() {
			t.Errorf("Wrong number of bytes appended by Sum(): %d when should be %d", len(buf), h.Size())
		}

		calculated := uint64(0)
		for _, b := range buf {
			calculated <<= 8
			calculated += uint64(b)
		}

		if calculated != crc {
			t.Errorf("Incorrect CRC 0x%04x calculated for %s (should be 0x%04x)", calculated, data, crc)
		}
	}

	doTest(&Parameters{Width: 8, Polynomial: 0x07, Init: 0x00, ReflectIn: false, ReflectOut: false, FinalXor: 0x00}, "123456789", 0xf4)
	doTest(CCITT, "12345678901234567890", 0xDA31)
	doTest(CRC64ECMA, "Introduction on CRC calculations", 0xCF8C40119AE90DCB)
	doTest(CRC32C, "Whenever digital data is stored or interfaced, data corruption might occur. Since the beginning of computer science, people have been thinking of ways to deal with this type of problem. For serial data they came up with the solution to attach a parity bit to each sent byte. This simple detection mechanism works if an odd number of bits in a byte changes, but an even number of false bits in one byte will not be detected by the parity check. To overcome this problem people have searched for mathematical sound mechanisms to detect multiple false bits.", 0x864FDAFC)
}

func BenchmarkCCITT(b *testing.B) {
	data := []byte("Whenever digital data is stored or interfaced, data corruption might occur. Since the beginning of computer science, people have been thinking of ways to deal with this type of problem. For serial data they came up with the solution to attach a parity bit to each sent byte. This simple detection mechanism works if an odd number of bits in a byte changes, but an even number of false bits in one byte will not be detected by the parity check. To overcome this problem people have searched for mathematical sound mechanisms to detect multiple false bits.")
	for i := 0; i < b.N; i++ {
		tableDriven := NewHash(CCITT)
		tableDriven.Update(data)
		tableDriven.CRC()
	}
}
