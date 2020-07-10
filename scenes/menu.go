package scenes

import (
	"image/color"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/otraore/space-shooter/gui"
)

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
	w.AddSystem(&gui.ButtonSystem{})
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
	y := ((engo.GameHeight() / 2) - (btnHeight / 2) - btnHeight/2) + 50

	playBtn, _ := gui.NewButton(gui.Button{
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

	optionsBtn, _ := gui.NewButton(gui.Button{
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
		engo.SetScene(OptionsScene{}, true)
	})

	exitBtn, _ := gui.NewButton(gui.Button{
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
