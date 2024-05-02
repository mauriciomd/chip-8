package emu

import (
	"fmt"
	"math/rand"
)

type CPU struct {
	Memory   *Memory
	Stack    *Stack
	Screen   *Screen
	Keyboard *Keyboard

	DelayTimer byte
	SoundTimer byte
	PC         uint16
	I          uint16
	Registers  [16]byte

	LagacyMode bool
	Paused     bool
}

func NewCpu(m *Memory, s *Stack, d *Screen, k *Keyboard) *CPU {
	return &CPU{
		Memory:     m,
		Screen:     d,
		Stack:      s,
		Keyboard:   k,
		PC:         0x200,
		LagacyMode: false,
		Paused:     false,
	}
}

func (c *CPU) decreaseTimers() {
	if c.DelayTimer > 0 {
		c.DelayTimer -= 1
	}

	if c.SoundTimer > 0 {
		c.SoundTimer -= 1
	}
}

func (c *CPU) ExecuteNextInstruction() {
	op := c.fetch()
	fmt.Printf("Executing [%.4x] - PC [%d]\n", op, c.PC)
	c.execute(op)
	c.decreaseTimers()
}

func (c *CPU) fetch() uint16 {
	firstByte := uint16(c.Memory[c.PC])
	secondByte := uint16(c.Memory[c.PC+1])

	if !c.Paused {
		c.PC += 2
	}
	instruction := firstByte<<8 + secondByte
	return instruction
}

func (c *CPU) execute(instruction uint16) {
	op := c.decodeType(instruction)

	switch op {
	case 0x0000:
		switch instruction {
		case 0x00E0:
			c.Screen.Clear()

		case 0x00EE:
			c.PC = c.Stack.Pop()
		default:
			fmt.Printf("Executing instruction 0x0NNN => %.4x\n", instruction)
		}
	case 0x1000:
		c.PC = c.decodeNNN(instruction)
	case 0x2000:
		c.Stack.Push(c.PC)
		c.PC = c.decodeNNN(instruction)
	case 0x3000:
		r, value := c.decodeXNN(instruction)
		if c.getRegisterValue(r) == value {
			c.PC += 2
		}
	case 0x4000:
		r, value := c.decodeXNN(instruction)
		if c.getRegisterValue(r) != value {
			c.PC += 2
		}
	case 0x5000:
		r1, r2, _ := c.decodeXYN(instruction)
		vx := c.getRegisterValue(r1)
		vy := c.getRegisterValue(r2)

		if vx == vy {
			c.PC += 2
		}
	case 0x6000:
		register, value := c.decodeXNN(instruction)
		c.setRegisterValue(register, value)
	case 0x7000:
		register, value := c.decodeXNN(instruction)
		current := c.getRegisterValue(register)
		c.setRegisterValue(register, current+value)
	case 0x8000:
		t := instruction & 0x000F
		r1, r2, _ := c.decodeXYN(instruction)
		switch t {
		case 0x0000:
			vy := c.getRegisterValue(r2)
			c.setRegisterValue(r1, vy)

		case 0x0001:
			vx := c.getRegisterValue(r1)
			vy := c.getRegisterValue(r2)

			c.setRegisterValue(r1, vx|vy)

		case 0x0002:
			vx := c.getRegisterValue(r1)
			vy := c.getRegisterValue(r2)

			c.setRegisterValue(r1, vx&vy)

		case 0x0003:
			vx := c.getRegisterValue(r1)
			vy := c.getRegisterValue(r2)

			c.setRegisterValue(r1, vx^vy)

		case 0x0004:
			vx := c.getRegisterValue(r1)
			vy := c.getRegisterValue(r2)
			vf := byte(0)
			sum := uint16(vx) + uint16(vy)
			if sum > 0xFF {
				vf = 1
			}

			c.setRegisterValue(r1, vx+vy)
			c.setRegisterValue(0xF, vf)

		case 0x0005:
			vy := c.getRegisterValue(r2)
			vx := c.getRegisterValue(r1)
			vf := byte(0)

			if vx > vy {
				vf = 1
			}

			c.setRegisterValue(r1, vx-vy)
			c.setRegisterValue(0xF, vf)

		case 0x0006:
			if !c.LagacyMode {
				vy := c.getRegisterValue(r2)
				c.setRegisterValue(r1, vy)
			}

			vx := c.getRegisterValue(r1)
			vf := vx & 0x1

			c.setRegisterValue(r1, vx>>1)
			c.setRegisterValue(0xF, vf)

		case 0x0007:
			vx := c.getRegisterValue(r1)
			vy := c.getRegisterValue(r2)
			vf := byte(0)

			if vy > vx {
				vf = 1
			}

			c.setRegisterValue(r1, vy-vx)
			c.setRegisterValue(0xF, vf)

		case 0x000E:
			if !c.LagacyMode {
				vy := c.getRegisterValue(r2)
				c.setRegisterValue(r1, vy)
			}

			vx := c.getRegisterValue(r1)
			vf := vx & 0x8

			c.setRegisterValue(r1, vx<<1)
			c.setRegisterValue(0xF, vf)
		}

	case 0x9000:
		r1, r2, _ := c.decodeXYN(instruction)
		vx := c.getRegisterValue(r1)
		vy := c.getRegisterValue(r2)

		if vx != vy {
			c.PC += 2
		}

	case 0xA000:
		c.I = c.decodeNNN(instruction)

	case 0xB000:
		v0 := uint16(c.getRegisterValue(0x0))
		value := c.decodeNNN(instruction)
		c.PC = uint16(v0 + value)

	case 0xC000:
		r, value := c.decodeXNN(instruction)
		number := rand.Uint32() % uint32(value)
		vx := uint8(number) & value
		c.setRegisterValue(r, vx)

	case 0xD000:
		r1, r2, value := c.decodeXYN(instruction)
		rx := c.getRegisterValue(r1)
		ry := c.getRegisterValue(r2)
		c.setRegisterValue(0xF, 0)

		for y := uint16(0); y < uint16(value); y++ {
			pixel := c.Memory[c.I+y]
			py := (uint16(ry) + y) % uint16(c.Screen.Height)

			for x := uint16(0); x < 8; x++ {
				if pixel&(0x80>>x) != 0 {
					px := (uint16(rx) + x) % uint16(c.Screen.Width)
					if c.Screen.Pixels[py][px] == 1 {
						c.setRegisterValue(0xF, 1)
					}
					c.Screen.Pixels[py][px] ^= 0x01
				}

			}
		}

	case 0xE000:
		r, v1, v2 := c.decodeXYN(instruction)
		vx := c.getRegisterValue(r)

		i := v1<<4 + v2

		switch i {
		case 0x9E:
			if c.Keyboard.IsPressed(vx) {
				c.PC += 2
			}
		case 0xA1:
			if !c.Keyboard.IsPressed(vx) {
				c.PC += 2
			}
		}
	case 0xF000:
		r, v1, v2 := c.decodeXYN(instruction)
		vx := c.getRegisterValue(r)

		i := v1<<4 + v2

		switch i {
		case 0x07:
			c.setRegisterValue(r, c.DelayTimer)
		case 0x15:
			c.DelayTimer = vx
		case 0x18:
			c.SoundTimer = vx
		case 0x29:
			c.I = uint16(vx)*5 + 0x50
		case 0x33:
			h := vx / 100
			c.Memory[c.I] = h
			c.Memory[c.I+1] = (vx - h*100) / 10
			c.Memory[c.I+2] = (vx - h*100) % 10
		case 0x55:
			for i := uint8(0); i <= r; i++ {
				addr := c.I + uint16(i)
				c.Memory[addr] = c.getRegisterValue(i)
			}
		case 0x65:
			for i := uint8(0); i <= r; i++ {
				addr := c.I + uint16(i)
				c.setRegisterValue(i, c.Memory[addr])
			}
		case 0x0A:
			isPressed, key := c.Keyboard.IsAnyKeyPressed()
			fmt.Println(isPressed, key)
			if !isPressed {
				c.Paused = true
			} else {
				c.Paused = false
				c.setRegisterValue(r, key)
			}
		case 0x1E:
			c.I += uint16(vx)
		}
	}
}

func (c *CPU) decodeType(instruction uint16) uint16 {
	return instruction & 0xF000
}

func (c *CPU) decodeXNN(instruction uint16) (byte, byte) {
	r := byte(instruction & 0x0F00 >> 8)
	value := byte(instruction & 0x00FF)

	return r, value
}

func (c *CPU) decodeNNN(instruction uint16) uint16 {
	return instruction & 0x0FFF
}

func (c *CPU) decodeXYN(instruction uint16) (byte, byte, byte) {
	r1 := byte(instruction & 0x0F00 >> 8)
	r2 := byte(instruction & 0x00F0 >> 4)
	value := byte(instruction & 0x000F)

	return r1, r2, value
}

func (c *CPU) setRegisterValue(register, value byte) {
	c.Registers[register] = value
}

func (c *CPU) getRegisterValue(register byte) byte {
	return c.Registers[register]
}
