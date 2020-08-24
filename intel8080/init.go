package intel8080

import (
	"fmt"
	"io/ioutil"
)

type CPU struct {
	DEBUG bool

	Memory [65536]byte
	// Registers
	A uint8
	B uint8
	C uint8
	D uint8
	E uint8
	H uint8
	L uint8

	PC uint16
	SP uint16

	// Flags
	Sign, Zero, Parity, Carry, AuxCarry bool
	InterruptsEnabled                   bool

	// op table
	table []func(*stepInfo) uint
}

type stepInfo struct {
	PC     uint16
	opcode uint8
}

var instructionNames = []string{
	"nop", "lxi b,#", "stax b", "inx b", "inr b", "dcr b", "mvi b,#", "rlc",
	"ill", "dad b", "ldax b", "dcx b", "inr c", "dcr c", "mvi c,#", "rrc",
	"ill", "lxi d,#", "stax d", "inx d", "inr d", "dcr d", "mvi d,#", "ral",
	"ill", "dad d", "ldax d", "dcx d", "inr e", "dcr e", "mvi e,#", "rar",
	"ill", "lxi h,#", "shld", "inx h", "inr h", "dcr h", "mvi h,#", "daa",
	"ill", "dad h", "lhld", "dcx h", "inr l", "dcr l", "mvi l,#", "cma",
	"ill", "lxi sp,#", "sta $", "inx sp", "inr M", "dcr M", "mvi M,#", "stc",
	"ill", "dad sp", "lda $", "dcx sp", "inr a", "dcr a", "mvi a,#", "cmc",
	"mov b,b", "mov b,c", "mov b,d", "mov b,e", "mov b,h", "mov b,l",
	"mov b,M", "mov b,a", "mov c,b", "mov c,c", "mov c,d", "mov c,e",
	"mov c,h", "mov c,l", "mov c,M", "mov c,a", "mov d,b", "mov d,c",
	"mov d,d", "mov d,e", "mov d,h", "mov d,l", "mov d,M", "mov d,a",
	"mov e,b", "mov e,c", "mov e,d", "mov e,e", "mov e,h", "mov e,l",
	"mov e,M", "mov e,a", "mov h,b", "mov h,c", "mov h,d", "mov h,e",
	"mov h,h", "mov h,l", "mov h,M", "mov h,a", "mov l,b", "mov l,c",
	"mov l,d", "mov l,e", "mov l,h", "mov l,l", "mov l,M", "mov l,a",
	"mov M,b", "mov M,c", "mov M,d", "mov M,e", "mov M,h", "mov M,l", "hlt",
	"mov M,a", "mov a,b", "mov a,c", "mov a,d", "mov a,e", "mov a,h",
	"mov a,l", "mov a,M", "mov a,a", "add b", "add c", "add d", "add e",
	"add h", "add l", "add M", "add a", "adc b", "adc c", "adc d", "adc e",
	"adc h", "adc l", "adc M", "adc a", "sub b", "sub c", "sub d", "sub e",
	"sub h", "sub l", "sub M", "sub a", "sbb b", "sbb c", "sbb d", "sbb e",
	"sbb h", "sbb l", "sbb M", "sbb a", "ana b", "ana c", "ana d", "ana e",
	"ana h", "ana l", "ana M", "ana a", "xra b", "xra c", "xra d", "xra e",
	"xra h", "xra l", "xra M", "xra a", "ora b", "ora c", "ora d", "ora e",
	"ora h", "ora l", "ora M", "ora a", "cmp b", "cmp c", "cmp d", "cmp e",
	"cmp h", "cmp l", "cmp M", "cmp a", "rnz", "pop b", "jnz $", "jmp $",
	"cnz $", "push b", "adi #", "rst 0", "rz", "ret", "jz $", "ill", "cz $",
	"call $", "aci #", "rst 1", "rnc", "pop d", "jnc $", "out p", "cnc $",
	"push d", "sui #", "rst 2", "rc", "ill", "jc $", "in p", "cc $", "ill",
	"sbi #", "rst 3", "rpo", "pop h", "jpo $", "xthl", "cpo $", "push h",
	"ani #", "rst 4", "rpe", "pchl", "jpe $", "xchg", "cpe $", "ill", "xri #",
	"rst 5", "rp", "pop psw", "jp $", "di", "cp $", "push psw", "ori #",
	"rst 6", "rm", "sphl", "jm $", "ei", "cm $", "ill", "cpi #", "rst 7",
}

var instructionBytes = []uint8{
	1, 3, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1,
	1, 3, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1,
	1, 3, 3, 1, 1, 1, 2, 1, 1, 1, 3, 1, 1, 1, 2, 1,
	1, 3, 3, 1, 1, 1, 2, 1, 1, 1, 3, 1, 1, 1, 2, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 3, 3, 3, 1, 2, 1, 1, 1, 3, 3, 3, 3, 2, 1,
	1, 1, 3, 2, 3, 1, 2, 1, 1, 1, 3, 2, 3, 3, 2, 1,
	1, 1, 3, 1, 3, 1, 2, 1, 1, 1, 3, 1, 3, 3, 2, 1,
	1, 1, 3, 1, 3, 1, 2, 1, 1, 1, 3, 1, 3, 3, 2, 1,
}

func (cpu *CPU) createInstructionTable() {
	cpu.table = []func(*stepInfo) uint{
		cpu.nop, cpu.lxi, cpu.stax, cpu.inx, cpu.inr, cpu.dcr, cpu.mvi, cpu.rlc, nil, cpu.dad, cpu.ldax, cpu.dcx, cpu.inr, cpu.dcr, cpu.mvi, cpu.rrc,
		nil, cpu.lxi, cpu.stax, cpu.inx, cpu.inr, cpu.dcr, cpu.mvi, cpu.ral, nil, cpu.dad, cpu.ldax, cpu.dcx, cpu.inr, cpu.dcr, cpu.mvi, cpu.rar,
		nil, cpu.lxi, cpu.shld, cpu.inx, cpu.inr, cpu.dcr, cpu.mvi, cpu.daa, nil, cpu.dad, cpu.lhld, cpu.dcx, cpu.inr, cpu.dcr, cpu.mvi, cpu.cma,
		nil, cpu.lxi, cpu.sta, cpu.inx, cpu.inr, cpu.dcr, cpu.mvi, cpu.stc, nil, cpu.dad, cpu.lda, cpu.dcx, cpu.inr, cpu.dcr, cpu.mvi, cpu.cmc,

		cpu.mov, cpu.mov, cpu.mov, cpu.mov, cpu.mov, cpu.mov, cpu.mov, cpu.mov, cpu.mov, cpu.mov, cpu.mov, cpu.mov, cpu.mov, cpu.mov, cpu.mov, cpu.mov,
		cpu.mov, cpu.mov, cpu.mov, cpu.mov, cpu.mov, cpu.mov, cpu.mov, cpu.mov, cpu.mov, cpu.mov, cpu.mov, cpu.mov, cpu.mov, cpu.mov, cpu.mov, cpu.mov,
		cpu.mov, cpu.mov, cpu.mov, cpu.mov, cpu.mov, cpu.mov, cpu.mov, cpu.mov, cpu.mov, cpu.mov, cpu.mov, cpu.mov, cpu.mov, cpu.mov, cpu.mov, cpu.mov,
		cpu.mov, cpu.mov, cpu.mov, cpu.mov, cpu.mov, cpu.mov, cpu.hlt, cpu.mov, cpu.mov, cpu.mov, cpu.mov, cpu.mov, cpu.mov, cpu.mov, cpu.mov, cpu.mov,

		cpu.add, cpu.add, cpu.add, cpu.add, cpu.add, cpu.add, cpu.add, cpu.add, cpu.adc, cpu.adc, cpu.adc, cpu.adc, cpu.adc, cpu.adc, cpu.adc, cpu.adc,
		cpu.sub, cpu.sub, cpu.sub, cpu.sub, cpu.sub, cpu.sub, cpu.sub, cpu.sub, cpu.sbb, cpu.sbb, cpu.sbb, cpu.sbb, cpu.sbb, cpu.sbb, cpu.sbb, cpu.sbb,
		cpu.ana, cpu.ana, cpu.ana, cpu.ana, cpu.ana, cpu.ana, cpu.ana, cpu.ana, cpu.xra, cpu.xra, cpu.xra, cpu.xra, cpu.xra, cpu.xra, cpu.xra, cpu.xra,
		cpu.ora, cpu.ora, cpu.ora, cpu.ora, cpu.ora, cpu.ora, cpu.ora, cpu.ora, cpu.cmp, cpu.cmp, cpu.cmp, cpu.cmp, cpu.cmp, cpu.cmp, cpu.cmp, cpu.cmp,

		cpu.rnz, cpu.pop, cpu.jnz, cpu.jmp, cpu.cnz, cpu.push, cpu.adi, cpu.rst, cpu.rz, cpu.ret, cpu.jz, nil, cpu.cz, cpu.call, cpu.aci, cpu.rst,
		cpu.rnc, cpu.pop, cpu.jnc, cpu.out, cpu.cnc, cpu.push, cpu.sui, cpu.rst, cpu.rc, nil, cpu.jc, cpu.in, cpu.cc, nil, cpu.sbi, cpu.rst,
		cpu.rpo, cpu.pop, cpu.jpo, cpu.xthl, cpu.cpo, cpu.push, cpu.ani, cpu.rst, cpu.rpe, cpu.pchl, cpu.jpe, cpu.xchg, cpu.cpe, nil, cpu.xri, cpu.rst,
		cpu.rp, cpu.pop, cpu.jp, cpu.di, cpu.cp, cpu.push, cpu.ori, cpu.rst, cpu.rm, cpu.sphl, cpu.jm, cpu.ei, cpu.cm, nil, cpu.cpi, cpu.rst,
	}
}

func (cpu *CPU) Reset() {
	// TODO
}

func NewCPU() *CPU {
	cpu := CPU{}
	cpu.createInstructionTable()
	cpu.Reset()
	return &cpu
}

func (cpu *CPU) LoadInvaders(filepath string) error {
	files := []string{
		"invaders.h",
		"invaders.g",
		"invaders.f",
		"invaders.e",
	}
	offset := 0
	for _, filename := range files {
		romPath := filepath + filename
		fmt.Println("loading %s", romPath)
		data, err := ioutil.ReadFile(romPath)
		if err != nil {
			return fmt.Errorf("loadRom: failed reading file: %v", err)
		}

		for _, b := range data {
			cpu.Memory[offset] = b
			offset++
		}
	}
	return nil
}
