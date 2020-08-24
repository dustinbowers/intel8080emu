package intel8080

import (
	"fmt"
	"log"
)

func (cpu *CPU) Step() (uint, error) {
	opcode := cpu.Read(cpu.PC)
	opcodeFunc := cpu.table[opcode]
	stepInfo := stepInfo{
		PC:     cpu.PC,
		opcode: opcode,
	}

	if cpu.DEBUG {
		log.Println(cpu.GetInstructionInfo())
	}

	// Execute current opcode
	cpu.PC += uint16(instructionBytes[opcode])
	cycles := opcodeFunc(&stepInfo)

	return cycles, nil
}

func (cpu *CPU) Read(pc uint16) uint8 {
	return cpu.Memory[pc]
}

func (cpu *CPU) GetInstructionInfo() string {
	opcode := cpu.Read(cpu.PC)
	bytes := instructionBytes[opcode]
	name := instructionNames[opcode]

	return fmt.Sprintf("0x%x (%d) - %s", opcode, bytes, name)
}

// TODO stubs
/////////////
func (cpu *CPU) sub(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) ana(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) ora(info *stepInfo) uint {

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

func (cpu *CPU) SHLD(info *stepInfo) uint {

	return 0
}

func (cpu *CPU) inx(info *stepInfo) uint {

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
