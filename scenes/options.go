package scenes

import (
	"image/color"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/otraore/space-shooter/gui"
)

const optionsBtnWidth = 230

type OptionsScene struct{}

func (OptionsScene) Preload() {
	// The assets should've already been loaded
	// err := engo.Files.Load(gameSheet, uiFont)
	// if err != nil {
	// 	panic(err)
	// }
}

func (OptionsScene) Setup(u engo.Updater) {
	w, _ := u.(*ecs.World)

	w.AddSystem(&common.RenderSystem{})
	w.AddSystem(&gui.ButtonSystem{})

	fnt := &common.Font{
		URL:  uiFont,
		FG:   color.White,
		Size: 25,
	}

	err := fnt.CreatePreloaded()
	handleErr(err)

	mainTexture := loadedSprite(mainBtnImage)
	mainTextureClicked := loadedSprite(mainBtnImageClicked)
	// panelTexture := loadedSprite(greyPanelImage)
	// panelHeaderTexture := loadedSprite(panelHeaderImage)

	texture := loadedSprite(btnImage)
	textureClicked := loadedSprite(btnImageClicked)

	var x float32 = 20
	y := engo.GameHeight() - btnHeight - 20

	cancelBtn, _ := gui.NewButton(gui.Button{
		Text:         "Cancel",
		World:        w,
		Image:        texture,
		ImageClicked: textureClicked,
		Font:         fnt,
		Position:     engo.Point{X: x, Y: y},
		Width:        optionsBtnWidth,
		Height:       btnHeight,
	})

	cancelBtn.OnClick(func() {
		engo.SetScene(MenuScene{}, true)
	})

	saveBtn, _ := gui.NewButton(gui.Button{
		Text:         "Save",
		World:        w,
		Image:        mainTexture,
		ImageClicked: mainTextureClicked,
		Font:         fnt,
		Position:     engo.Point{X: x + optionsBtnWidth + 20, Y: y},
		Width:        optionsBtnWidth,
		Height:       btnHeight,
	})

	saveBtn.OnClick(func() {
		err := saveSettings()
		if err != nil {
			handleErr(err)
		}

		engo.SetScene(MenuScene{}, true)
	})

	common.SetBackground(colorDirtyBlue)
	engo.SetCursor(engo.CursorNone)
}

func (OptionsScene) Type() string { return "OptionsScene" }

func saveSettings() error {
	return nil
}
