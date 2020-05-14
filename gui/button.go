package gui

import (
	"fmt"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type Graphic struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	common.MouseComponent
}

type Button struct {
	Base
	Label        *Label
	Graphic      Graphic
	Image        *common.Texture
	ImageClicked *common.Texture
	Position     engo.Point
	World        *ecs.World
	Enabled      bool
	Text         string
	Font         *common.Font
	OnMouseOut   func(*Button)
}

func (b *Button) OnClick(f func()) {
	b.EventListeners[EventMouseClicked{}.Type()] = append(b.EventListeners[EventMouseClicked{}.Type()], f)
}

func (b *Button) OnMouseOver(f func()) {
	b.EventListeners[EventMouseOver{}.Type()] = append(b.EventListeners[EventMouseOver{}.Type()], f)
}

func (b *Button) Init() error {
	b.EventListeners = make(map[string][]func())

	b.Graphic.BasicEntity = ecs.NewBasic()

	b.Graphic.MouseComponent = common.MouseComponent{}

	b.Graphic.RenderComponent = common.RenderComponent{
		Drawable: b.Image,
		Scale:    engo.Point{X: 1, Y: 1}, //Todo: make this editable
	}

	b.Graphic.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: b.Position.X, Y: b.Position.Y},
		Width:    b.Image.Width(),
		Height:   b.Image.Height(),
	}
	width, height, _ := b.Font.TextDimensions(b.Text)

	b.Graphic.SpaceComponent.Position = b.Position

	b.World.AddSystem(&common.MouseSystem{})
	b.World.AddSystem(&ButtonSystem{})

	for _, system := range b.World.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&b.Graphic.BasicEntity, &b.Graphic.RenderComponent, &b.Graphic.SpaceComponent)
		case *common.MouseSystem:
			sys.Add(&b.Graphic.BasicEntity, &b.Graphic.MouseComponent, &b.Graphic.SpaceComponent, &b.Graphic.RenderComponent)
		case *ButtonSystem:
			sys.Add(b)
		}
	}

	b.Label = NewLabel(Label{
		World: b.World,
		Font:  b.Font,
		Text:  b.Text,
		Position: engo.Point{
			X: b.Graphic.SpaceComponent.Position.X + float32(((b.Graphic.SpaceComponent.Width - float32(width)) / 2)),
			Y: b.Graphic.SpaceComponent.Position.Y + float32(height/2),
		},
	})
	return nil
}

type buttonEntity struct {
	*Button
}

type ButtonSystem struct {
	entities []buttonEntity
}

func (c *ButtonSystem) New(w *ecs.World) {
	fmt.Println("Button Created")
}

func (c *ButtonSystem) Add(b *Button) {
	c.entities = append(c.entities, buttonEntity{b})
}

func (c *ButtonSystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, e := range c.entities {
		if e.Graphic.BasicEntity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		c.entities = append(c.entities[:delete], c.entities[delete+1:]...)
	}
}

func (c *ButtonSystem) Update(float32) {
	curPos := engo.Point{X: engo.Input.Mouse.X, Y: engo.Input.Mouse.Y}
	cursorHand := false
	for _, e := range c.entities {
		if e.Graphic.Contains(curPos) {
			e.Graphic.RenderComponent.Drawable = e.ImageClicked
			cursorHand = true
			e.DispatchEvents(EventMouseOver{})

			if e.Graphic.MouseComponent.Clicked {
				e.Graphic.RenderComponent.Drawable = e.ImageClicked
				e.DispatchEvents(EventMouseClicked{})
			}
		} else {
			e.Graphic.RenderComponent.Drawable = e.Image
		}
	}

	if cursorHand {
		engo.SetCursor(engo.CursorHand)
	} else {
		engo.SetCursor(engo.CursorNone)
	}
}
