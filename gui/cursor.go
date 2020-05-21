package gui

import (
	"log"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

var cursorSystemAdded = false

func SetCursorImage(w *ecs.World, imageUrl string) error {
	cursor := basic{BasicEntity: ecs.NewBasic()}

	texture, err := common.LoadedSprite(imageUrl)
	if err != nil {
		log.Printf("Unable to load cursor image: %v \n", err)
	}

	if !cursorSystemAdded {
		w.AddSystem(&CursorSystem{})
	}

	cursor.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{0, 0},
		Width:    texture.Width() * cursor.RenderComponent.Scale.X,
		Height:   texture.Height() * cursor.RenderComponent.Scale.Y,
	}

	cursor.RenderComponent = common.RenderComponent{
		Drawable:    texture,
		Scale:       engo.Point{1, 1},
		StartZIndex: 99999,
	}

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&cursor.BasicEntity, &cursor.RenderComponent, &cursor.SpaceComponent)
		case *CursorSystem:
			sys.Add(&cursor.BasicEntity, &cursor.RenderComponent, &cursor.SpaceComponent)
		}
	}
	return nil
}

type cursorEntity struct {
	*ecs.BasicEntity
	*common.RenderComponent
	*common.SpaceComponent
}

type CursorSystem struct {
	entities []cursorEntity
}

func (s *CursorSystem) Add(basic *ecs.BasicEntity, render *common.RenderComponent, space *common.SpaceComponent) {
	s.entities = append(s.entities, cursorEntity{basic, render, space})
}

func (s *CursorSystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, e := range s.entities {
		if e.BasicEntity.ID() == basic.ID() {
			delete = index
			break
		}
	}

	if delete >= 0 {
		s.entities = append(s.entities[:delete], s.entities[delete+1:]...)
	}
}

func (s *CursorSystem) Update(dt float32) {
	for _, e := range s.entities {
		e.SpaceComponent.Position.X += engo.Input.Axis(engo.DefaultMouseXAxis).Value()
		e.SpaceComponent.Position.Y += engo.Input.Axis(engo.DefaultMouseYAxis).Value()
	}
}
