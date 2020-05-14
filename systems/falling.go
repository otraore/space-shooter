package systems

import (
	"fmt"
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type fallingEntity struct {
	*ecs.BasicEntity
	*common.SpaceComponent
}

type FallingSystem struct {
	entities []fallingEntity
}

func (f *FallingSystem) New(w *ecs.World) {
	engo.Mailbox.Listen(ClearRocks{}.Type(), func(message engo.Message) {
		for _, e := range f.entities {
			w.RemoveEntity(*e.BasicEntity)
		}

		go func() {
			time.Sleep(5 * time.Second)
			engo.Mailbox.Dispatch(SpawnRocks{})
		}()
		fmt.Println("ClearRocks")
	})
}

func (f *FallingSystem) Add(basic *ecs.BasicEntity, space *common.SpaceComponent) {
	f.entities = append(f.entities, fallingEntity{basic, space})
}

func (f *FallingSystem) Remove(basic ecs.BasicEntity) {
	delete := -1
	for index, e := range f.entities {
		if e.BasicEntity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		f.entities = append(f.entities[:delete], f.entities[delete+1:]...)
	}
}

func (f *FallingSystem) Update(dt float32) {
	speed := 400 * dt

	for _, e := range f.entities {
		e.SpaceComponent.Position.Y += speed

		if e.SpaceComponent.Position.Y > engo.GameHeight() {
			f.Remove(*e.BasicEntity)
		}
	}
}
