package gui

import (
	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
)

type Label struct {
	Font     *common.Font
	cache    string
	Position engo.Point
	World    *ecs.World
	Text     string

	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

func (l *Label) Init() {
	l.BasicEntity = ecs.NewBasic()

	width, height, _ := l.Font.TextDimensions(l.Text)

	l.SpaceComponent = common.SpaceComponent{
		Width:  float32(width),
		Height: float32(height),
	}
	l.SpaceComponent.Position = l.Position
	l.RenderComponent.Drawable = common.Text{
		Font: l.Font,
		Text: l.Text,
	}

	l.SetShader(common.TextHUDShader)

	for _, system := range l.World.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&l.BasicEntity, &l.RenderComponent, &l.SpaceComponent)
		}

	}
}

func (l *Label) SetText(s string) bool {
	if l.Font == nil {
		panic("Label.SetText called without setting Label.Font")
	}

	if l.cache == s {
		return false
	}

	if l.RenderComponent.Drawable == nil {
		l.RenderComponent.Drawable = common.Text{Font: l.Font}
	}

	fnt := l.RenderComponent.Drawable.(common.Text)
	fnt.Text = s

	return true
}
