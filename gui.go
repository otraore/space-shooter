package main

import (
	"fmt"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
)

type Label struct {
	Font  *common.Font
	cache string

	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
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

type Graphic struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	common.MouseComponent
}

type Button struct {
	Label        Label
	Graphic      Graphic
	Image        *common.Texture
	ImageClicked *common.Texture

	OnClick     func(*Button)
	OnMouseOver func(*Button)
	OnMouseOut  func(*Button)
}

func NewButton(w *ecs.World, bg *common.Texture, bgClicked *common.Texture, f *common.Font, label string) *Button {
	b := new(Button)
	b.Label.BasicEntity = ecs.NewBasic()
	b.Graphic.BasicEntity = ecs.NewBasic()
	b.Image = bg
	b.ImageClicked = bgClicked

	b.Graphic.MouseComponent = common.MouseComponent{}

	b.Graphic.RenderComponent = common.RenderComponent{
		Drawable: bg,
		Scale:    engo.Point{1, 1},
	}

	width, height, _ := f.TextDimensions(label)
	b.Graphic.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{0, 0},
		Width:    bg.Width(),
		Height:   bg.Height(),
	}

	b.Label.SpaceComponent = common.SpaceComponent{
		Width:  float32(width),
		Height: float32(height),
	}

	b.Graphic.SpaceComponent.Position.X = (engo.GameWidth() / 2) - b.Graphic.SpaceComponent.Width/2
	b.Graphic.SpaceComponent.Position.Y = (engo.GameHeight() / 2) - b.Graphic.SpaceComponent.Height/2

	b.Label.SpaceComponent.Position.X = b.Graphic.SpaceComponent.Position.X + float32(width)/2
	b.Label.SpaceComponent.Position.Y = b.Graphic.SpaceComponent.Position.Y + float32(height)

	b.Label.RenderComponent.Drawable = common.Text{
		Font: f,
		Text: label,
	}

	b.Label.SetShader(common.TextHUDShader)
	w.AddSystem(&common.MouseSystem{})
	w.AddSystem(&ButtonSystem{})

	fmt.Println(b.Label.SpaceComponent.Width)
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.MouseSystem:
			sys.Add(&b.Graphic.BasicEntity, &b.Graphic.MouseComponent, &b.Graphic.SpaceComponent, &b.Graphic.RenderComponent)
		case *ButtonSystem:
			sys.Add(b)
		}
	}
	return b
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
	fmt.Println("Created")
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
		pos := engo.Point{e.MouseX, e.MouseY}
		if e.Contains(pos) {
			e.Graphic.RenderComponent.Drawable = e.ImageClicked
			engo.SetCursor(engo.CursorHand)
			if e.MouseComponent.Clicked {
				e.Graphic.RenderComponent.Drawable = e.ImageClicked
				engo.SetScene(GameScene{}, true)
				fmt.Println("clicked")
			}
		} else {
			engo.SetCursor(engo.CursorNone)
			e.Graphic.RenderComponent.Drawable = e.Image
		}

	}
}
