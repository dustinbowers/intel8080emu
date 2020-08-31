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

func main() {
	fmt.Println("Launching...")

	cpu = intel8080.NewCPU()
	cpu.DEBUG = true
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
	vram := cpu.Memory[0x2400:0x4000]
	for running {
		err := display.Draw(vram)
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
		fmt.Printf(".")
		time.Sleep(1000/60 * time.Millisecond)
	}

}

func startCPU() {
	go func() {
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
	}()
}