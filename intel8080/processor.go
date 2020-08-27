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
