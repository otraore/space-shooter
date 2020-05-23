package gui

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type Panel struct {
	Header                      Graphic
	Body                        Graphic
	Label                       *Label
	Text                        string
	Font                        *common.Font
	Position                    engo.Point
	Width, Height, HeaderHeight float32
	World                       *ecs.World
	HeaderImage                 *common.Texture
	BodyImage                   *common.Texture
}

func NewPanel(p Panel) (*Panel, error) {
	p.Header.BasicEntity = ecs.NewBasic()
	p.Body.BasicEntity = ecs.NewBasic()

	realPanelHeight := p.Height - p.HeaderHeight
	realHeaderHeight := p.HeaderHeight + 25

	if p.Width == 0 {
		p.Width = p.BodyImage.Width()
	}

	if p.Height == 0 {
		p.Height = p.BodyImage.Height()
	}

	p.Body.RenderComponent = common.RenderComponent{
		Drawable: p.BodyImage,
		Scale:    engo.Point{X: p.Width / p.BodyImage.Width(), Y: (realPanelHeight) / p.BodyImage.Height()},
	}
	p.Body.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: p.Position.X, Y: p.Position.Y + p.HeaderHeight - 15},
		Width:    p.Width,
		Height:   realPanelHeight,
	}

	p.Header.RenderComponent = common.RenderComponent{
		Drawable: p.HeaderImage,
		Scale:    engo.Point{X: p.Width / p.HeaderImage.Width(), Y: realHeaderHeight / p.HeaderImage.Height()},
	}
	p.Header.SpaceComponent = common.SpaceComponent{
		Position: p.Position,
		Width:    p.Width,
		Height:   realHeaderHeight,
	}

	for _, system := range p.World.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&p.Header.BasicEntity, &p.Header.RenderComponent, &p.Header.SpaceComponent)
			sys.Add(&p.Body.BasicEntity, &p.Body.RenderComponent, &p.Body.SpaceComponent)
		}
	}
	width, height, _ := p.Font.TextDimensions(p.Text)

	p.Label = NewLabel(Label{
		World: p.World,
		Font:  p.Font,
		Text:  p.Text,
		Position: engo.Point{
			X: p.Header.SpaceComponent.Position.X + float32(((p.Header.SpaceComponent.Width - float32(width)) / 2)),
			Y: p.Header.SpaceComponent.Position.Y + float32(height/2) + 5,
		},
	})
	return &p, nil
}
