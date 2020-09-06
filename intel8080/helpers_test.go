package intel8080

import (
	"testing"
)

func TestGetOpcodeRegPtr(t *testing.T) {
	ioBus := NewIOBus()
	memory := NewMemory(0xFFFF)
	var tCpu = NewCPU(ioBus, memory)
	tCpu.H = 0xAA
	tCpu.L = 0xBB
	memOffset := (uint16(tCpu.H) << 8) | uint16(tCpu.L)
	tests := []struct {
		inRegIndicator   uint8
		wantPtr          *uint8
		wantMemoryAccess bool
	}{
		{0b111, &tCpu.A, false},
		{0b000, &tCpu.B, false},
		{0b001, &tCpu.C, false},
		{0b010, &tCpu.D, false},
		{0b011, &tCpu.E, false},
		{0b100, &tCpu.H, false},
		{0b101, &tCpu.L, false},
		{0b110, memory.GetOffsetPtr(memOffset), true},
	}
	for i, tt := range tests {
		gotPtr, gotMemoryAccess := tCpu.getOpcodeRegPtr(tt.inRegIndicator)
		if gotPtr != tt.wantPtr ||
			gotMemoryAccess != tt.wantMemoryAccess {
			t.Fatalf("ind: %d\n", i)
		}
	}
}

func TestGetOpcodeArgs(t *testing.T) {
	ioBus := NewIOBus()
	memory := NewMemory(0xFFFF)
	var tCpu = NewCPU(ioBus, memory)
	PC := uint16(100)
	tCpu.PC = PC
	memory.Write(PC+0, 0xAA)
	memory.Write(PC+1, 0xBB)
	memory.Write(PC+2, 0xCC)

	wantLb := uint8(0xBB)
	wantHb := uint8(0xCC)

	gotLb, gotHb := tCpu.getOpcodeArgs(tCpu.PC)

	if gotLb != wantLb ||
		gotHb != wantHb {
		t.Fatalf("expected: (0b%08b, 0b%08b), got: (0b%08b, 0b%08b)\n", wantLb, wantHb, gotLb, gotHb)
	}
}

func TestSetProgramStatus(t *testing.T) {
	ioBus := NewIOBus()
	memory := NewMemory(0xFFFF)
	var tCpu = NewCPU(ioBus, memory)
	tests := []struct {
		inPsw                                                   uint8
		wantSign, wantZero, wantAuxCarry, wantParity, wantCarry bool
	}{
		{0b00000010, false, false, false, false, false},
		{0b10000010, true, false, false, false, false},
		{0b01000010, false, true, false, false, false},
		{0b00010010, false, false, true, false, false},
		{0b00000110, false, false, false, true, false},
		{0b00000011, false, false, false, false, true},
		{0b11010111, true, true, true, true, true},
	}
	for _, tt := range tests {
		tCpu.setProgramStatus(tt.inPsw)
		if tCpu.Sign != tt.wantSign ||
			tCpu.Zero != tt.wantZero ||
			tCpu.AuxCarry != tt.wantAuxCarry ||
			tCpu.Parity != tt.wantParity ||
			tCpu.Carry != tt.wantCarry {
			got := tCpu.getProgramStatus() // not the best idea to use this in a test...
			t.Fatalf("inPSW: 0b%08b, got: 0b%08b\n", tt.inPsw, got)
		}
	}
}

func TestGetProgramStatus(t *testing.T) {
	ioBus := NewIOBus()
	memory := NewMemory(0xFFFF)
	var tCpu = NewCPU(ioBus, memory)
	tests := []struct {
		inSign     bool
		inZero     bool
		inAuxCarry bool
		inParity   bool
		inCarry    bool
		want       uint8
	}{
		// PSW format: S, Z, 0, AC, 0, P, 1, CY
		{false, false, false, false, false, 0b00000010},
		{true, true, true, true, true, 0b11010111},
		{true, false, false, false, false, 0b10000010},
		{false, true, false, false, false, 0b01000010},
		{false, false, true, false, false, 0b00010010},
		{false, false, false, true, false, 0b00000110},
		{false, false, false, false, true, 0b00000011},
	}
	for _, tt := range tests {
		tCpu.Sign = tt.inSign
		tCpu.Zero = tt.inZero
		tCpu.AuxCarry = tt.inAuxCarry
		tCpu.Parity = tt.inParity
		tCpu.Carry = tt.inCarry
		got := tCpu.getProgramStatus()
		if got != tt.want {
			t.Fatalf("expected: 0b%08b, got: 0b%08b\n", tt.want, got)
		}
	}
}

func TestGetParity(t *testing.T) {
	tests := []struct {
		input byte
		want  bool
	}{
		{0b00000000, true},
		{0b00000001, false},
		{0b00000010, false},
		{0b00000100, false},
		{0b00001000, false},
		{0b00010000, false},
		{0b00100000, false},
		{0b01000000, false},
		{0b10000000, false},
		{0b10000001, true},
		{0b10000011, false},
		{0b10000111, true},
		{0b10001111, false},
		{0b10011111, true},
		{0b10111111, false},
		{0b11111111, true},
	}
	for _, tt := range tests {
		got := getParity(tt.input)
		if got != tt.want {
			t.Fatalf("input: 0b%08b, expected: %v, got: %v\n", tt.input, tt.want, got)
		}
	}
}

func TestGetOpcodeRP(t *testing.T) {
	tests := []struct {
		input uint8
		want  uint8
	}{
		//     0b--RP----
		{0b00000000, 0b00},
		{0b11001111, 0b00},
		{0b00010000, 0b01},
		{0b11011111, 0b01},
		{0b00100000, 0b10},
		{0b11101111, 0b10},
		{0b00110000, 0b11},
		{0b11111111, 0b11},
	}
	for _, tt := range tests {
		got := getOpcodeRP(tt.input)
		if got != tt.want {
			t.Fatalf("input: 0b%08b, expected: %v, got: %v\n", tt.input, tt.want, got)
		}
	}
}

func TestGetOpcodeDDDSSS(t *testing.T) {
	tests := []struct {
		input   uint8
		wantDDD uint8
		wantSSS uint8
	}{
		{0b00000000, 0b000, 0b000},
		{0b00111111, 0b111, 0b111},
		{0b00000001, 0b000, 0b001},
		{0b00000010, 0b000, 0b010},
		{0b00000100, 0b000, 0b100},
		{0b00001000, 0b001, 0b000},
		{0b00010000, 0b010, 0b000},
		{0b00100000, 0b100, 0b000},
		{0b00010010, 0b010, 0b010},
	}
	for _, tt := range tests {
		gotDDD, gotSSS := getOpcodeDDDSSS(tt.input)
		if gotDDD != tt.wantDDD {
			t.Fatalf("bad DDD: input 0b%08b, expected: %v, got: %v\n", tt.input, tt.wantDDD, gotDDD)
		}
		if gotSSS != tt.wantSSS {
			t.Fatalf("bad SSS: input 0b%08b, expected: %v, got: %v\n", tt.input, tt.wantSSS, gotSSS)
		}
	}
}

func TestIncRegisterPair(t *testing.T) {
	tests := []struct {
		inputHi uint8
		inputLo uint8
		wantHi  uint8
		wantLo  uint8
	}{
		{inputHi: 0b00000000, inputLo: 0b00000000, wantHi: 0b00000000, wantLo: 0b00000001},
		{inputHi: 0b11111111, inputLo: 0b11111111, wantHi: 0b00000000, wantLo: 0b00000000},
		{inputHi: 0b00000000, inputLo: 0b11111111, wantHi: 0b00000001, wantLo: 0b00000000},
		{inputHi: 0b11111111, inputLo: 0b00000001, wantHi: 0b11111111, wantLo: 0b00000010},
	}
	for _, tt := range tests {
		gotHi, gotLo := incRegisterPair(tt.inputHi, tt.inputLo)
		if gotHi != tt.wantHi || gotLo != tt.wantLo {
			t.Fatalf("input: (0b%08b, 0b%08b), expected: (0b%08b, 0b%08b), got: (0b%08b, 0b%08b)\n", tt.inputHi, tt.inputLo, tt.wantHi, tt.wantLo, gotHi, gotLo)
		}
	}
}

func TestDecRegisterPair(t *testing.T) {
	tests := []struct {
		inputHi uint8
		inputLo uint8
		wantHi  uint8
		wantLo  uint8
	}{
		{inputHi: 0b00000000, inputLo: 0b00000000, wantHi: 0b11111111, wantLo: 0b11111111},
		{inputHi: 0b11111111, inputLo: 0b11111111, wantHi: 0b11111111, wantLo: 0b11111110},
		{inputHi: 0b00000001, inputLo: 0b00000000, wantHi: 0b00000000, wantLo: 0b11111111},
		{inputHi: 0b11111111, inputLo: 0b00000001, wantHi: 0b11111111, wantLo: 0b00000000},
	}
	for _, tt := range tests {
		gotHi, gotLo := decRegisterPair(tt.inputHi, tt.inputLo)
		if gotHi != tt.wantHi || gotLo != tt.wantLo {
			t.Fatalf("input: (0b%08b, 0b%08b), expected: (0b%08b, 0b%08b), got: (0b%08b, 0b%08b)\n", tt.inputHi, tt.inputLo, tt.wantHi, tt.wantLo, gotHi, gotLo)
		}
	}
}
