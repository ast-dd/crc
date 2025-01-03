package crc

import (
	"fmt"
	"strings"
)

var parametersMap = map[string]*Parameters{
	"CRC8":         CRC8,
	"CRC8CDMA2000": CRC8CDMA2000,
	"CRC8DARC":     CRC8DARC,
	"CRC8DVBS2":    CRC8DVBS2,
	"CRC8EBU":      CRC8EBU,
	"CRC8ICODE":    CRC8ICODE,
	"CRC8ITU":      CRC8ITU,
	"CRC8MAXIM":    CRC8MAXIM,
	"CRC8ROHC":     CRC8ROHC,
	"CRC8WCDMA":    CRC8WCDMA,

	"CRC8AUTOSAR":    CRC8AUTOSAR,
	"CRC8BLUETOOTH":  CRC8BLUETOOTH,
	"CRC8GSMA":       CRC8GSMA,
	"CRC8GSMB":       CRC8GSMB,
	"CRC8HITAG":      CRC8HITAG,
	"CRC8LTE":        CRC8LTE,
	"CRC8MIFAREMAD":  CRC8MIFAREMAD,
	"CRC8NRSC5":      CRC8NRSC5,
	"CRC8OPENSAFETY": CRC8OPENSAFETY,
	"CRC8SAEJ1850":   CRC8SAEJ1850,

	"CRC16ARC":        CRC16ARC,
	"CRC16AUGCCITT":   CRC16AUGCCITT,
	"CRC16BUYPASS":    CRC16BUYPASS,
	"CRC16CCITTFALSE": CRC16CCITTFALSE,
	"CCITT":           CCITT,
	"CRC16CDMA2000":   CRC16CDMA2000,
	"CRC16DDS110":     CRC16DDS110,
	"CRC16DECTR":      CRC16DECTR,
	"CRC16DECTX":      CRC16DECTX,
	"CRC16DNP":        CRC16DNP,
	"CRC16EN13757":    CRC16EN13757,
	"CRC16GENIBUS":    CRC16GENIBUS,
	"CRC16KERMIT":     CRC16KERMIT,
	"CRC16MAXIM":      CRC16MAXIM,
	"CRC16MCRF4XX":    CRC16MCRF4XX,
	"CRC16MODBUS":     CRC16MODBUS,
	"CRC16RIELLO":     CRC16RIELLO,
	"CRC16T10DIF":     CRC16T10DIF,
	"CRC16TELEDISK":   CRC16TELEDISK,
	"CRC16TMS37157":   CRC16TMS37157,
	"CRC16USB":        CRC16USB,
	"CRC16X25":        CRC16X25,
	"X25":             X25,
	"CRC16XMODEM":     CRC16XMODEM,
	"XMODEM2":         XMODEM2,
	"CRCA":            CRCA,

	"CRC16CMS":         CRC16CMS,
	"CRC16GSM":         CRC16GSM,
	"CRC16LJ1200":      CRC16LJ1200,
	"CRC16M17":         CRC16M17,
	"CRC16NRSC5":       CRC16NRSC5,
	"CRC16OPENSAFETYA": CRC16OPENSAFETYA,
	"CRC16OPENSAFETYB": CRC16OPENSAFETYB,
	"CRC16PROFIBUS":    CRC16PROFIBUS,

	"CRC32":       CRC32,
	"IEEE":        IEEE,
	"CRC32BZIP2":  CRC32BZIP2,
	"CRC32JAMCRC": CRC32JAMCRC,
	"CRC32MPEG2":  CRC32MPEG2,
	"CRC32POSIX":  CRC32POSIX,
	"CRC32SATA":   CRC32SATA,
	"CRC32XFER":   CRC32XFER,
	"CRC32C":      CRC32C,
	"Castagnoli":  Castagnoli,
	"CRC32D":      CRC32D,
	"CRC32Q":      CRC32Q,

	"CRC32AUTOSAR":  CRC32AUTOSAR,
	"CRC32CDROMEDC": CRC32CDROMEDC,
	"CRC32MEF":      CRC32MEF,

	"Koopman": Koopman,

	"CRC64ISO":  CRC64ISO,
	"CRC64ECMA": CRC64ECMA,
}

// GetParameters returns the CRC parameters for given string s
func GetParameters(s string) (parameters *Parameters, err error) {
	var ok bool
	s = strings.ToUpper(s)
	if parameters, ok = parametersMap[s]; !ok {
		err = fmt.Errorf("unknown CRC type: %q", s)
		return
	}
	return
}

// GetParametersName returns the name for given CRC parameters by checking the pointer
func GetParametersName(parameters *Parameters) (name string, err error) {
	for s, p := range parametersMap {
		if p == parameters {
			name = s
			return
		}
	}
	err = fmt.Errorf("parameters not from known list")
	return
}
