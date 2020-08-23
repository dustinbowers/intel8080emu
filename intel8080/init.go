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

	// op table
	table []func(*stepInfo) uint
}

type stepInfo struct{
	PC uint16
	opcode uint8
}

// TODO
var instructionNames = []string{
	"NOP", "LXI", "STAX", "INX", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---",
	"---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---",
	"---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---",
	"---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---",

	"MOV", "MOV", "MOV", "MOV", "MOV", "MOV", "MOV", "MOV", "MOV", "MOV", "MOV", "MOV", "MOV", "MOV", "MOV",
	"MOV", "MOV", "MOV", "MOV", "MOV", "MOV", "MOV", "MOV", "MOV", "MOV", "MOV", "MOV", "MOV", "MOV", "MOV",
	"MOV", "MOV", "MOV", "MOV", "MOV", "MOV", "MOV", "MOV", "MOV", "MOV", "MOV", "MOV", "MOV", "MOV", "MOV",
	"MOV", "MOV", "MOV", "MOV", "MOV", "MOV", "HLT", "MOV", "MOV", "MOV", "MOV", "MOV", "MOV", "MOV", "MOV",

	"---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---",
	"---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---",
	"---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---",
	"---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---",

	"---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---",
	"---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---",
	"---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---",
	"---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---", "---",
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
		nil, cpu.lxi, cpu.SHLD, cpu.inx, cpu.inr, cpu.dcr, cpu.mvi, cpu.daa, nil, cpu.dad, cpu.lhld, cpu.dcx, cpu.inr, cpu.dcr, cpu.mvi, cpu.cma,
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

func (cpu *CPU) add(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) sub(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) ana(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) ora(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) adc(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) sbb(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) xra(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) cmp(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) rnz(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) rnc(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) rpo(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) rp(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) lxi(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) stax(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) SHLD(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) sta(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) inx(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) inr(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) dcr(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) mvi(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) rlc(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) ral(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) daa(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) stc(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) dad(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) ldax(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) lhld(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) lda(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) dcx(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) rrc(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) rar(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) cma(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) cmc(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) pop(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) jnz(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) jnc(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) jpo(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) jp(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) out(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) xthl(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) di(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) cnz(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) cnc(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) cpo(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) cp(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) push(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) ori(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) adi(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) sui(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) ani(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) rst(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) rz(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) rm(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) ret(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) rc(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) rpe(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) sphl(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) jz(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) jc(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) pchl(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) jm(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) in(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) jpe(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) ei(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) cz(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) cc(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) xchg(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) cm(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) call(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) aci(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) sbi(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) cpe(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) xri(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) cpi(info *stepInfo) uint {

	return 0
}
