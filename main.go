package main

import (
	"github.com/EngoEngine/engo"
	"github.com/otraore/space-shooter/scenes"
)

func main() {
	opts := engo.RunOptions{
		Title:          "Space Shooter",
		Width:          512,
		Height:         512,
		StandardInputs: true,
		NotResizable:   true,
	}

	engo.Run(opts, &scenes.MenuScene{})
}
