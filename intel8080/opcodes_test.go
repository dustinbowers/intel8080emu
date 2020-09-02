package intel8080

import (
	"testing"
)

func TestAci(t *testing.T) {
	var ioBus = NewIOBus()
	var tCpu = NewCPU(ioBus)
	var PC uint16 = 0x4
	var opcode uint8 = 0b11001110
	var immediate uint8 = 0x42
	var stepInfo = stepInfo{PC, opcode}

	tCpu.PC = PC
	tCpu.A = 0x14
	tCpu.setProgramStatus(0b01010111)
	tCpu.Memory[PC+0] = opcode
	tCpu.Memory[PC+1] = immediate

	cycles := tCpu.aci(&stepInfo) // tCpu.Step()

	wantA := uint8(0x57)
	wantCarry := false
	wantCycles := uint(7)

	if tCpu.A != wantA {
		t.Errorf("tCpu.A != wantA - expected: %02x, got: %02x\n", wantA, tCpu.A)
	}
	if tCpu.Carry != wantCarry {
		t.Errorf("tCpu.Carry != wantCarry\n")
	}
	if cycles != wantCycles {
		t.Errorf("cycles != wantCycles\n")
	}
}
