package main

import (
	"bytes"
	"flag"
	"io"
	"os"

	"github.com/mauriciomd/go-chip-8/emu"
)

func parseArgs() (int, string) {
	scale := flag.Int("scale", 10, "Set display scale.")
	file := flag.String("file", "./roms/test-opcode.ch8", "Rom to be loaded.")
	flag.Parse()

	return *scale, *file
}

func ReadFile(filename string) []byte {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	data, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	return data
}

func main() {
	scale, filename := parseArgs()
	emu, err := emu.New(int32(scale))
	if err != nil {
		panic(err)
	}

	emu.LoadRom(bytes.NewReader(ReadFile(filename)))
	emu.Run()
}
