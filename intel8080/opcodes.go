package intel8080

import (
	"log"
)

func (cpu *CPU) nop(_ *stepInfo) uint {
	return 4
}

// HLT       01110110          -       Halt processor
func (cpu *CPU) hlt(_ *stepInfo) uint {
	// Note: Doesn't increase cpu.PC, which starts an infinite loop
	return 7
}

// MOV D,S   01DDDSSS          -       Move register to register
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

// INR D     00DDD100          ZSPA    Increment register
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

// DCR D     00DDD101          ZSPA    Decrement register
func (cpu *CPU) dcr(info *stepInfo) uint {
	var cycles uint = 5
	ddd, _ := getOpcodeDDDSSS(info.opcode)
	destPtr, memAccess := cpu.getOpcodeRegPtr(ddd)
	if memAccess {
		cycles = 10
	}
	*destPtr--

	cpu.Sign = (*destPtr & 0b1000000) != 0
	cpu.AuxCarry = (*destPtr & 0b1111) == 0
	cpu.Zero = *destPtr == 0
	cpu.Parity = getParity(*destPtr)
	return cycles
}

// LXI RP,#  00RP0001 lb hb    -       Load register pair immediate
func (cpu *CPU) lxi(info *stepInfo) uint {
	lb, hb := cpu.getOpcodeArgs(info.PC)
	rp := getOpcodeRP(info.opcode)

	switch rp {
	case 0b00:
		cpu.B = hb
		cpu.C = lb
	case 0b01:
		cpu.D = hb
		cpu.E = lb
	case 0b10:
		cpu.H = hb
		cpu.L = lb
	case 0b11:
		cpu.SP = (uint16(hb) << 8) | uint16(lb)
	}
	return 10
}

// LDA a     00111010 lb hb    -       Load A from memory
func (cpu *CPU) lda(info *stepInfo) uint {
	lb, hb := cpu.getOpcodeArgs(info.PC)
	address := (uint16(hb) << 8) | uint16(lb)
	cpu.A = cpu.Memory[address]
	return 13
}

// STA a     00110010 lb hb    -       Store A to memory
func (cpu *CPU) sta(info *stepInfo) uint {
	lb, hb := cpu.getOpcodeArgs(info.PC)
	address := (uint16(hb) << 8) | uint16(lb)
	cpu.Memory[address] = cpu.A
	return 13
}

// LHLD a    00101010 lb hb    -       Load H:L from memory
func (cpu *CPU) lhld(info *stepInfo) uint {
	lb, hb := cpu.getOpcodeArgs(info.PC)
	address := (uint16(hb) << 8) | uint16(lb)
	cpu.L = cpu.Memory[address]
	cpu.H = cpu.Memory[address+1]
	return 16
}

// SHLD a    00100010 lb hb    -       Store H:L to memory
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
	case 0b00:
		address = (uint16(cpu.B) << 8) | uint16(cpu.C)
	case 0b01:
		address = (uint16(cpu.D) << 8) | uint16(cpu.E)
	default:
		log.Fatalf("Invalid opcode: %v", info)
	}
	cpu.A = cpu.Memory[address]
	return 7
}

// STAX RP   00RP0010 *1       -       Store indirect through BC or DE
func (cpu *CPU) stax(info *stepInfo) uint {
	var address uint16
	rp := getOpcodeRP(info.opcode)
	switch rp {
	case 0b00:
		address = (uint16(cpu.B) << 8) | uint16(cpu.C)
	case 0b01:
		address = (uint16(cpu.D) << 8) | uint16(cpu.E)
	default:
		log.Fatalf("Invalid opcode: %v", info)
	}
	cpu.Memory[address] = cpu.A
	return 7
}

// XCHG      11101011          -       Exchange DE and HL content
func (cpu *CPU) xchg(_ *stepInfo) uint {
	var tempD, tempE uint8
	tempD = cpu.D
	tempE = cpu.E
	cpu.D = cpu.H
	cpu.E = cpu.L
	cpu.H = tempD
	cpu.L = tempE
	return 5
}

// ADD S     10000SSS          ZSPCA   Add register to A
func (cpu *CPU) add(info *stepInfo) uint {
	var cycles uint = 4
	_, sss := getOpcodeDDDSSS(info.opcode)
	regPtr, memAccess := cpu.getOpcodeRegPtr(sss)
	if memAccess {
		cycles = 7
	}

	result := uint16(cpu.A) - uint16(*regPtr)
	cpu.Zero = (result & 0xFF) == 0
	cpu.Sign = result&0b10000000 != 0
	cpu.Carry = result&0b100000000 != 0
	cpu.Parity = getParity(uint8(result & 0b11111111))
	cpu.AuxCarry = ((cpu.A ^ uint8(result) ^ *regPtr) & 0b00010000) > 0 // ?? TODO: verify
	cpu.A = uint8(result & 0xFF)
	return cycles
}

// ADC S     10001SSS          ZSCPA   Add register to A with carry
func (cpu *CPU) adc(info *stepInfo) uint {
	var cycles uint = 4
	_, sss := getOpcodeDDDSSS(info.opcode)
	regPtr, memAccess := cpu.getOpcodeRegPtr(sss)
	if memAccess {
		cycles = 7
	}

	carryVal := uint16(0)
	if cpu.Carry {
		carryVal = 1
	}
	result := uint16(cpu.A) + uint16(*regPtr) + carryVal
	cpu.Zero = (result & 0xFF) == 0
	cpu.Sign = result&0b10000000 != 0
	cpu.Carry = result&0b100000000 != 0
	cpu.Parity = getParity(uint8(result & 0b11111111))
	cpu.AuxCarry = ((cpu.A ^ uint8(result) ^ *regPtr) & 0b00010000) > 0 // ?? TODO: verify
	cpu.A = uint8(result & 0xFF)
	return cycles
}

// ADI #     11000110 db       ZSCPA   Add immediate to A
func (cpu *CPU) adi(info *stepInfo) uint {
	db, _ := cpu.getOpcodeArgs(info.PC)

	result := uint16(cpu.A) + uint16(db)
	cpu.Zero = (result & 0xFF) == 0
	cpu.Sign = result&0b10000000 != 0
	cpu.Carry = result&0b100000000 != 0
	cpu.Parity = getParity(uint8(result & 0b11111111))
	cpu.AuxCarry = ((cpu.A ^ uint8(result) ^ db) & 0b00010000) > 0 // ?? TODO: verify

	cpu.A = uint8(result)
	return 7
}

// ACI #     11001110 db       ZSCPA   Add immediate to A with carry
func (cpu *CPU) aci(info *stepInfo) uint {
	db, _ := cpu.getOpcodeArgs(info.PC)

	carryVal := uint16(0)
	if cpu.Carry {
		carryVal = 1
	}
	result := uint16(cpu.A) + uint16(db) + carryVal
	cpu.Zero = (result & 0xFF) == 0
	cpu.Sign = result&0b10000000 != 0
	cpu.Carry = result&0b100000000 != 0
	cpu.Parity = getParity(uint8(result & 0b11111111))
	cpu.AuxCarry = ((cpu.A ^ uint8(result) ^ db) & 0b00010000) > 0 // ?? TODO: verify

	cpu.A = uint8(result)
	return 7
}

// DAD RP    00RP1001          C       Add register pair to HL (16 bit add)
func (cpu *CPU) dad(info *stepInfo) uint {
	rp := getOpcodeRP(info.opcode)
	var currentHL, resultHL, addend uint32
	currentHL = (uint32(cpu.H) << 8) | uint32(cpu.L)
	switch rp {
	case 0b00:
		addend = (uint32(cpu.B) << 8) | uint32(cpu.C)
	case 0b01:
		addend = (uint32(cpu.D) << 8) | uint32(cpu.E)
	case 0b10:
		addend = (uint32(cpu.H) << 8) | uint32(cpu.L)
	case 0b11:
		addend = uint32(cpu.SP)
	}
	resultHL = currentHL + addend
	cpu.Carry = resultHL&0x10000 > 0
	cpu.H = uint8(currentHL >> 8)
	cpu.L = uint8(currentHL & 0xFF)
	return 10
}

// SUB S     10010SSS          ZSCPA   Subtract register from A
func (cpu *CPU) sub(info *stepInfo) uint {
	_, sss := getOpcodeDDDSSS(info.opcode)
	regPtr, memAccess := cpu.getOpcodeRegPtr(sss)

	result := uint16(cpu.A) - uint16(*regPtr)
	cpu.Zero = (result & 0xFF) == 0
	cpu.Sign = result&0b10000000 != 0
	cpu.Carry = result&0b100000000 != 0
	cpu.Parity = getParity(uint8(result & 0b11111111))
	cpu.AuxCarry = ((cpu.A ^ uint8(result) ^ *regPtr) & 0b00010000) > 0 // ?? TODO: verify
	cpu.A = uint8(result & 0xFF)
	if memAccess {
		return 7
	} else {
		return 4
	}
}

// SUI #     11010110 db       ZSCPA   Subtract immediate from A
func (cpu *CPU) sui(info *stepInfo) uint {
	db, _ := cpu.getOpcodeArgs(info.PC)

	result := uint16(cpu.A) - uint16(db)
	cpu.Zero = (result & 0xFF) == 0
	cpu.Sign = result&0b10000000 != 0
	cpu.Carry = result&0b100000000 != 0
	cpu.Parity = getParity(uint8(result & 0b11111111))
	cpu.AuxCarry = ((cpu.A ^ uint8(result) ^ db) & 0b00010000) > 0 // ?? TODO: verify

	cpu.A = uint8(result)
	return 7
}

// SBI #     11011110 db       ZSCPA   Subtract immediate from A with borrow
func (cpu *CPU) sbi(info *stepInfo) uint {
	db, _ := cpu.getOpcodeArgs(info.PC)

	carryVal := uint16(0)
	if cpu.Carry {
		carryVal = 1
	}
	result := uint16(cpu.A) - uint16(db) - carryVal
	cpu.Zero = (result & 0xFF) == 0
	cpu.Sign = result&0b10000000 != 0
	cpu.Carry = result&0b100000000 != 0
	cpu.Parity = getParity(uint8(result & 0b11111111))
	cpu.AuxCarry = ((cpu.A ^ uint8(result) ^ db) & 0b00010000) > 0 // ?? TODO: verify

	cpu.A = uint8(result)
	return 7
}

// SBB S     10011SSS          ZSCPA   Subtract register from A with borrow
func (cpu *CPU) sbb(info *stepInfo) uint {
	_, sss := getOpcodeDDDSSS(info.opcode)
	regPtr, memAccess := cpu.getOpcodeRegPtr(sss)

	carryVal := uint16(0)
	if cpu.Carry {
		carryVal = 1
	}
	result := uint16(cpu.A) - uint16(*regPtr) - carryVal
	cpu.Zero = (result & 0xFF) == 0
	cpu.Sign = result&0b10000000 != 0
	cpu.Carry = result&0b100000000 != 0
	cpu.Parity = getParity(uint8(result & 0b11111111))
	cpu.AuxCarry = ((cpu.A ^ uint8(result) ^ *regPtr) & 0b00010000) > 0 // ?? TODO: verify
	cpu.A = uint8(result & 0xFF)
	if memAccess {
		return 7
	} else {
		return 4
	}
}

// CMP S     10111SSS          ZSPCA   Compare register with A
func (cpu *CPU) cmp(info *stepInfo) uint {
	_, sss := getOpcodeDDDSSS(info.opcode)
	ptr, memoryAccess := cpu.getOpcodeRegPtr(sss)

	result := uint16(cpu.A) - uint16(*ptr)
	cpu.Zero = (result & 0xFF) == 0
	cpu.Sign = result&0b10000000 != 0
	cpu.Carry = result&0b100000000 != 0
	cpu.Parity = getParity(uint8(result & 0b11111111))
	cpu.AuxCarry = ((cpu.A ^ uint8(result) ^ *ptr) & 0b00010000) > 0 // ?? TODO: verify

	if memoryAccess {
		return 7
	} else {
		return 4
	}
}

// PUSH RP   11RP0101 *2       -       Push register pair on the stack
func (cpu *CPU) push(info *stepInfo) uint {
	rp := getOpcodeRP(info.opcode)
	var hb, lb uint8
	switch rp {
	case 0b00:
		hb, lb = cpu.B, cpu.C
	case 0b01:
		hb, lb = cpu.D, cpu.E
	case 0b10:
		hb, lb = cpu.H, cpu.L
	case 0b11:
		hb, lb = cpu.A, cpu.getProgramStatus()
	default:
		panic("PUSH bad register pair. (this shouldn't ever happen)")
	}
	cpu.SP--
	cpu.Memory[cpu.SP] = hb
	cpu.SP--
	cpu.Memory[cpu.SP] = lb

	return 11
}

// POP RP    11RP0001 *2       *2      Pop  register pair from the stack
func (cpu *CPU) pop(info *stepInfo) uint {
	rp := getOpcodeRP(info.opcode)

	lb := cpu.Memory[cpu.SP]
	cpu.SP++
	hb := cpu.Memory[cpu.SP]
	cpu.SP++

	switch rp {
	case 0b00:
		cpu.B, cpu.C = hb, lb
	case 0b01:
		cpu.D, cpu.E = hb, lb
	case 0b10:
		cpu.H, cpu.L = hb, lb
	case 0b11:
		cpu.A = hb
		cpu.setProgramStatus(lb)
	default:
		panic("pop - invalid register pair") // this should be impossible
	}
	return 11
}

///////////////////////
// CALLs
///////////////////////
func (cpu *CPU) call(info *stepInfo) uint {
	lb, hb := cpu.getOpcodeArgs(info.PC)

	nextPC := info.PC + 3
	pclo, pchi := uint8(nextPC & 0xFF), uint8(nextPC >> 8 & 0xFF)
	cpu.SP--
	cpu.Memory[cpu.SP] = pchi // TODO: make sure this byte order is correct....
	cpu.SP--
	cpu.Memory[cpu.SP] = pclo

	cpu.PC = (uint16(hb) << 8) | uint16(lb)
	return 17
}
func (cpu *CPU) cm(info *stepInfo) uint {
	if cpu.Zero {
		cpu.call(info)
		return 17
	} else {
		cpu.PC += 2
	}
	return 11
}
func (cpu *CPU) cnc(info *stepInfo) uint {
	if cpu.Carry == false {
		cpu.call(info)
		return 17
	} else {
		cpu.PC += 3
	}
	return 11
}
func (cpu *CPU) cpe(info *stepInfo) uint {
	if cpu.Parity {
		cpu.call(info)
		return 17
	} else {
		cpu.PC += 3
	}
	return 11
}
func (cpu *CPU) cpo(info *stepInfo) uint {
	if cpu.Parity == false {
		cpu.call(info)
		return 17
	} else {
		cpu.PC += 3
	}
	return 11
}
func (cpu *CPU) cp(info *stepInfo) uint {
	if cpu.Sign == false {
		cpu.call(info)
		return 17
	} else {
		cpu.PC += 3
	}
	return 11
}
func (cpu *CPU) cc(info *stepInfo) uint {
	if cpu.Carry {
		cpu.call(info)
		return 17
	} else {
		cpu.PC += 3
	}
	return 11
}
func (cpu *CPU) cz(info *stepInfo) uint {
	if cpu.Zero {
		cpu.call(info)
		return 17
	} else {
		cpu.PC += 3
	}
	return 11
}
func (cpu *CPU) cnz(info *stepInfo) uint {
	if cpu.Zero == false {
		cpu.call(info)
		return 17
	} else {
		cpu.PC += 3
	}
	return 11
}

/////////////////////
// JMPs
////////////////////
func (cpu *CPU) jmp(info *stepInfo) uint {
	// TODO: Ensure lb,hb order is correct...
	lb, hb := cpu.getOpcodeArgs(info.PC)
	cpu.PC = (uint16(hb) << 8) | uint16(lb)
	return 10
}
func (cpu *CPU) jz(info *stepInfo) uint {
	if cpu.Zero {
		lb, hb := cpu.getOpcodeArgs(info.PC)
		cpu.PC = (uint16(hb) << 8) | uint16(lb)
	} else {
		cpu.PC += 3
	}
	return 10
}
func (cpu *CPU) jnz(info *stepInfo) uint {
	if cpu.Zero == false {
		lb, hb := cpu.getOpcodeArgs(info.PC)
		cpu.PC = (uint16(hb) << 8) | uint16(lb)
	} else {
		cpu.PC += 3
	}
	return 10
}
func (cpu *CPU) jc(info *stepInfo) uint {
	if cpu.Carry {
		lb, hb := cpu.getOpcodeArgs(info.PC)
		cpu.PC = (uint16(hb) << 8) | uint16(lb)
	} else {
		cpu.PC += 3
	}
	return 10
}
func (cpu *CPU) jnc(info *stepInfo) uint {
	if cpu.Carry == false {
		lb, hb := cpu.getOpcodeArgs(info.PC)
		cpu.PC = (uint16(hb) << 8) | uint16(lb)
	} else {
		cpu.PC += 3
	}
	return 10
}
func (cpu *CPU) jm(info *stepInfo) uint {
	if cpu.Sign {
		lb, hb := cpu.getOpcodeArgs(info.PC)
		cpu.PC = (uint16(hb) << 8) | uint16(lb)
	} else {
		cpu.PC += 3
	}
	return 10
}
func (cpu *CPU) jp(info *stepInfo) uint {
	if cpu.Sign == false {
		lb, hb := cpu.getOpcodeArgs(info.PC)
		cpu.PC = (uint16(hb) << 8) | uint16(lb)
	} else {
		cpu.PC += 3
	}
	return 10
}
func (cpu *CPU) jpe(info *stepInfo) uint {
	if cpu.Parity {
		lb, hb := cpu.getOpcodeArgs(info.PC)
		cpu.PC = (uint16(hb) << 8) | uint16(lb)
	} else {
		cpu.PC += 3
	}
	return 10
}
func (cpu *CPU) jpo(info *stepInfo) uint {
	if cpu.Parity == false {
		lb, hb := cpu.getOpcodeArgs(info.PC)
		cpu.PC = (uint16(hb) << 8) | uint16(lb)
	} else {
		cpu.PC += 3
	}
	return 10
}

// STC       00110111          C       Set Carry flag
func (cpu *CPU) stc(_ *stepInfo) uint {
	cpu.Carry = true
	return 4
}

// EI        11111011          -       Enable interrupts
func (cpu *CPU) ei(_ *stepInfo) uint {
	// TODO: this action should be deferred until the next instruction has completed
	cpu.deferInterruptsEnable = true
	return 4
}

// DI        11110011          -       Disable interrupts
func (cpu *CPU) di(_ *stepInfo) uint {
	cpu.InterruptsEnabled = false
	return 4
}

// MVI D,#   00DDD110 db       -       Move immediate to register
func (cpu *CPU) mvi(info *stepInfo) uint {
	ddd, _ := getOpcodeDDDSSS(info.opcode)
	regPtr, memoryAccess := cpu.getOpcodeRegPtr(ddd)
	db, _ := cpu.getOpcodeArgs(info.PC)
	*regPtr = db
	if memoryAccess {
		return 10
	} else {
		return 7
	}
}

// CPI #     11111110          ZSPCA   Compare immediate with A
func (cpu *CPU) cpi(info *stepInfo) uint {
	// TODO: test this...
	db, _ := cpu.getOpcodeArgs(info.PC)
	result := int16(cpu.A) - int16(db)
	cpu.Carry = (result & 0b100000000) != 0
	cpu.AuxCarry = ^(int16(cpu.A)^result^int16(db))&0x10 > 0
	cpu.Zero = (result & 0xFF) == 0
	cpu.Sign = (result & 0x80) != 0
	cpu.Parity = getParity(uint8(result))
	return 0
}

// INX RP    00RP0011          -       Increment register pair
func (cpu *CPU) inx(info *stepInfo) uint {
	rp := getOpcodeRP(info.opcode)
	switch rp {
	case 0b00:
		cpu.B, cpu.C = incRegisterPair(cpu.B, cpu.C)
	case 0b01:
		cpu.D, cpu.E = incRegisterPair(cpu.D, cpu.E)
	case 0b10:
		cpu.H, cpu.L = incRegisterPair(cpu.H, cpu.L)
	case 0b11:
		cpu.SP++
	}
	return 5
}

// DCX RP    00RP1011          -       Decrement register pair
func (cpu *CPU) dcx(info *stepInfo) uint {
	rp := getOpcodeRP(info.opcode)
	switch rp {
	case 0b00:
		cpu.B, cpu.C = decRegisterPair(cpu.B, cpu.C)
	case 0b01:
		cpu.D, cpu.E = decRegisterPair(cpu.D, cpu.E)
	case 0b10:
		cpu.H, cpu.L = decRegisterPair(cpu.H, cpu.L)
	case 0b11:
		cpu.SP--
	}
	return 5
}

// ORA S     10110SSS          ZSPCA   OR  register with A
func (cpu *CPU) ora(info *stepInfo) uint {
	_, sss := getOpcodeDDDSSS(info.opcode)
	ptr, memoryAccess := cpu.getOpcodeRegPtr(sss)
	cpu.A |= *ptr

	cpu.Zero = false
	cpu.Sign = cpu.A&0b10000000 > 0
	cpu.Carry = false
	cpu.AuxCarry = false
	cpu.Parity = getParity(cpu.A)

	if memoryAccess {
		return 7
	} else {
		return 4
	}
}

// ORI #     11110110          ZSPCA   OR  immediate with A
func (cpu *CPU) ori(info *stepInfo) uint {
	db, _ := cpu.getOpcodeArgs(info.PC)
	cpu.A |= db

	cpu.Zero = false
	cpu.Sign = cpu.A&0b10000000 > 0
	cpu.Carry = false
	cpu.AuxCarry = false
	cpu.Parity = getParity(cpu.A)
	return 7
}

///////////////////
// RETs
///////////////////
func (cpu *CPU) ret(_ *stepInfo) uint {
	hb, lb := cpu.Memory[cpu.SP+1], cpu.Memory[cpu.SP]
	cpu.SP += 2
	cpu.PC = (uint16(hb) << 8) | uint16(lb)
	return 10
}
func (cpu *CPU) rc(info *stepInfo) uint {
	if cpu.Carry {
		cpu.ret(info)
		return 11
	}
	return 5
}
func (cpu *CPU) rm(info *stepInfo) uint {
	if cpu.Sign {
		cpu.ret(info)
		return 11
	} else {
		cpu.PC += 1
	}
	return 5
}
func (cpu *CPU) rnc(info *stepInfo) uint {
	if cpu.Carry == false {
		cpu.ret(info)
		return 11
	} else {
		cpu.PC += 1
	}
	return 5
}
func (cpu *CPU) rnz(info *stepInfo) uint {
	if cpu.Zero == false {
		cpu.ret(info)
		return 11
	} else {
		cpu.PC += 1
	}
	return 5
}
func (cpu *CPU) rz(info *stepInfo) uint {
	if cpu.Zero {
		cpu.ret(info)
		return 11
	} else {
		cpu.PC += 1
	}
	return 5
}
func (cpu *CPU) rp(info *stepInfo) uint {
	if cpu.Sign == false {
		cpu.ret(info)
		return 11
	} else {
		cpu.PC += 1
	}
	return 5
}
func (cpu *CPU) rpe(info *stepInfo) uint {
	if cpu.Parity == true {
		cpu.ret(info)
		return 11
	} else {
		cpu.PC += 1
	}
	return 50
}
func (cpu *CPU) rpo(info *stepInfo) uint {
	if cpu.Parity == false {
		cpu.ret(info)
		return 11
	} else {
		cpu.PC += 1
	}
	return 5
}

//////////////////
// Rotates
//////////////////
// RLC       00000111          C       Rotate A left
func (cpu *CPU) rlc(_ *stepInfo) uint {
	highBit := cpu.A >> 7
	cpu.Carry = highBit == 1
	cpu.A = cpu.A << 1
	cpu.A = cpu.A | highBit
	return 4
}

// RAL       00010111          C       Rotate A left through carry
func (cpu *CPU) ral(_ *stepInfo) uint {
	var oldCarry uint8
	if cpu.Carry {
		oldCarry = 1
	}
	highBit := cpu.A >> 7
	cpu.Carry = highBit == 1
	cpu.A = cpu.A << 1
	cpu.A = cpu.A | oldCarry
	return 4
}

// RRC       00001111          C       Rotate A right
func (cpu *CPU) rrc(_ *stepInfo) uint {
	lowBit := cpu.A & 0b1
	cpu.Carry = lowBit == 1
	cpu.A = cpu.A >> 1
	cpu.A = cpu.A | (lowBit << 7)
	return 4
}

// RAR       00011111          C       Rotate A right through carry
func (cpu *CPU) rar(_ *stepInfo) uint {
	var oldCarry uint8
	if cpu.Carry {
		oldCarry = 1
	}
	lowBit := cpu.A & 0b1
	cpu.Carry = lowBit == 1
	cpu.A = cpu.A >> 1
	cpu.A = cpu.A | (oldCarry << 7)
	return 4
}

// CMA       00101111          -       Complement A
func (cpu *CPU) cma(_ *stepInfo) uint {
	cpu.A = ^cpu.A
	return 4
}

// CMC       00111111          C       Complement Carry flag
func (cpu *CPU) cmc(_ *stepInfo) uint {
	if cpu.Carry {
		cpu.Carry = false
	} else {
		cpu.Carry = true
	}
	return 4
}

// ANA S     10100SSS          ZSCPA   AND register with A
func (cpu *CPU) ana(info *stepInfo) uint {
	_, sss := getOpcodeDDDSSS(info.opcode)
	ptr, memoryAccess := cpu.getOpcodeRegPtr(sss)

	cpu.A &= *ptr
	cpu.Carry = false

	if memoryAccess {
		return 7
	} else {
		return 4
	}
}

// XRA S     10101SSS          ZSPCA   ExclusiveOR register with A
func (cpu *CPU) xra(info *stepInfo) uint {
	_, sss := getOpcodeDDDSSS(info.opcode)
	ptr, memoryAccess := cpu.getOpcodeRegPtr(sss)

	cpu.A ^= *ptr

	cpu.Zero = cpu.A != 0
	cpu.Sign = cpu.A&0b10000000 > 0
	cpu.Carry = false
	cpu.AuxCarry = false
	cpu.Parity = getParity(cpu.A)
	if memoryAccess {
		return 7
	} else {
		return 4
	}
}

// DAA       00100111          ZSPCA   Decimal Adjust accumulator
func (cpu *CPU) daa(_ *stepInfo) uint {

	var addend uint8
	msb := cpu.A >> 4
	lsb := cpu.A & 0x0F

	if lsb > 9 || cpu.AuxCarry == true {
		addend += 0x06
	}
	if msb > 9 || cpu.Carry == true {
		addend += 0x60
	}

	result := uint16(cpu.A) + uint16(addend)
	cpu.Zero = (result & 0xFF) == 0
	cpu.Sign = result&0b10000000 != 0
	cpu.Carry = result&0b100000000 != 0
	cpu.Parity = getParity(uint8(result & 0b11111111))
	cpu.AuxCarry = ((cpu.A ^ uint8(result) ^ addend) & 0b00010000) > 0 // ?? TODO: verify

	// TODO: above could definitely be wrong for BCD formatting
	cpu.A = uint8(result)
	return 4
}

// XTHL      11100011          -       Swap H:L with top word on stack
func (cpu *CPU) xthl(_ *stepInfo) uint {
	stackLo := cpu.Memory[cpu.SP]
	stackHi := cpu.Memory[cpu.SP+1]
	cpu.Memory[cpu.SP] = cpu.L
	cpu.Memory[cpu.SP+1] = cpu.H
	cpu.L = stackLo
	cpu.H = stackHi
	return 18
}

// ANI #     11100110 db       ZSPCA   AND immediate with A
func (cpu *CPU) ani(info *stepInfo) uint {
	db, _ := cpu.getOpcodeArgs(info.PC)
	cpu.A &= db

	cpu.Zero = false
	cpu.Sign = cpu.A&0b10000000 > 0
	cpu.Carry = false
	cpu.AuxCarry = false
	cpu.Parity = getParity(cpu.A)
	return 7
}

// XRI #     11101110 db       ZSPCA   ExclusiveOR immediate with A
func (cpu *CPU) xri(info *stepInfo) uint {
	db, _ := cpu.getOpcodeArgs(info.PC)
	cpu.A ^= db

	cpu.Zero = false
	cpu.Sign = cpu.A&0b10000000 > 0
	cpu.Carry = false
	cpu.AuxCarry = false
	cpu.Parity = getParity(cpu.A)
	return 7
}

// SPHL      11111001          -       Set SP to content of H:L
func (cpu *CPU) sphl(_ *stepInfo) uint {
	cpu.SP = (uint16(cpu.H) << 8) | uint16(cpu.L)
	return 0
}

// PCHL      11101001          -       Jump to address in H:L
func (cpu *CPU) pchl(_ *stepInfo) uint {
	cpu.PC = (uint16(cpu.H) << 8) | uint16(cpu.L)
	return 0
}

func (cpu *CPU) Interrupt(interruptType uint) {
	if cpu.InterruptsEnabled {
		// Push PC on to the stack
		pcHi, pcLo := uint8(cpu.PC >> 8), uint8(cpu.PC & 0xFF)
		cpu.SP--
		cpu.Memory[cpu.SP] = pcHi
		cpu.SP--
		cpu.Memory[cpu.SP] = pcLo

		// Move PC to type * 8
		cpu.PC = 8 * uint16(interruptType)
		cpu.InterruptsEnabled = false
	}
}