package co2mon_test

import (
	"testing"

	"github.com/gaiaz-iusipov/rpi-violet/internal/monitor/co2mon"
)

func TestDecrypt(t *testing.T) {
	type args struct {
		data [8]byte
		key  [8]byte
	}
	tests := []struct {
		name string
		args args
		want [8]byte
	}{
		{
			name: "without key",
			args: args{
				data: [8]byte{0x35, 0xa4, 0x32, 0xb6, 0xcf, 0x9a, 0x9c, 0xd0},
			},
			want: [8]byte{0x42, 0x12, 0x90, 0xe4, 0x0d, 0x00, 0x00, 0x00},
		},
		{
			name: "with key",
			args: args{
				data: [8]byte{0x35, 0xa4, 0x32, 0xb6, 0xcf, 0x9a, 0x9c, 0xd0},
				key:  [8]byte{0xd3, 0x89, 0xc9, 0xc8, 0xa4, 0x9c, 0x9a, 0xb6},
			},
			want: [8]byte{0x98, 0xe1, 0x89, 0xad, 0xf9, 0x6d, 0x6d, 0xaa},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := co2mon.Decrypt(tt.args.data, tt.args.key)

			if got != tt.want {
				t.Errorf("Decrypt() = %# x, want %# x", got, tt.want)
			}
		})
	}
}
