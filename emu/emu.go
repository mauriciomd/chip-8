package emu

import (
	"io"

	"github.com/veandco/go-sdl2/sdl"
)

type Chip8Emu struct {
	isRunning bool
	cpu       *CPU
	display   *Display
	screen    *Screen
	memory    *Memory
	keyboard  *Keyboard
}

func New(scale int32) (*Chip8Emu, error) {
	display := NewDisplay(scale)
	if err := display.Create(); err != nil {
		return nil, err
	}

	memory, stack := NewMemory()
	screen := NewSceen()
	keyboard := NewKeyboard()
	cpu := NewCpu(memory, stack, screen, keyboard)

	emu := &Chip8Emu{
		isRunning: true,
		cpu:       cpu,
		display:   display,
		screen:    screen,
		memory:    memory,
		keyboard:  keyboard,
	}

	return emu, nil
}

func (c *Chip8Emu) LoadRom(r io.Reader) {
	c.memory.LoadRom(r)
}

func (c *Chip8Emu) Run() {
	defer c.display.Destroy()

	for c.isRunning {
		c.checkSDLEvents()
		c.cpu.ExecuteNextInstruction()
		c.display.Draw(c.screen.Pixels)

		sdl.Delay(3)
	}
}

func (c *Chip8Emu) checkSDLEvents() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.QuitEvent:
			println("Quit")
			c.isRunning = false
		case *sdl.KeyboardEvent:
			t := event.GetType()
			key := event.(*sdl.KeyboardEvent).Keysym.Scancode

			switch key {
			case sdl.SCANCODE_1:
				c.keyboard.Set(0x1, t == sdl.KEYDOWN)

			case sdl.SCANCODE_2:
				c.keyboard.Set(0x2, t == sdl.KEYDOWN)

			case sdl.SCANCODE_3:
				c.keyboard.Set(0x3, t == sdl.KEYDOWN)

			case sdl.SCANCODE_4:
				c.keyboard.Set(0xC, t == sdl.KEYDOWN)

			case sdl.SCANCODE_Q:
				c.keyboard.Set(0x4, t == sdl.KEYDOWN)

			case sdl.SCANCODE_W:
				c.keyboard.Set(0x5, t == sdl.KEYDOWN)

			case sdl.SCANCODE_E:
				c.keyboard.Set(0x6, t == sdl.KEYDOWN)

			case sdl.SCANCODE_R:
				c.keyboard.Set(0xD, t == sdl.KEYDOWN)

			case sdl.SCANCODE_A:
				c.keyboard.Set(0x7, t == sdl.KEYDOWN)

			case sdl.SCANCODE_S:
				c.keyboard.Set(0x8, t == sdl.KEYDOWN)

			case sdl.SCANCODE_D:
				c.keyboard.Set(0x9, t == sdl.KEYDOWN)

			case sdl.SCANCODE_F:
				c.keyboard.Set(0xE, t == sdl.KEYDOWN)

			case sdl.SCANCODE_Z:
				c.keyboard.Set(0xA, t == sdl.KEYDOWN)

			case sdl.SCANCODE_X:
				c.keyboard.Set(0x0, t == sdl.KEYDOWN)

			case sdl.SCANCODE_C:
				c.keyboard.Set(0xB, t == sdl.KEYDOWN)

			case sdl.SCANCODE_V:
				c.keyboard.Set(0xF, t == sdl.KEYDOWN)

			case sdl.SCANCODE_K:
				if t == sdl.KEYUP {
					c.isRunning = false
				}
			}
		}
	}
}
