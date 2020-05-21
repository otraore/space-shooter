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
	world    *ecs.World
}

func (f *FallingSystem) New(w *ecs.World) {
	f.world = w
	engo.Mailbox.Listen(ClearRocks{}.Type(), func(message engo.Message) {
		for _, e := range f.entities {
			f.Remove(*e.BasicEntity)
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
	for i, e := range f.entities {
		if e.BasicEntity.ID() == basic.ID() {
			for _, system := range f.world.Systems() {
				switch system.(type) {
				case *FallingSystem:
				case *common.RenderSystem:
				default:
					system.Remove(*e.BasicEntity)
				}
			}
			f.entities = append(f.entities[:i], f.entities[i+1:]...)
			break
		}
	}
}

func (f *FallingSystem) Update(dt float32) {
	speed := 400 * dt

	// log.Println("num rocks:", len(f.entities))
	for _, e := range f.entities {
		e.SpaceComponent.Position.Y += speed

		if e.SpaceComponent.Position.Y > engo.GameHeight() {
			f.Remove(*e.BasicEntity)
		}
	}
}
