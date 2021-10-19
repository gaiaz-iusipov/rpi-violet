package co2mon_test

import (
	"testing"

	"github.com/gaiaz-iusipov/rpi-violet/internal/monitor/co2mon"
)

func Test_decrypt(t *testing.T) {
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
			name: "with empty key",
			args: args{
				data: [8]byte{0x35, 0xa4, 0x32, 0xb6, 0xcf, 0x9a, 0x9c, 0xd0},
			},
			want: [8]byte{0x42, 0x12, 0x90, 0xe4, 0x0d, 0x00, 0x00, 0x00},
		},
		{
			name: "common",
			args: args{
				data: [8]byte{0x35, 0xa4, 0x32, 0xb6, 0xcf, 0x9a, 0x9c, 0xd0},
				key: [8]byte{1, 2, 1, 4, 3, 2, 5, 12},
			},
			want: [8]byte{194, 50, 80, 196, 141, 96, 64, 161},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := co2mon.Decrypt(tt.args.data, tt.args.key)

			if got != tt.want {
				t.Errorf("Decrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}
