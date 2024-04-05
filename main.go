package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"intel8080/display"
	"intel8080/intel8080"
)

var testRomPath = flag.String("test", "", "Run a test ROM")

var cpu *intel8080.CPU

func main() {
	fmt.Println("Launching...")

	flag.Parse()
	if *testRomPath != "" {
		fmt.Printf("Running a test ROM - %s\n", *testRomPath)
		intel8080.RunTestRom(*testRomPath)
		return
	}

	memory := intel8080.NewMemory(0x4000)
	romDir := "roms/"
	count, err := memory.LoadRomFiles([]string{
		romDir + "invaders.h",
		romDir + "invaders.g",
		romDir + "invaders.f",
		romDir + "invaders.e",
	}, 0x0, true)
	fmt.Printf("%d bytes loaded\n", count)
	if err != nil {
		fmt.Printf("LoadRomFiles failed: %v\n", err)
		os.Exit(1)
	}

	ioBus := intel8080.NewIOBus()
	cpu = intel8080.NewCPU(ioBus, memory)

	if err != nil {
		log.Fatalf("load invaders failed: %v", err)
		os.Exit(0)
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

	screenCols := uint(256)
	screenRows := uint(224)
	screenWidth := uint(512)
	screenHeight := screenRows * screenWidth / screenCols

	display.Init(screenHeight, screenWidth, screenRows, screenCols)
	defer display.Cleanup()
	cm := display.NewColorMask()
	cm.AddBoxMask(0, screenCols, 48, 64, 0xff00ff00)   // Green
	cm.AddBoxMask(0, screenCols, 192, 224, 0xffff0000) // Red
	display.SetColorMask(cm)

	vram := cpu.GetVram()
	_ = display.DrawRotated(vram)
	fmt.Println("Starting CPU")
	var holdCycles uint
	var currCycles uint
	var interruptType uint = 1
	sleepTime := time.Duration(0)
	for running {
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
			if interruptType == 2 {
				interruptType = 1
			} else if interruptType == 1 {
				interruptType = 2
			}
			_ = display.DrawRotated(vram)
			// trigger interrupt (this happens in hardware at VBlank and ~1/2 VBlank)
			cpu.Interrupt(interruptType)

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
						ioBus.HandleInput(1, 0, pressed)
					case sdl.K_SPACE:
						ioBus.HandleInput(1, 2, pressed)
					case sdl.K_w:
						ioBus.HandleInput(1, 4, pressed)
					case sdl.K_a:
						ioBus.HandleInput(1, 5, pressed)
					case sdl.K_d:
						ioBus.HandleInput(1, 6, pressed)
					case sdl.K_t:
						ioBus.HandleInput(2, 2, pressed)
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
						case sdl.K_p:
							if memory.DEBUG {
								memory.DEBUG = false
							} else {
								memory.DEBUG = true
							}
						case sdl.K_COMMA:
							sleepTime += 1 * time.Nanosecond
							fmt.Printf("sleepTime: %d\n", sleepTime)
						case sdl.K_PERIOD:
							sleepTime -= 1 * time.Nanosecond
							if sleepTime < 0 {
								sleepTime = 0
							}
							fmt.Printf("sleepTime: %d\n", sleepTime)
						}
					}
				case *sdl.QuitEvent:
					println("Quit")
					running = false
				}
			}
		}
		time.Sleep(sleepTime)
	}
}

func dumpCoreToFile(filename string, memory *intel8080.Memory) error {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer f.Close()

	_, err = f.Write(memory.GetMemoryCopy())
	if err != nil {
		return fmt.Errorf("writing to file failed: %v", err)
	}
	return nil
}
