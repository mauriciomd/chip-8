package emu

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Display struct {
	window  *sdl.Window
	rederer *sdl.Renderer
	scale   int32
	width   int32
	height  int32
}

func NewDisplay(scale int32) *Display {
	return &Display{
		width:  64,
		height: 32,
		scale:  scale,
	}
}

func (d *Display) Create() error {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return err
	}

	var err error
	d.window, err = sdl.CreateWindow(
		"CHIP-8",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		int32(d.width)*d.scale,
		int32(d.height)*d.scale,
		sdl.WINDOW_SHOWN,
	)

	if err != nil {
		return err
	}

	d.rederer, err = sdl.CreateRenderer(d.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return err
	}

	return nil
}

func (d *Display) Destroy() {
	d.window.Destroy()
	d.rederer.Destroy()
	sdl.Quit()
}

func (d *Display) Draw(buffer [32][64]byte) {
	d.rederer.SetDrawColor(0, 0, 0, 0)
	d.rederer.Clear()
	d.rederer.SetDrawColor(255, 255, 255, 200)

	for i := 0; i < int(d.height); i++ {
		for j := 0; j < int(d.width); j++ {
			x := int32(j) * d.scale
			y := int32(i) * d.scale

			if buffer[i][j] == 1 {
				rect := sdl.Rect{X: x, Y: y, W: d.scale, H: d.scale}
				d.rederer.FillRect(&rect)
			}
		}
	}
	d.rederer.Present()
}
