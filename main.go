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

	ioBus := intel8080.NewIOBus()

	memory := intel8080.NewMemory(0x4000)
	romDir := "roms/"
	count, err := memory.LoadRomFiles([]string{
		romDir + "invaders.h",
		romDir + "invaders.g",
		romDir + "invaders.f",
		romDir + "invaders.e",
	})
	fmt.Printf("%d bytes loaded\n", count)
	if err != nil {
		fmt.Printf("LoadRomFiles failed: %v\n", err)
	}

	cpu = intel8080.NewCPU(ioBus, memory)

	if err != nil {
		log.Fatalf("load invaders failed: %v", err)
	}

	running := true

	// Trap SIGINT for debugging purposes
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Printf("(sig %v) Dumping memory... ", sig)
			err := dumpCoreToFile("core.dump", memory)
			if err != nil {
				log.Printf("dumpCoreToFile error: %v\n", err)
			}
			running = false
		}
	}()

	screenCols := 256
	screenRows := 224
	screenWidth := 512
	screenHeight := screenRows * screenWidth / screenCols

	display.Init(screenHeight, screenWidth, screenRows, screenCols)
	defer display.Cleanup()

	vram := cpu.GetVram()
	_ = display.Draw(vram)
	fmt.Println("Starting tick loop")
	var holdCycles uint
	var currCycles uint
	var interruptType uint = 1
	sleepTime := (1000 / 2000) * time.Millisecond
	for running != false {
		if holdCycles > 0 {
			holdCycles--
			time.Sleep(sleepTime)
			continue
		}

		holdCycles, err = cpu.Step()
		if err != nil {
			fmt.Printf("CPU Execution error: %v\n", err)
			running = false
		}
		currCycles += holdCycles

		if currCycles > 16666 {
			currCycles = 0
			// toggle interrupt type between 1 and 2
			if interruptType == 2 {
				interruptType = 1
			} else if interruptType == 1 {
				interruptType = 2
			}
			_ = display.DrawRotated(vram)
			// trigger interrupt (this happens in hardware at VBlank and ~1/2 VBlank)
			cpu.Interrupt(interruptType)

		}
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.KeyboardEvent:
				// Game input
				pressed := false
				if t.Type == sdl.KEYDOWN {
					pressed = true
				} else if t.Type == sdl.KEYUP {
					pressed = false
				}
				switch t.Keysym.Sym {
				case sdl.K_c:
					ioBus.HandleInput(0, pressed)
				case sdl.K_SPACE:
					ioBus.HandleInput(2, pressed)
				case sdl.K_w:
					ioBus.HandleInput(4, pressed)
				case sdl.K_a:
					ioBus.HandleInput(5, pressed)
				case sdl.K_d:
					ioBus.HandleInput(6, pressed)
				}

				// Misc input
				if t.Type == sdl.KEYDOWN {
					switch t.Keysym.Sym {
					case sdl.K_ESCAPE:
						running = false
					case sdl.K_LEFTBRACKET:
						if cpu.DEBUG {
							cpu.DEBUG = false
						} else {
							cpu.DEBUG = true
						}
					case sdl.K_RIGHTBRACKET:
						if ioBus.DEBUG {
							ioBus.DEBUG = false
						} else {
							ioBus.DEBUG = true
						}
					case sdl.K_p: //sdl.K_BACKSLASH:
						if memory.DEBUG {
							memory.DEBUG = false
						} else {
							memory.DEBUG = true
						}
					case sdl.K_COMMA:
						sleepTime += 10 * time.Millisecond
						fmt.Printf("sleepTime: %d\n", sleepTime)
					case sdl.K_PERIOD:
						sleepTime -= 10 * time.Millisecond
						if sleepTime < 0 {
							sleepTime = 0
						}
						fmt.Printf("sleepTime: %d\n", sleepTime)
					}
				}
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}
		time.Sleep(sleepTime)
	}
}

func dumpCoreToFile(filename string, memory *intel8080.Memory) error {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("Failed to open file: %v\n", err)
	}
	defer f.Close()

	_, err = f.Write(memory.GetMemoryCopy())
	if err != nil {
		return fmt.Errorf("Writing to file failed: %v\n", err)
	}
	return nil
}
