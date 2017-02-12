package main

import "engo.io/engo"

func main() {
	opts := engo.RunOptions{
		Title:          "Dodger",
		Width:          1024,
		Height:         640,
		StandardInputs: true,
	}

	engo.Run(opts, &MenuScene{})
}
