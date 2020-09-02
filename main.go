package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"intel8080/display"
	"intel8080/intel8080"
	"log"
	"os"
	"os/signal"
	"time"
)

var cpu *intel8080.CPU

var renderFlag bool

func main() {
	fmt.Println("Launching...")

	ioBus := intel8080.NewIOBus()
	cpu = intel8080.NewCPU(ioBus)
	//cpu.DEBUG = true
	err := cpu.LoadInvaders("roms/")
	if err != nil {
		log.Fatalf("load invaders failed: %v", err)
	}

	// Trap SIGINT for debugging purposes
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func(){
		for sig := range c {
			log.Printf("(sig %v) Dumping stuff... ", sig)
			// TODO:
			os.Exit(1)
		}
	}()

	screenCols := 256
	screenRows := 224
	screenWidth := 512
	screenHeight := screenRows * screenWidth / screenCols

	display.Init(screenWidth, screenHeight, screenCols, screenRows)
	defer display.Cleanup()

	startCPU()

	running := true
	//vram := cpu.Memory[0x2400:0x4000]
	for running {

		if err != nil {
			log.Printf("display draw error: %v", err)
			os.Exit(1)
		}

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}
		//fmt.Printf(".")
		time.Sleep(1000/60 * time.Millisecond)
	}

}

func startCPU() {
	go func() {
		vram := cpu.Memory[0x2400:0x4000]
		fmt.Println("Starting tick loop")
		var holdCycles uint
		var currCycles uint
		var interruptType uint = 1
		sleepTime := (1000 / 2000) * time.Millisecond
		for {
			if holdCycles > 0 {
				holdCycles--
				time.Sleep(sleepTime)
				continue
			}

			holdCycles, _ = cpu.Step()
			currCycles += holdCycles


			if currCycles > 16666 {
				currCycles = 0
				// toggle interrupt type between 1 and 2
				if interruptType == 2 {
					interruptType = 1
					_ = display.Draw(vram)
					fmt.Print(".")
				} else if interruptType == 1 {
					interruptType = 2
				}
				// trigger interrupt (this happens in hardware at VBlank and ~1/2 VBlank)
				cpu.Interrupt(interruptType)
			}

			// TODO: update this to 2 MHz
			time.Sleep(sleepTime)

		}
	}()
}