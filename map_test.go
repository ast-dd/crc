package crc_test

import (
	"reflect"
	"testing"

	"github.com/ast-dd/crc"
)

func TestGetParameters(t *testing.T) {
	tests := []struct {
		name           string
		s              string
		wantParameters *crc.Parameters
		wantErr        bool
	}{
		{"empty", "", nil, true},
		{"unknown", "CRCsomeUnknown", nil, true},
		{"CRC8SAEJ1850", "CRC8SAEJ1850", crc.CRC8SAEJ1850, false},
		{"Crc8Saej1850", "Crc8Saej1850", crc.CRC8SAEJ1850, false},
		{"CRC64ECMA", "CRC64ECMA", crc.CRC64ECMA, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotParameters, err := crc.GetParameters(tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetParameters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotParameters, tt.wantParameters) {
				t.Errorf("GetParameters() gotParameters = %v, want %v", gotParameters, tt.wantParameters)
			}
		})
	}
}
