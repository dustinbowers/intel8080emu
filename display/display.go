package display

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

var (
	width       int32
	height      int32
	rows        int32
	cols        int32
	blockWidth  int32
	blockHeight int32

	colorMask *ColorMask
)

var window *sdl.Window

func Init(screenWidth uint, screenHeight uint, screenCols uint, screenRows uint) {
	if err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_AUDIO); err != nil {
		panic(err)
	}

	width = int32(screenWidth)
	height = int32(screenHeight)
	cols = int32(screenCols)
	rows = int32(screenRows)
	blockWidth = width / cols
	blockHeight = height / rows

	win, err := sdl.CreateWindow("Intel 8080 Emulator", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		width, height, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	window = win
}

func SetColorMask(c *ColorMask) {
	colorMask = c
}

func Draw(cells []byte) error {
	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}
	// if err != nil {
	// 	return fmt.Errorf("draw: FillRect failed: %v", err)
	// }

	for i, b := range cells {
		y := i / 32
		for bit := 0; bit < 8; bit++ {
			x := (i%32)*8 + bit

			xPos := int32(x) * blockWidth
			yPos := int32(y) * blockHeight

			// Yes, it is inefficient to re-draw the entire screen when not needed.
			// It's done to ensure that each frame's blitting ops take approximately
			// the same amount of time to complete per frame
			color := uint(b) & (0x1 << (bit))
			if color > 0 {
				color = 0xffffffff
			}

			rect := sdl.Rect{
				X: xPos,
				Y: yPos,
				W: blockWidth,
				H: blockHeight,
			}
			_ = surface.FillRect(&rect, uint32(color))
		}
	}

	err = window.UpdateSurface()
	if err != nil {
		return fmt.Errorf("draw: UpdateSurface failed: %v", err)
	}
	return nil
}

func DrawRotated(cells []byte) error {
	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}
	err = surface.FillRect(nil, 0)
	if err != nil {
		return fmt.Errorf("draw: FillRect failed: %v", err)
	}

	for i, b := range cells {
		x := i / 32
		for bit := 0; bit < 8; bit++ {
			y := (i%32)*8 + bit

			xPos := int32(x) * blockWidth
			yPos := height - int32(y)*blockHeight

			// Yes, it is inefficient to re-draw the entire screen when not needed.
			// It's done to ensure that each frame's blitting ops take approximately
			// the same amount of time to complete per frame
			on := uint(b) & (0x1 << (bit))

			var color uint32
			if on > 0 {
				masked, ok := colorMask.Color[uint(x)][uint(y)]
				if ok {
					color = masked
				} else {
					color = 0xffffffff
				}
			}

			rect := sdl.Rect{
				X: xPos,
				Y: yPos,
				W: blockWidth,
				H: blockHeight,
			}
			_ = surface.FillRect(&rect, color)
		}
	}

	err = window.UpdateSurface()
	if err != nil {
		return fmt.Errorf("draw: UpdateSurface failed: %v", err)
	}
	return nil
}

func Cleanup() {
	sdl.Quit()
	_ = window.Destroy()
}
