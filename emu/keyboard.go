package emu

type Keyboard struct {
	keys [16]bool
}

func NewKeyboard() *Keyboard {
	return &Keyboard{}
}

func (k *Keyboard) Set(key byte, pressed bool) {
	k.keys[key] = pressed
}

func (k *Keyboard) IsPressed(key byte) bool {
	return k.keys[key]
}

func (k *Keyboard) IsAnyKeyPressed() (bool, uint8) {
	for idx, key := range k.keys {
		if key {
			return true, uint8(idx)
		}
	}

	return false, 0
}
