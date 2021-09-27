package co2mon

import (
	"testing"
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
			name: "equal",
			args: args{
				data: [8]byte{0x35, 0xa4, 0x32, 0xb6, 0xcf, 0x9a, 0x9c, 0xd0},
			},
			want: [8]byte{0x42, 0x12, 0x90, 0xe4, 0x0d, 0x00, 0x00, 0x00},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := decrypt(tt.args.data, tt.args.key)

			if got != tt.want {
				t.Errorf("decrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_phase1(t *testing.T) {
	tests := []struct {
		name       string
		orig, want [8]byte
	}{
		{
			name: "equal",
			orig: [8]byte{0xc9, 0xa4, 0xd3, 0xb6, 0x89, 0x9a, 0x9c, 0xc8},
			want: [8]byte{0xd3, 0x89, 0xc9, 0xc8, 0xa4, 0x9c, 0x9a, 0xb6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := phase1(tt.orig)

			if tt.want != got {
				t.Errorf("phase1() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_phase3(t *testing.T) {
	tests := []struct {
		name       string
		orig, want [8]byte
	}{
		{
			name: "equal",
			orig: [8]byte{0xd3, 0x89, 0xc9, 0xc8, 0xa4, 0x9c, 0x9a, 0xb6},
			want: [8]byte{0xda, 0x71, 0x39, 0x39, 0x14, 0x93, 0x93, 0x56},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := phase3(tt.orig)

			if tt.want != got {
				t.Errorf("phase3() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_phase4(t *testing.T) {
	tests := []struct {
		name       string
		orig, want [8]byte
	}{
		{
			name: "equal",
			orig: [8]byte{0xf1, 0x4e, 0x1c, 0x10, 0x14, 0x93, 0x93, 0x56},
			want: [8]byte{0x6d, 0x07, 0xc6, 0x3a, 0x0d, 0x00, 0x00, 0x00},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := phase4(tt.orig)

			if tt.want != got {
				t.Errorf("phase4() = %v, want %v", got, tt.want)
			}
		})
	}
}
