package main

import (
	"github.com/EngoEngine/engo"
)

func main() {
	opts := engo.RunOptions{
		Title:          "Space Shooter",
		Width:          512,
		Height:         512,
		StandardInputs: true,
		NotResizable:   true,
	}

	engo.Run(opts, &MenuScene{})
}
