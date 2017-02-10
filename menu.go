package main

import (
	"image/color"
	"log"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
)

var playBtn *Button

type MenuScene struct{}

func (*MenuScene) Preload() {
	engo.Files.Load("images/button_silver.png", "images/button_gold.png", "fonts/kenvector_future.ttf")
}

func (*MenuScene) Setup(w *ecs.World) {
	common.SetBackground(color.White)

	w.AddSystem(&common.RenderSystem{})

	fnt := &common.Font{
		URL:  "fonts/kenvector_future.ttf",
		FG:   color.White,
		Size: 64,
	}
	err := fnt.CreatePreloaded()
	if err != nil {
		panic(err)
	}

	//Retrieve a texture
	texture, err := common.LoadedSprite("images/button_silver.png")
	if err != nil {
		log.Println(err)
	}

	textureClicked, err := common.LoadedSprite("images/button_gold.png")
	if err != nil {
		log.Println(err)
	}

	playBtn = NewButton(w, texture, textureClicked, fnt, "Play")

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&playBtn.Graphic.BasicEntity, &playBtn.Graphic.RenderComponent, &playBtn.Graphic.SpaceComponent)
			sys.Add(&playBtn.Label.BasicEntity, &playBtn.Label.RenderComponent, &playBtn.Label.SpaceComponent)
		}
	}

	// Create an entity

}

func (*MenuScene) Type() string { return "MenuScene" }
