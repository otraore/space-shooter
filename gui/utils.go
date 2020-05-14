package gui

import (
	"log"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type BackgroundImage struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

func SetBackgroundImage(w *ecs.World, imageUrl string) error {
	bg := &BackgroundImage{BasicEntity: ecs.NewBasic()}

	texture, err := common.LoadedSprite(imageUrl)
	if err != nil {
		log.Printf("Unable to load background image: %v \n", err)
	}

	bg.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: 1, Y: 1},
		Width:    texture.Width(),
		Height:   texture.Height(),
	}

	bg.RenderComponent = common.RenderComponent{
		Drawable:    texture,
		Scale:       engo.Point{X: engo.GameWidth() / texture.Width(), Y: engo.GameHeight() / texture.Height()},
		StartZIndex: -1,
	}

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&bg.BasicEntity, &bg.RenderComponent, &bg.SpaceComponent)
		}
	}
	return nil
}

type UpdateMsg struct{}

func (UpdateMsg) Type() string { return "UpdateMsg" }
