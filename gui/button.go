package gui

import (
	"fmt"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

var btnSystemsAdded = false

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

func NewButton(b Button) (*Button, error) {
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

	b.Graphic.SpaceComponent.Position = b.Position

	// Make sure only one instance of the systems are added
	if !btnSystemsAdded {
		b.World.AddSystem(&ButtonSystem{})
		btnSystemsAdded = true
	}

	for _, system := range b.World.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&b.Graphic.BasicEntity, &b.Graphic.RenderComponent, &b.Graphic.SpaceComponent)
		case *ButtonSystem:
			sys.Add(&b)
		}
	}

	width, height, _ := b.Font.TextDimensions(b.Text)

	b.Label = NewLabel(Label{
		World: b.World,
		Font:  b.Font,
		Text:  b.Text,
		Position: engo.Point{
			X: b.Graphic.SpaceComponent.Position.X + float32(((b.Graphic.SpaceComponent.Width - float32(width)) / 2)),
			Y: b.Graphic.SpaceComponent.Position.Y + float32(height/2),
		},
	})
	return &b, nil
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
	btnHovered := false
	curPos := engo.Point{X: engo.Input.Mouse.X, Y: engo.Input.Mouse.Y}
	for _, e := range c.entities {
		if e.Graphic.Contains(curPos) {
			e.Graphic.RenderComponent.Drawable = e.ImageClicked
			btnHovered = true
			e.DispatchEvents(EventMouseOver{})

			if engo.Input.Mouse.Action == engo.Press && engo.Input.Mouse.Button == engo.MouseButtonLeft {
				e.Graphic.RenderComponent.Drawable = e.Image
				e.DispatchEvents(EventMouseClicked{})
			}
		} else {
			e.Graphic.RenderComponent.Drawable = e.Image
		}
	}

	if btnHovered {
		engo.SetCursor(engo.CursorHand)
	} else {
		engo.SetCursor(engo.CursorNone)
	}
}
