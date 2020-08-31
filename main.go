package main

import (
	"fmt"
	"intel8080/intel8080"
	"log"
	"time"
)

func main() {
	fmt.Println("Launching...")

	cpu := intel8080.NewCPU()
	cpu.DEBUG = true
	err := cpu.LoadInvaders("roms/")
	if err != nil {
		log.Fatalf("load invaders failed: %v", err)
	}

	fmt.Println("Starting tick loop")
	var holdCycles uint
	sleepTime := (1000 / 1000) * time.Millisecond
	for {
		if holdCycles > 0 {
			holdCycles--
			time.Sleep(sleepTime)
			continue
		}

		holdCycles, _ = cpu.Step()

		// TODO: update this to 2 MHz
		time.Sleep(sleepTime)

	}

}
