package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/otraore/space-shooter/gui"
)

const (
	btnImage        = "images/ui/button_silver.png"
	btnImageClicked = "images/ui/button_gold.png"
	uiFont          = "fonts/kenvector_future.ttf"
	gameSheet       = "spritesheets/game.xml"
)

type MenuScene struct{}

func (MenuScene) Preload() {
	err := engo.Files.Load(gameSheet, btnImage, btnImageClicked, uiFont)
	if err != nil {
		panic(err)
	}
}

func (MenuScene) Setup(u engo.Updater) {
	w, _ := u.(*ecs.World)

	w.AddSystem(&common.RenderSystem{})

	fnt := &common.Font{
		URL:  uiFont,
		FG:   color.White,
		Size: 64,
	}

	err := fnt.CreatePreloaded()
	if err != nil {
		panic(err)
	}

	w.AddSystem(&common.FPSSystem{Display: true, Font: fnt})

	texture, err := common.LoadedSprite(btnImage)
	if err != nil {
		log.Println(err)
	}

	textureClicked, err := common.LoadedSprite(btnImageClicked)
	if err != nil {
		log.Println(err)
	}

	x := (engo.GameWidth() / 2) - texture.Width()/2
	y := (engo.GameHeight() / 2) - (texture.Height() / 2) - texture.Height()/2

	fmt.Println(texture.Width())

	playBtn, err := gui.NewButton(gui.Button{
		Text:         "Play",
		World:        w,
		Image:        texture,
		ImageClicked: textureClicked,
		Font:         fnt,
		Position:     engo.Point{X: x, Y: y},
	})

	playBtn.OnClick(func() {
		engo.SetScene(GameScene{}, true)
	})

	exitBtn, err := gui.NewButton(gui.Button{
		Text:         "Exit",
		World:        w,
		Image:        texture,
		ImageClicked: textureClicked,
		Font:         fnt,
		Position:     engo.Point{X: x, Y: y + texture.Height() + 30},
	})

	exitBtn.OnClick(func() {
		engo.Exit()
	})

	err = gui.SetBackgroundImage(w, "backgrounds/darkPurple.png")
	if err != nil {
		log.Println(err)
	}
}

func (MenuScene) Type() string { return "MenuScene" }
