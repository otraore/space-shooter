package gui

import (
	"errors"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
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

func NewLabel(l Label) *Label {
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

	return &l
}

func (l *Label) SetText(s string) error {
	if l.Font == nil {
		return errors.New("Label.SetText called without setting Label.Font")
	}
	if l.cache == s {
		return nil
	}

	l.Text = s
	l.cache = s

	l.RenderComponent.Drawable = common.Text{
		Font: l.Font,
		Text: s,
	}

	width, height, _ := l.Font.TextDimensions(l.Text)

	l.SpaceComponent.Width = float32(width)
	l.SpaceComponent.Height = float32(height)

	return nil
}
