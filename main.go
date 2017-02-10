package main

import "engo.io/engo"

var guy Guy
var playing = true

func main() {
	opts := engo.RunOptions{
		Title:          "Dodger",
		Width:          1024,
		Height:         640,
		StandardInputs: true,
	}

	engo.Run(opts, &MenuScene{})
}
