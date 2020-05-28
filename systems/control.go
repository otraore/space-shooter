package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/otraore/space-shooter/gui"
)

type controlEntity struct {
	*ecs.BasicEntity
	*common.SpaceComponent
}

type ControlSystem struct {
	entities  []controlEntity
	MenuScene engo.Scene
	Ship      *Ship
}

func (c *ControlSystem) Add(basic *ecs.BasicEntity, space *common.SpaceComponent) {
	c.entities = append(c.entities, controlEntity{basic, space})
}

func (c *ControlSystem) Remove(basic ecs.BasicEntity) {
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

func (c *ControlSystem) Update(dt float32) {
	speed := 500 * dt
	hori := engo.Input.Axis(engo.DefaultHorizontalAxis).Value()
	vert := engo.Input.Axis(engo.DefaultVerticalAxis).Value()

	// Out of bounds
	if c.Ship.SpaceComponent.Position.Y+c.Ship.SpaceComponent.Height > engo.GameHeight() {
		if engo.Input.Axis(engo.DefaultVerticalAxis).Value() > 0 {
			return
		}
	}

	if c.Ship.SpaceComponent.Position.X+c.Ship.SpaceComponent.Width > engo.GameWidth() {
		if engo.Input.Axis(engo.DefaultHorizontalAxis).Value() > 0 {
			return
		}
	}

	engo.Mailbox.Dispatch(gui.UpdateMsg{})

	for _, e := range c.entities {
		posX := e.SpaceComponent.Position.X + (speed * hori)
		posY := e.SpaceComponent.Position.Y + (speed * vert)

		if posX > 0 && posX < engo.GameWidth() {
			e.SpaceComponent.Position.X = posX
		}

		if posY > 0 && posY < engo.GameHeight() {
			e.SpaceComponent.Position.Y = posY
		}
	}

	if btn := engo.Input.Button("quit"); btn.JustPressed() {
		engo.Files.Unload("spritesheets/game.xml")
		engo.SetCursorVisibility(true)
		engo.SetScene(c.MenuScene, false)
	}
}
