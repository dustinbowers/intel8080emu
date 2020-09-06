package intel8080

import (
	"testing"
)

func TestNop(t *testing.T) {
	var ioBus = NewIOBus()
	var memory = NewMemory(0xFF)
	var tCpu = NewCPU(ioBus, memory)
	tCpu.PC = 0x4
	var opcode uint8 = 00000000
	memory.Write(tCpu.PC, opcode)
	//var stepInfo = stepInfo{PC, opcode}

	//tCpu.memory.Write(PC+0, opcode)
	cycles, err := tCpu.Step()
	if err != nil {
		t.Fail()
	}
	//cycles := tCpu.nop(&stepInfo)

	wantCycles := uint(4)
	if cycles != wantCycles {
		t.Errorf("cycles != wantCycles\n")
	}
}

func TestAci(t *testing.T) {
	var ioBus = NewIOBus()
	var memory = NewMemory(0xFF)
	var tCpu = NewCPU(ioBus, memory)
	var PC uint16 = 0x4
	var opcode uint8 = 0b11001110
	var immediate uint8 = 0x42
	var stepInfo = stepInfo{PC, opcode}

	tCpu.PC = PC
	tCpu.A = 0x14
	tCpu.setProgramStatus(0b01010111)
	tCpu.memory.Write(PC+0, opcode)
	tCpu.memory.Write(PC+1, immediate)

	cycles := tCpu.aci(&stepInfo)

	wantA := uint8(0x57)
	wantCarry := false
	wantCycles := uint(7)

	if tCpu.A != wantA {
		t.Errorf("tCpu.A != wantA - expected: 0x%02x, got: 0x%02x\n", wantA, tCpu.A)
	}
	if tCpu.Carry != wantCarry {
		t.Errorf("tCpu.Carry != wantCarry\n")
	}
	if cycles != wantCycles {
		t.Errorf("cycles != wantCycles\n")
	}
}

func TestDaa(t *testing.T) {
	var ioBus = NewIOBus()
	var memory = NewMemory(0xFF)
	var tCpu = NewCPU(ioBus, memory)
	var PC uint16 = 0x4
	var opcode uint8 = 0b00100111
	var stepInfo = stepInfo{PC, opcode}

	tCpu.PC = PC
	tCpu.A = 0x9b
	tCpu.setProgramStatus(0b00000010)
	tCpu.memory.Write(PC+0, opcode)

	cycles := tCpu.daa(&stepInfo)

	wantA := uint8(0x01)
	wantCarry := true
	wantAuxCarry := true
	wantCycles := uint(4)

	if tCpu.A != wantA {
		t.Errorf("tCpu.A != wantA - expected: 0x%02x, got: 0x%02x\n", wantA, tCpu.A)
	}
	if tCpu.Carry != wantCarry {
		t.Errorf("tCpu.Carry != wantCarry\n")
	}
	if tCpu.Carry != wantAuxCarry {
		t.Errorf("tCpu.AuxCarry != wantAuxCarry\n")
	}
	if cycles != wantCycles {
		t.Errorf("cycles != wantCycles\n")
	}
}
