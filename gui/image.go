package gui

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type Image struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	Texture  *common.Texture
	Scale    engo.Point
	World    *ecs.World
	Position engo.Point
}

func (i *Image) Init() {
	i.BasicEntity = ecs.NewBasic()

	i.RenderComponent = common.RenderComponent{
		Drawable: i.Texture,
		Scale:    i.Scale,
	}

	i.SpaceComponent = common.SpaceComponent{
		Position: i.Position,
		Width:    i.Texture.Width(),
		Height:   i.Texture.Height(),
	}

	for _, system := range i.World.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&i.BasicEntity, &i.RenderComponent, &i.SpaceComponent)
		}

	}
}
