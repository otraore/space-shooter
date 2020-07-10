package scenes

import (
	"image/color"
	"log"

	"github.com/EngoEngine/engo/common"
)

func loadedSprite(url string) *common.Texture {
	tex, err := common.LoadedSprite(url)
	if err != nil {
		handleErr(err)
	}

	return tex
}

func handleErr(err error) {
	if err != nil {
		log.Fatal("error: ", err)
	}
}

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
