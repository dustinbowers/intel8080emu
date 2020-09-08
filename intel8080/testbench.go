package intel8080

import (
	"fmt"
	"os"
)

var cpu *CPU

func RunTestRom(testRomPath string) {

	ioBus := NewIOBus()
	memory := NewMemory(0xFFFF)
	count, err := memory.LoadRomFiles([]string{
		testRomPath,
	}, 0x100, false)
	if err != nil {
		fmt.Printf("Error loading ROM file: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("%d bytes loaded\n", count)
	if err != nil {
		fmt.Printf("LoadRomFiles failed: %v\n", err)
	}

	cpu = NewCPU(ioBus, memory)

	// Debug setup
	cpu.PC = 0x0100
	retOpcode := byte(0b11001001) // RET
	inOpcode := byte(0b11011011)  // IN pa
	outOpcode := byte(0b11010011) // OUT pa
	cpu.memory.Write(0x0000, inOpcode)
	cpu.memory.Write(0x0005, outOpcode)
	cpu.memory.Write(0x0007, retOpcode)

	inCallback := func(info *stepInfo) {
		os.Exit(0)
	}
	outCallback := func(info *stepInfo) {
		//C = 0x02 signals printing the value of register E as an ASCII value
		//C = 0x09 signals printing the value of memory pointed to by DE until a '$' character is encountered
		switch cpu.C {
		case 0x02:
			fmt.Printf("%s", string(cpu.E))
		case 0x09:
			start := (uint16(cpu.D) << 8) | uint16(cpu.E)
			end := start
			for {
				c := cpu.memory.Read(end)
				if string(c) == "$" {
					break
				}
				fmt.Printf("%s", string(c))
				end++
			}
		}
	}

	// Set callbacks on IN / OUT
	cpu.inCallback = inCallback
	cpu.outCallback = outCallback

	// Run the test
	for {
		_, _ = cpu.Step()
	}
}
