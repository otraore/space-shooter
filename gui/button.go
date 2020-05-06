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
	Label          Label
	Graphic        Graphic
	Image          *common.Texture
	ImageClicked   *common.Texture
	Position       engo.Point
	World          *ecs.World
	Enabled        bool
	Text           string
	Font           *common.Font
	OnMouseOut     func(*Button)
	EventListeners map[string][]func()
}

func (b *Button) OnClick(f func()) {
	b.EventListeners["click"] = append(b.EventListeners["click"], f)
}

func (b *Button) OnMouseOver(f func()) {
	b.EventListeners["mouse_over"] = append(b.EventListeners["mouse_over"], f)
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
		Position: engo.Point{X: 0, Y: 0},
		Width:    b.Image.Width(),
		Height:   b.Image.Height(),
	}
	width, height, _ := b.Font.TextDimensions(b.Text)

	b.Graphic.SpaceComponent.Position = b.Position

	b.World.AddSystem(&common.MouseSystem{})
	b.World.AddSystem(&ButtonSystem{})

	for _, system := range b.World.Systems() {
		switch sys := system.(type) {
		case *common.MouseSystem:
			sys.Add(&b.Graphic.BasicEntity, &b.Graphic.MouseComponent, &b.Graphic.SpaceComponent, &b.Graphic.RenderComponent)
		case *ButtonSystem:
			sys.Add(b)
		case *common.RenderSystem:
			sys.Add(&b.Graphic.BasicEntity, &b.Graphic.RenderComponent, &b.Graphic.SpaceComponent)
		}

	}

	b.Label = Label{
		World: b.World,
		Font:  b.Font,
		Text:  b.Text,
		Position: engo.Point{
			X: b.Graphic.SpaceComponent.Position.X + float32(((b.Graphic.SpaceComponent.Width - float32(width)) / 2)),
			Y: b.Graphic.SpaceComponent.Position.Y + float32(height/2),
		},
	}
	b.Label.Init()
	return nil
}

type buttonEntity struct {
	*ecs.BasicEntity
	*common.SpaceComponent
	*common.MouseComponent
	*Button
}

type ButtonSystem struct {
	entities []buttonEntity
}

func (c *ButtonSystem) New(w *ecs.World) {
	fmt.Println("Button Created")
}

func (c *ButtonSystem) Add(b *Button) {
	c.entities = append(c.entities, buttonEntity{&b.Graphic.BasicEntity, &b.Graphic.SpaceComponent, &b.Graphic.MouseComponent, b})
}

func (c *ButtonSystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, e := range c.entities {
		if e.BasicEntity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		c.entities = append(c.entities[:delete], c.entities[delete+1:]...)
	}
}

func (c *ButtonSystem) Update(float32) {
	for _, e := range c.entities {
		pos := engo.Point{X: e.MouseX, Y: e.MouseY}
		if e.Contains(pos) {

			e.Graphic.RenderComponent.Drawable = e.ImageClicked
			engo.SetCursor(engo.CursorHand)
			for _, f := range e.EventListeners["mouse_over"] {
				f()
			}
			if e.MouseComponent.Clicked {
				e.Graphic.RenderComponent.Drawable = e.ImageClicked
				for _, f := range e.EventListeners["click"] {
					f()
				}
			}
		} else {
			engo.SetCursor(engo.CursorNone)
			e.Graphic.RenderComponent.Drawable = e.Image
		}

	}
}
