package intel8080

import "log"

func (cpu *CPU) nop(info *stepInfo) uint {
	return 4
}

func (cpu *CPU) hlt(info *stepInfo) uint {
	// TODO
	return 7
}

func (cpu *CPU) mov(info *stepInfo) uint {
	var cycles uint = 5
	ddd, sss := getOpcodeDDDSSS(info.opcode)

	var destPtr, srcPtr *uint8
	destPtr, memAccess := cpu.getOpcodeRegPtr(ddd)
	if memAccess {
		cycles = 7
	}
	srcPtr, memAccess = cpu.getOpcodeRegPtr(sss)
	if memAccess {
		cycles = 7
	}

	*destPtr = *srcPtr // Maybe check this for nil-dereference? (it should be impossible, though)
	return cycles
}

func (cpu *CPU) jmp(info *stepInfo) uint {
	// TODO: Ensure lb,hb order is correct...
	lb, hb := cpu.getOpcodeArgs(info.PC)
	cpu.PC = (uint16(hb) << 8) | uint16(lb)
	return 10
}

func (cpu *CPU) adc(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) inr(info *stepInfo) uint {
	var cycles uint = 5
	ddd, _ := getOpcodeDDDSSS(info.opcode)
	destPtr, memAccess := cpu.getOpcodeRegPtr(ddd)
	if memAccess {
		cycles = 10
	}
	*destPtr++

	cpu.Sign = (*destPtr & 0b1000000) != 0
	cpu.AuxCarry = (*destPtr & 0b1111) == 0
	cpu.Zero = *destPtr == 0
	cpu.Parity = getParity(*destPtr)
	return cycles
}

func (cpu *CPU) lxi(info *stepInfo) uint {
	lb, hb := cpu.getOpcodeArgs(info.PC)
	rp := getOpcodeRP(info.opcode)

	switch rp {
	case 0x00:
		cpu.B = hb
		cpu.C = lb
	case 0x01:
		cpu.D = hb
		cpu.E = lb
	case 0x10:
		cpu.H = hb
		cpu.L = lb
	case 0x11:
		cpu.SP = (uint16(hb) << 8) | uint16(lb)
	}
	return 10
}

func (cpu *CPU) lda(info *stepInfo) uint {
	lb, hb := cpu.getOpcodeArgs(info.PC)
	address := (uint16(hb) << 8) | uint16(lb)
	cpu.A = cpu.Memory[address]
	return 13
}

func (cpu *CPU) sta(info *stepInfo) uint {
	lb, hb := cpu.getOpcodeArgs(info.PC)
	address := (uint16(hb) << 8) | uint16(lb)
	cpu.Memory[address] = cpu.A
	return 13
}

func (cpu *CPU) lhld(info *stepInfo) uint {
	lb, hb := cpu.getOpcodeArgs(info.PC)
	address := (uint16(hb) << 8) | uint16(lb)
	cpu.L = cpu.Memory[address]
	cpu.H = cpu.Memory[address+1]
	return 16
}

func (cpu *CPU) shld(info *stepInfo) uint {
	lb, hb := cpu.getOpcodeArgs(info.PC)
	address := (uint16(hb) << 8) | uint16(lb)
	cpu.Memory[address] = cpu.L
	cpu.Memory[address+1] = cpu.H
	return 16
}

func (cpu *CPU) ldax(info *stepInfo) uint {
	var address uint16
	rp := getOpcodeRP(info.opcode)
	switch rp {
	case 0x00:
		address = (uint16(cpu.B) << 8) | uint16(cpu.C)
	case 0x01:
		address = (uint16(cpu.D) << 8) | uint16(cpu.E)
	default:
		log.Fatalf("Invalid opcode: %v", info)
	}
	cpu.A = cpu.Memory[address]
	return 7
}

func (cpu *CPU) stax(info *stepInfo) uint {
	var address uint16
	rp := getOpcodeRP(info.opcode)
	switch rp {
	case 0x00:
		address = (uint16(cpu.B) << 8) | uint16(cpu.C)
	case 0x01:
		address = (uint16(cpu.D) << 8) | uint16(cpu.E)
	default:
		log.Fatalf("Invalid opcode: %v", info)
	}
	cpu.Memory[address] = cpu.A
	return 7
}

func (cpu *CPU) xchg(info *stepInfo) uint {
	var tempD, tempE uint8
	tempD = cpu.D
	tempE = cpu.E
	cpu.D = cpu.H
	cpu.E = cpu.L
	cpu.H = tempD
	cpu.L = tempE
	return 5
}

func (cpu *CPU) add(info *stepInfo) uint {
	var cycles uint = 4
	_, sss := getOpcodeDDDSSS(info.opcode)
	regPtr, memAccess := cpu.getOpcodeRegPtr(sss)
	if memAccess {
		cycles = 7
	}
	result := cpu.A + *regPtr
	cpu.Zero = result == 0
	cpu.Sign = result&0b10000000 != 0
	cpu.Parity = getParity(result)
	cpu.AuxCarry = ((cpu.A ^ result ^ *regPtr) & 0b00010000) > 0 // ?? TODO: test this

	cpu.A = result
	return cycles
}

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