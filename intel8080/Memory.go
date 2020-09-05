package intel8080

import (
	"fmt"
	"io/ioutil"
)

type Memory struct {
	bytes []byte
	readOnlyBlocks []protectedBlocks
}

type protectedBlocks struct {
	start uint16
	end uint16
}

func NewMemory(size uint16) *Memory {
	m := Memory{}
	m.bytes = make([]byte, size)
	return &m
}

func (m *Memory) Protect(startAddress, endAddress uint16) {
	block := protectedBlocks { startAddress, endAddress }
	m.readOnlyBlocks = append(m.readOnlyBlocks, block)
}

func (m *Memory) GetOffsetPtr(address uint16) *byte {
	return &m.bytes[address]
}

func (m *Memory) GetMemorySlice(start uint16, end uint16) []byte {
	return m.bytes[start:end]
}

func (m* Memory) GetMemoryCopy() []byte {
	bytesCopy := make([]byte, len(m.bytes))
	copy(m.bytes, bytesCopy)
	return bytesCopy
}

func (m *Memory) Read(address uint16) byte {
	return m.bytes[address]
}

func (m *Memory) Write(address uint16, b byte) {
	// Protect read only boundaries
	for _, protectedBlock := range m.readOnlyBlocks {
		start := protectedBlock.start
		end := protectedBlock.end
		if address >= start && address <= end {
			panic(fmt.Sprintf("write to read-only memory location %d in protected block [%d:%d]\n", address, start, end))
		}
	}

	m.bytes[address] = b
}

func (m *Memory) LoadRomFiles(filenames []string) (uint, error) {
	offset := uint(0)
	for _, romPath := range filenames {
		fmt.Printf("loading %s\n", romPath)
		data, err := ioutil.ReadFile(romPath)
		if err != nil {
			return 0, fmt.Errorf("loadRom: failed reading file: %v", err)
		}

		for _, b := range data {
			m.Write(uint16(offset), b)
			offset++
		}
	}
	m.Protect(0, uint16(offset))
	return offset, nil
}
