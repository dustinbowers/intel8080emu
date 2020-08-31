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
	if opcodeFunc == nil {
		return 0, fmt.Errorf("Invalid opcode: 0x%x\n", opcode)
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

	return fmt.Sprintf("PC: %05d, SP: %05d, Flags: %08b, Opcode: 0x%02x (%d) - %s\t[%d %d %d %d %d %d %d]",
		cpu.PC, cpu.SP, cpu.getProgramStatus(), opcode, bytes, name,
		cpu.A, cpu.B, cpu.C, cpu.D, cpu.E, cpu.H, cpu.L)
}
