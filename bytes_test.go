package crc_test

import (
	"reflect"
	"testing"

	"github.com/ast-dd/crc"
)

var (
	byteTestStrings = [4]string{
		"123456789",
		"12345678901234567890",
		"Introduction on CRC calculations",
		"Whenever digital data is stored or interfaced, data corruption might occur. Since the beginning of computer science, people have been thinking of ways to deal with this type of problem. For serial data they came up with the solution to attach a parity bit to each sent byte. This simple detection mechanism works if an odd number of bits in a byte changes, but an even number of false bits in one byte will not be detected by the parity check. To overcome this problem people have searched for mathematical sound mechanisms to detect multiple false bits.",
	}

	byteTests = []struct {
		name      string
		crcParams *crc.Parameters
		wantBytes [4][]byte
	}{
		{"CRC8SAEJ1850", crc.CRC8SAEJ1850, [4][]byte{{0x4B}, {0x91}, {0x8D}, {0x41}}},
		{"CRC16MODBUS", crc.CRC16MODBUS, [4][]byte{{0x37, 0x4B}, {0x13, 0x80}, {0x8B, 0xE6}, {0xDF, 0xFF}}},
		{"CRC32", crc.CRC32, [4][]byte{{0x26, 0x39, 0xF4, 0xCB}, {0xF2, 0x19, 0x63, 0x90}, {0x45, 0x2B, 0x4F, 0x81}, {0x17, 0x38, 0x27, 0x8F}}},
	}
)

func TestCalculateCRCBytes(t *testing.T) {
	for _, tt := range byteTests {
		t.Run(tt.name, func(t *testing.T) {
			for i, testString := range byteTestStrings {
				data := []byte(testString)
				if got, want := crc.CalculateCRCBytes(tt.crcParams, data), tt.wantBytes[i]; !reflect.DeepEqual(got, want) {
					t.Errorf("CalculateCRCBytes(%q) = %v, want %v", testString, got, want)
				}
			}
		})
	}

	t.Run("modbus", func(t *testing.T) {
		data := []byte{3, 0x04, 0, 2, 0, 1}
		got, want := crc.CalculateCRCBytes(crc.CRC16MODBUS, data), []byte{0x91, 0xe8}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("CalculateCRCBytes(%#v) = %#v, want %#v", data, got, want)
		}
	})
}

func TestAppendCRCBytes(t *testing.T) {
	for _, tt := range byteTests {
		t.Run(tt.name, func(t *testing.T) {
			for i, testString := range byteTestStrings {
				data := []byte(testString)

				got, want := crc.AppendCRCBytes(tt.crcParams, data), append(data, tt.wantBytes[i]...)
				if !reflect.DeepEqual(got, want) {
					t.Errorf("CalculateCRCBytes(%q) = %v, want %v", testString, got, want)
				}
			}
		})
	}
}

func TestCheckCRCBytes(t *testing.T) {
	for _, tt := range byteTests {
		t.Run(tt.name, func(t *testing.T) {
			for i, testString := range byteTestStrings {
				var want bool
				data := []byte(testString)

				// correct
				correct := tt.wantBytes[i]
				got := crc.CheckCRCBytes(tt.crcParams, data, correct)
				if want = true; got != want {
					t.Errorf("CheckCRCBytes(%q, %#v) = %v, want %v", testString, correct, got, want)
				}

				// incorrect
				incorrect := make([]byte, len(correct))
				copy(incorrect, correct)
				incorrect[0] += 1
				got = crc.CheckCRCBytes(tt.crcParams, data, incorrect)
				if want = false; got != want {
					t.Errorf("CheckCRCBytes(%q, %#v) = %v, want %v", testString, incorrect, got, want)
				}

				// checksum too short
				short := make([]byte, len(correct)-1)
				copy(short, correct)
				got = crc.CheckCRCBytes(tt.crcParams, data, short)
				if want = false; got != want {
					t.Errorf("CheckCRCBytes(%q, %#v) = %v, want %v", testString, short, got, want)
				}

				// checksum too long
				long := make([]byte, len(correct)+1)
				copy(long, correct)
				long = append(long, 0x00)
				got = crc.CheckCRCBytes(tt.crcParams, data, long)
				if want = false; got != want {
					t.Errorf("CheckCRCBytes(%q, %#v) = %v, want %v", testString, long, got, want)
				}
			}
		})
	}
}
