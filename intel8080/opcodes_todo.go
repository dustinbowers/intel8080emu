package intel8080

// TODO stubs
/////////////
///*

func (cpu *CPU) rst(info *stepInfo) uint {
	panic("RST not implemented")
	return 0
}

func (cpu *CPU) in(info *stepInfo) uint {
	db, _ := cpu.getOpcodeArgs(info.PC)
	cpu.ioBus.Read(db)
	return 0
}

func (cpu *CPU) out(info *stepInfo) uint {
	db, _ := cpu.getOpcodeArgs(info.PC)
	cpu.ioBus.Write(db, cpu.A)
	return 0
}

/**/
