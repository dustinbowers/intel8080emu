package intel8080


func (cpu *CPU) nop(info *stepInfo) uint {
	return 4
}

func (cpu *CPU) hlt(info *stepInfo) uint {
	// TODO
	return 7
}

func (cpu *CPU) mov(info *stepInfo) uint {
	ddd, sss := opcodeDDDSSS(info.opcode)
	var cycles uint = 5

	var destPtr, srcPtr *uint8
	switch ddd {
	case 0b111: destPtr = &cpu.A
	case 0b000: destPtr = &cpu.B
	case 0b001: destPtr = &cpu.C
	case 0b010: destPtr = &cpu.D
	case 0b011: destPtr = &cpu.E
	case 0b100: destPtr = &cpu.H
	case 0b101: destPtr = &cpu.L
	case 0b110:
		cycles = 7
		memOffset := (uint16(cpu.H) << 8) | uint16(cpu.L)
		destPtr = &cpu.Memory[memOffset]
	}

	switch sss {
	case 0b111: srcPtr = &cpu.A
	case 0b000: srcPtr = &cpu.B
	case 0b001: srcPtr = &cpu.C
	case 0b010: srcPtr = &cpu.D
	case 0b011: srcPtr = &cpu.E
	case 0b100: srcPtr = &cpu.H
	case 0b101: srcPtr = &cpu.L
	case 0b110:
		cycles = 7
		memOffset := (uint16(cpu.H) << 8) | uint16(cpu.L)
		srcPtr = &cpu.Memory[memOffset]
	}
	*destPtr = *srcPtr // Maybe check this for nil-dereference? (it shouldn't be possible though)
	return cycles
}

func (cpu *CPU) jmp(info *stepInfo) uint {
	// TODO: Ensure lb,hb order is correct...
	lb, hb := cpu.getOpcodeArgs(info.PC)
	cpu.PC = (uint16(hb) << 8) | uint16(lb)
	return 10
}

func (cpu *CPU) getOpcodeArgs(PC uint16) (byte1, byte2 uint8) {
	return cpu.Memory[PC+1], cpu.Memory[PC+2]
}

func opcodeDDDSSS(opcode uint8) (ddd uint8, sss uint8) {
	ddd = (opcode >> 3) & 0b111
	sss = opcode & 0b111
	return ddd, sss
}

