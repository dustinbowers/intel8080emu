package intel8080

import "fmt"

func (cpu *CPU) setFlagSZP(result uint8) {
	cpu.Zero = result == 0
	cpu.Sign = result>>7 > 0
	cpu.Parity = getParity(result)
}

func (cpu *CPU) getOpcodeRegPtr(regIndicator uint8) *uint8 {
	switch regIndicator {
	case 0b111:
		return &cpu.A
	case 0b000:
		return &cpu.B
	case 0b001:
		return &cpu.C
	case 0b010:
		return &cpu.D
	case 0b011:
		return &cpu.E
	case 0b100:
		return &cpu.H
	case 0b101:
		return &cpu.L
	default:
		panic(fmt.Sprintf("Bad register pair indicator 0b%03b\n", regIndicator))
	}
}

func (cpu *CPU) getOpcodeArgs(PC uint16) (uint8, uint8) {
	return cpu.memory.Read(PC + 1), cpu.memory.Read(PC + 2)
}

func (cpu *CPU) setProgramStatus(psw uint8) {
	cpu.Sign = (psw >> 7 & 0b1) > 0
	cpu.Zero = (psw >> 6 & 0b1) > 0
	cpu.AuxCarry = (psw >> 4 & 0b1) > 0
	cpu.Parity = (psw >> 2 & 0b1) > 0
	cpu.Carry = (psw >> 0 & 0b1) > 0
}

func (cpu *CPU) getProgramStatus() uint8 {
	// S, Z, 0, AC, 0, P, 1, CY
	status := uint8(0b00000010)
	if cpu.Sign {
		status |= 1 << 7
	}
	if cpu.Zero {
		status |= 1 << 6
	}
	if cpu.AuxCarry {
		status |= 1 << 4
	}
	if cpu.Parity {
		status |= 1 << 2
	}
	if cpu.Carry {
		status |= 1 << 0
	}
	return status
}

func getParity(b uint8) bool {
	ones := uint8(0)
	// TODO: this could be optimized...
	for i := 0; i < 8; i++ {
		ones += (b >> i) & 0b1
	}
	return (ones & 0b1) == 0
}

// getOpcodeRP returns the register pair indicator from an opcode
func getOpcodeRP(opcode uint8) uint8 {
	// e.g. opcode: 0b00RP0011
	return (opcode >> 4) & 0b11
}

// getOpcodeDDDSSS returns the source and destination indicators from an opcode
func getOpcodeDDDSSS(opcode uint8) (ddd uint8, sss uint8) {
	// e.g. opcode: 0b01DDDSSS
	ddd = (opcode >> 3) & 0b111
	sss = opcode & 0b111
	return ddd, sss
}

func incRegisterPair(hi uint8, lo uint8) (uint8, uint8) {
	result := (uint16(hi) << 8) | uint16(lo)
	result++
	return uint8(result >> 8), uint8(result & 0xFF)
}

func decRegisterPair(hi uint8, lo uint8) (uint8, uint8) {
	result := (uint16(hi) << 8) | uint16(lo)
	result--
	return uint8(result >> 8), uint8(result & 0xFF)
}
