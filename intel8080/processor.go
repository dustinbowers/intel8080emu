package intel8080

import (
	"fmt"
	"log"
)

func (cpu *CPU) Step() (uint, error) {
	opcode := cpu.memory.Read(cpu.PC)
	opcodeFunc := cpu.table[opcode]
	stepInfo := stepInfo{
		PC:     cpu.PC,
		opcode: opcode,
	}

	if cpu.DEBUG {
		log.Println(cpu.GetInstructionInfo())
	}
	if opcodeFunc == nil {
		return 0, fmt.Errorf("Invalid opcode: 0x%02x\n (%s)", opcode, instructionNames[opcode])
	}

	if cpu.deferInterruptsEnable {
		cpu.deferInterruptsEnable = false
		cpu.InterruptsEnabled = true
	}

	// Execute current opcode
	pcAdvanceAmt := uint16(instructionBytes[opcode])
	if pcAdvanceMask[opcode] == 1 {
		cpu.PC += pcAdvanceAmt
	}
	cycles := opcodeFunc(&stepInfo)

	return cycles, nil
}

func (cpu *CPU) GetInstructionInfo() string {
	opcode := cpu.memory.Read(cpu.PC)
	bytes := instructionBytes[opcode]
	name := instructionNames[opcode]

	args := "---- ----"
	if bytes == 2 {
		args = fmt.Sprintf("0x%02x ----", cpu.memory.Read(cpu.PC+1))
	}
	if bytes == 3 {
		args = fmt.Sprintf("0x%02x 0x%02x", cpu.memory.Read(cpu.PC+1), cpu.memory.Read(cpu.PC+2))
	}
	return fmt.Sprintf("PC: %04x, SP: %04x, Flags: %08b, Regs: [%02x %02x %02x %02x %02x %02x %02x], Opcode: 0b%08b / 0x%02x (%d) - %s,\tArgs [%s]",
		cpu.PC, cpu.SP, cpu.getProgramStatus(),
		cpu.A, cpu.B, cpu.C, cpu.D, cpu.E, cpu.H, cpu.L,
		opcode, opcode, bytes, name, args)
}
