package co2mon

import (
	"errors"
	"fmt"
)

func decrypt(data, key [8]byte) [8]byte {
	data = phase1(data)
	data = phase2(data, key)
	data = phase3(data)
	data = phase4(data)
	return data
}

var swaps = map[byte]byte{0: 2, 1: 4, 3: 7, 5: 6}

func phase1(b [8]byte) [8]byte {
	for l, r := range swaps {
		b[l], b[r] = b[r], b[l]
	}
	return b
}

func phase2(b, key [8]byte) [8]byte {
	for i := range b {
		b[i] ^= key[i]
	}
	return b
}

func phase3(b [8]byte) [8]byte {
	tmp := b[7] << 5
	for i := 7; i > 0; i-- {
		b[i] = b[i-1]<<5 | b[i]>>3
	}
	b[0] = tmp | b[0]>>3
	return b
}

func phase4(b [8]byte) [8]byte {
	const magicWord = "Htemp99e"
	for i := range b {
		b[i] -= magicWord[i]<<4 | magicWord[i]>>4
	}
	return b
}

var validTail = [...]byte{0x0d, 0x00, 0x00, 0x00}

func validate(b [8]byte) error {
	var tail [4]byte
	copy(tail[:], b[4:])
	if tail != validTail {
		return fmt.Errorf("unexpected tail: got %v, want %v", tail, validTail)
	}

	if b[0]+b[1]+b[2] != b[3] {
		return errors.New("head checksum")
	}
	return nil
}
