package co2mon

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

type Device struct {
	file *os.File
	opts *options
}

func Open(setters ...OptionSetter) (*Device, error) {
	opts := newOptions(setters)

	file, err := os.OpenFile(opts.path, os.O_RDWR, 0)
	if err != nil {
		return nil, fmt.Errorf("os.OpenFile(): %w", err)
	}

	data := [9]byte{}
	copy(data[1:], opts.key[:]) // First byte needs to be 0x00

	const hidiocsfeature9 uintptr = 0xc0094806
	_, _, errNo := syscall.Syscall(syscall.SYS_IOCTL, file.Fd(), hidiocsfeature9, uintptr(unsafe.Pointer(&data)))
	if errNo != 0 {
		file.Close()
		return nil, fmt.Errorf("syscall.Syscall(): %w", errNo)
	}

	return &Device{
		file: file,
		opts: opts,
	}, nil
}

func (d *Device) Close() error {
	return d.file.Close()
}

// ReadPacket reads one Packet from the monitor.
// An error will be returned on an I/O error or if a message could not be read or decoded.
func (d *Device) ReadPacket() (*Packet, error) {
	var pack *Packet
	for !pack.isValid() {
		buf := make([]byte, 8)
		_, err := d.file.Read(buf)
		if err != nil {
			return nil, fmt.Errorf("file.Read(): %w", err)
		}

		var data [8]byte
		copy(data[:], buf)

		// If the "magic byte" is present no decryption is necessary.
		// This is the case for AIRCO2NTROL COACH and newer AIRCO2NTROL MINIs in general.
		if data[4] != 0x0d && !d.opts.withoutDecrypt {
			data = decrypt(data, d.opts.key)
		}

		err = validate(data)
		if err != nil {
			return nil, fmt.Errorf("validation failed: %w", err)
		}

		pack = newPacket(data)
	}
	return pack, nil
}
