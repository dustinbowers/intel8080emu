package intel8080

import (
	"log"
)

type IOBus struct {
	DEBUG   bool
	bitMask byte
	shiftH  byte
	shiftL  byte
	offset  byte
	input   byte
}

func NewIOBus() *IOBus {
	bus := IOBus{}
	bus.bitMask = 0b111
	bus.input = 0xFF
	return &bus
}

func (bus *IOBus) Read(b byte) byte {
	switch b {
	case 0x01:
		if bus.DEBUG {
			log.Printf("IOBus.Read(0x%02x) = 0b%08b - read input\n", b, bus.input)
		}
		return bus.input
	case 0x02:
		//Read 2
		//BIT	0,1	dipswitch number of lives (0:3,1:4,2:5,3:6)
		//2	tilt 'button'
		//3	dipswitch bonus life at 1:1000,0:1500
		//4	P2 shoot button
		//5	P2 joystick left
		//6	P2 joystick right
		//7	dipswitch coin info 1:off,0:on
		return 0b00000000
	case 0x03:
		shift := uint16(bus.shiftH)<<8 | uint16(bus.shiftL)
		result := byte(shift >> (8 - bus.offset))
		if bus.DEBUG {
			log.Printf("IOBus.Read(0x%02x) = 0b%08b - read shift register\n", b, result)
		}
		return result
	default:
		if bus.DEBUG {
			log.Printf("IOBus.Read(0x%02x) ---------\n", b)
		}
	}
	return 0
}

func (bus *IOBus) Write(b byte, A byte) {
	switch b {
	case 0x02:
		bus.offset = A & bus.bitMask
		if bus.DEBUG {
			log.Printf("IOBus.Write(0x%02x, 0b%08b) - offset = 0b%08b\n", b, A, bus.offset)
		}
	case 0x04:
		bus.shiftL = bus.shiftH
		bus.shiftH = A
		if bus.DEBUG {
			log.Printf("IOBus.Write(0x%02x, 0b%08b) - H: 0b%08b L: 0b%08b\n", b, A, bus.shiftH, bus.shiftL)
		}
	default:
		//fmt.Printf("IOBus.Write(0x%02x, 0b%08b) ---------\n", b, A)
	}
}
