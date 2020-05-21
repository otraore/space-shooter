package main

import (
	"github.com/EngoEngine/engo"
)

func main() {
	opts := engo.RunOptions{
		Title:          "Space Shooter",
		Width:          900,
		Height:         550,
		StandardInputs: true,
		NotResizable:   true,
	}

	engo.Run(opts, &MenuScene{})
}
