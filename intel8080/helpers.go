package intel8080

// ----------------------------
// -------- Helpers -----------

func (cpu *CPU) getOpcodeRegPtr(regIndicator uint8) (*uint8, bool) {
	var ptr *uint8
	memoryAccess := false
	switch regIndicator {
	case 0b111:
		ptr = &cpu.A
	case 0b000:
		ptr = &cpu.B
	case 0b001:
		ptr = &cpu.C
	case 0b010:
		ptr = &cpu.D
	case 0b011:
		ptr = &cpu.E
	case 0b100:
		ptr = &cpu.H
	case 0b101:
		ptr = &cpu.L
	case 0b110:
		memoryAccess = true
		memOffset := (uint16(cpu.H) << 8) | uint16(cpu.L)
		ptr = &cpu.Memory[memOffset]
	}
	return ptr, memoryAccess
}

func (cpu *CPU) getOpcodeArgs(PC uint16) (byte1, byte2 uint8) {
	return cpu.Memory[PC+1], cpu.Memory[PC+2]
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
		ones += (b >> 0) & 0b1
	}
	return (ones & 0b1) == 0
}

func getOpcodeRP(opcode uint8) uint8 {
	return (opcode >> 4) & 0b11
}

func getOpcodeDDDSSS(opcode uint8) (ddd uint8, sss uint8) {
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
