package emu

type Screen struct {
	Width  uint8
	Height uint8
	Pixels [32][64]byte
}

func NewSceen() *Screen {
	return &Screen{
		Width:  64,
		Height: 32,
	}
}

func (d *Screen) Clear() {
	for i := uint8(0); i < d.Width; i++ {
		for j := uint8(0); j < d.Height; j++ {
			d.Pixels[j][i] = 0
		}
	}
}
