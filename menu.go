package main

import (
	"image/color"
	"log"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/otraore/space-shooter/gui"
)

const (
	mainBtnImage        = "buttons/yellow.png"
	mainBtnImageClicked = "buttons/yellow_pressed.png"
	btnImage            = "buttons/blue.png"
	btnImageClicked     = "buttons/blue_pressed.png"
	greyPanelImage      = "ui/grey_panel_hd.png"
	panelHeaderImage    = "ui/blue_panel_header.png"
	uiFont              = "fonts/kenvector_future_thin.ttf"
	gameSheet           = "spritesheets/game.xml"
	btnWidth            = 340
	btnHeight           = 70
	panelWidth          = 420
	panelHeight         = 400
	offsetY             = 5
	btnMargin           = 20
)

var colorYellow = &color.RGBA{R: 255, G: 204, B: 0, A: 255}
var colorDirtyBlue = &color.RGBA{R: 63, G: 124, B: 182, A: 255}

type MenuScene struct{}

func (MenuScene) Preload() {
	err := engo.Files.Load(gameSheet, uiFont)
	if err != nil {
		panic(err)
	}
}

func (MenuScene) Setup(u engo.Updater) {
	w, _ := u.(*ecs.World)

	w.AddSystem(&common.RenderSystem{})

	mainFnt := &common.Font{
		URL:  uiFont,
		FG:   color.RGBA{R: 153, G: 122, B: 0, A: 255},
		Size: 30,
	}

	fnt := &common.Font{
		URL:  uiFont,
		FG:   color.White,
		Size: 25,
	}

	err := fnt.CreatePreloaded()
	handleErr(err)
	err = mainFnt.CreatePreloaded()
	handleErr(err)

	mainTexture := loadedSprite(mainBtnImage)
	mainTextureClicked := loadedSprite(mainBtnImageClicked)
	panelTexture := loadedSprite(greyPanelImage)
	panelHeaderTexture := loadedSprite(panelHeaderImage)

	texture := loadedSprite(btnImage)
	textureClicked := loadedSprite(btnImageClicked)

	x := (engo.GameWidth() / 2) - btnWidth/2
	y := ((engo.GameHeight() / 2) - (btnHeight / 2) - btnHeight/2) + 30

	playBtn, err := gui.NewButton(gui.Button{
		Text:         "Start Game",
		World:        w,
		Image:        mainTexture,
		ImageClicked: mainTextureClicked,
		Font:         mainFnt,
		Position:     engo.Point{X: x, Y: y - btnMargin - btnHeight - 10},
		Width:        btnWidth,
		Height:       btnHeight + 10,
	})

	playBtn.OnClick(func() {
		engo.SetScene(GameScene{}, true)
	})

	optionsBtn, err := gui.NewButton(gui.Button{
		Text:         "Options",
		World:        w,
		Image:        texture,
		ImageClicked: textureClicked,
		Font:         fnt,
		Position:     engo.Point{X: x, Y: y},
		Width:        btnWidth,
		Height:       btnHeight,
		OffsetY:      offsetY,
	})

	optionsBtn.OnClick(func() {
		log.Println("options button")
	})

	exitBtn, err := gui.NewButton(gui.Button{
		Text:         "Exit",
		World:        w,
		Image:        texture,
		ImageClicked: textureClicked,
		Font:         fnt,
		Position:     engo.Point{X: x, Y: y + btnHeight + btnMargin},
		Width:        btnWidth,
		Height:       btnHeight,
		OffsetY:      offsetY,
	})

	offsetX := float32((panelWidth - btnWidth) / 2)
	_, err = gui.NewPanel(gui.Panel{
		Text:         "Space Shooter",
		World:        w,
		HeaderImage:  panelHeaderTexture,
		BodyImage:    panelTexture,
		Font:         fnt,
		Position:     engo.Point{X: x - offsetX, Y: 50},
		Width:        panelWidth,
		Height:       panelHeight,
		HeaderHeight: 73,
	})

	exitBtn.OnClick(func() {
		engo.Exit()
	})

	common.SetBackground(colorYellow)
}

func (MenuScene) Type() string { return "MenuScene" }
