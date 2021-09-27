package co2mon

// Packet represents single operation/value pair from the monitor
type Packet struct {
	Ops   ops
	Value int
}

func newPacket(data [8]byte) *Packet {
	value := int(data[1])<<8 | int(data[2])
	return &Packet{
		Ops:   newOperation(data[0]),
		Value: value,
	}
}

func (p *Packet) isValid() bool {
	return p != nil && p.Ops != invalidOps
}

type ops uint8

const (
	invalidOps ops = iota
	// HumOps indicates relative humidity in units of 0.01%.
	HumOps
	// TempOps indicates temperature in Kelvin (unit of 1/16th K).
	TempOps
	// Co2Ops indicates CO₂ concentration in ppm.
	Co2Ops
)

func newOperation(b byte) ops {
	switch rune(b) {
	case 'A':
		return HumOps
	case 'B':
		return TempOps
	case 'P':
		return Co2Ops
	default:
		return invalidOps
	}
}
