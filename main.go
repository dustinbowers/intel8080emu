package main

import (
	"fmt"
	"intel8080/intel8080"
	"log"
	"time"
)

func main() {
	fmt.Println("Hello, Invaders!")

	cpu := intel8080.NewCPU()
	cpu.DEBUG = true
	err := cpu.LoadInvaders("roms/")
	if err != nil {
		log.Fatalf("load invaders failed: %v", err)
	}

	fmt.Println("Starting tick loop")
	holdCycles := 0
	for {
		if holdCycles > 0 {
			holdCycles--
			continue
		}

		log.Println("Tick...")
		cpu.Step()

		// TODO: update this to 2 MHz
		time.Sleep((1000 / 30) * time.Millisecond)
	}

}
