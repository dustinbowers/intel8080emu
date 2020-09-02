package intel8080

type IOBus struct {
	bitMask byte
	shiftH byte
	shiftL byte
	offset byte
	input byte
}

func NewIOBus() *IOBus {
	bus := IOBus{}
	return &bus
}

func (bus *IOBus) Read(b byte) byte {
	switch b {
	case 0x01:
		return bus.input
	case 0x03:
		shift := uint16(bus.shiftH) << 8 | uint16(bus.shiftL)
		return byte(shift >> (8 - bus.offset))
	}
	return 0
}

func (bus *IOBus) Write(b byte, A byte) {
	switch b {
	case 0x02:
		bus.offset = byte(A & bus.bitMask)
	case 0x04:
		bus.shiftL = bus.shiftH
		bus.shiftH = A
	}
}
