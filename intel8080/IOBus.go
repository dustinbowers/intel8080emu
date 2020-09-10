package intel8080

import (
	"fmt"
)

type IOBus struct {
	DEBUG   bool
	bitMask byte
	shiftH  byte
	shiftL  byte
	offset  byte
	input1  byte
	input2  byte
}

func NewIOBus() *IOBus {
	bus := IOBus{}
	bus.bitMask = 0b111
	bus.input1 = 0x00
	bus.input2 = 0x00
	return &bus
}

func (bus *IOBus) Read(b byte) byte {
	switch b {
	case 0x01:
		//Read 1
		//BIT	0	coin (0 when active)
		//1	P2 start button
		//2	P1 start button
		//3	?
		//4	P1 shoot button
		//5	P1 joystick left
		//6	P1 joystick right
		//7	?
		if bus.DEBUG {
			fmt.Printf("IOBus.Read(0x%02x) = 0b%08b - read input\n", b, bus.input1)
		}
		return bus.input1
	case 0x02:
		//Read 2
		//BIT	0,1	dipswitch number of lives (0:3,1:4,2:5,3:6)
		//2	tilt 'button'
		//3	dipswitch bonus life at 1:1000,0:1500
		//4	P2 shoot button
		//5	P2 joystick left
		//6	P2 joystick right
		//7	dipswitch coin info 1:off,0:on
		if bus.DEBUG {
			fmt.Printf("IOBus.Read(0x%02x) = 0b%08b - read input\n", b, bus.input2)
		}
		return bus.input2
	case 0x03:
		shift := uint16(bus.shiftH)<<8 | uint16(bus.shiftL)
		result := byte(shift >> (8 - bus.offset))
		if bus.DEBUG {
			fmt.Printf("IOBus.Read(0x%02x) = 0b%08b - read shift register\n", b, result)
		}
		return result
	default:
		if bus.DEBUG {
			fmt.Printf("IOBus.Read(0x%02x) ---------\n", b)
		}
	}
	return 0
}

func (bus *IOBus) Write(b byte, A byte) {
	switch b {
	case 0x02:
		bus.offset = A & bus.bitMask
		if bus.DEBUG {
			fmt.Printf("IOBus.Write(0x%02x, 0b%08b) - offset = 0b%08b\n", b, A, bus.offset)
		}
	case 0x04:
		bus.shiftL = bus.shiftH
		bus.shiftH = A
		if bus.DEBUG {
			fmt.Printf("IOBus.Write(0x%02x, 0b%08b) - H: 0b%08b L: 0b%08b\n", b, A, bus.shiftH, bus.shiftL)
		}
	default:
		//fmt.Printf("IOBus.Write(0x%02x, 0b%08b) ---------\n", b, A)
	}
}

func (bus *IOBus) HandleInput(portNumber uint8, bitNumber uint8, pressed bool) {
	if portNumber == 1 {
		if pressed {
			bus.input1 |= 1 << bitNumber
		} else {
			bus.input1 &= ^(1 << bitNumber)
		}
	} else if portNumber == 2 {
		if pressed {
			if pressed {
				bus.input2 |= 1 << bitNumber
			} else {
				bus.input2 &= ^(1 << bitNumber)
			}
		}
	}
}
