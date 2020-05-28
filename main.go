package main

import (
	"github.com/EngoEngine/engo"
)

func main() {
	opts := engo.RunOptions{
		Title:          "Space Shooter",
		Width:          256 * 2,
		Height:         256 * 2,
		StandardInputs: true,
		NotResizable:   true,
	}

	engo.Run(opts, &MenuScene{})
}
