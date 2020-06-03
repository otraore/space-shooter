package systems

import (
	"log"
	"math/rand"
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/otraore/space-shooter/config"
)

type RockSpawnSystem struct {
	entities   []rockEntity
	world      *ecs.World
	texture    *common.Texture
	SpawnRocks bool
}

type Rock struct {
	ecs.BasicEntity
	common.CollisionComponent
	common.RenderComponent
	common.SpaceComponent
}

type rockEntity struct {
	*ecs.BasicEntity
	*common.SpaceComponent
}

func (rock *RockSpawnSystem) New(w *ecs.World) {
	rock.world = w

	texture, err := common.LoadedSprite("rocks/meteorBrown_big1.png")
	if err != nil {
		log.Println(err)
	}

	rock.texture = texture

	engo.Mailbox.Listen(ClearRocks{}.Type(), func(_ engo.Message) {
		rock.SpawnRocks = false
	})

	engo.Mailbox.Listen(SpawnRocks{}.Type(), func(_ engo.Message) {
		rock.SpawnRocks = true
	})

	engo.Mailbox.Listen(ClearRocks{}.Type(), func(message engo.Message) {
		for _, e := range rock.entities {
			rock.Remove(*e.BasicEntity)
			w.RemoveEntity(*e.BasicEntity)
		}

		go func() {
			time.Sleep(5 * time.Second)
			engo.Mailbox.Dispatch(SpawnRocks{})
		}()
	})
}

func (r *RockSpawnSystem) Add(basic *ecs.BasicEntity, space *common.SpaceComponent) {
	r.entities = append(r.entities, rockEntity{basic, space})
}

func (r *RockSpawnSystem) Remove(basic ecs.BasicEntity) {
	for i, e := range r.entities {
		if e.BasicEntity.ID() == basic.ID() {
			for _, system := range r.world.Systems() {
				switch system.(type) {
				case *common.CollisionSystem:
					system.Remove(*e.BasicEntity)
				case *common.RenderSystem:
					system.Remove(*e.BasicEntity)
				}
			}
			r.entities = append(r.entities[:i], r.entities[i+1:]...)
			break
		}
	}
}

func (r *RockSpawnSystem) Update(dt float32) {
	// Rock span logic
	if r.SpawnRocks {
		if rand.Float32() > .96 {
			position := engo.Point{
				X: rand.Float32() * engo.GameWidth(),
				Y: -32,
			}
			r.NewRock(position)
		}
	}

	// Falling logic
	for _, e := range r.entities {
		e.SpaceComponent.Position.Y += 5.5

		if e.SpaceComponent.Position.Y > engo.GameHeight() {
			r.Remove(*e.BasicEntity)
		}
	}
}

func (rs *RockSpawnSystem) NewRock(position engo.Point) {
	rock := Rock{BasicEntity: ecs.NewBasic()}
	rock.RenderComponent = common.RenderComponent{
		Drawable: rs.texture,
		Scale:    config.GlobalScale,
	}
	rock.SpaceComponent = common.SpaceComponent{
		Position: position,
		Width:    rs.texture.Width() * rock.RenderComponent.Scale.X,
		Height:   rs.texture.Height() * rock.RenderComponent.Scale.Y,
	}
	rock.CollisionComponent = common.CollisionComponent{Group: 1}

	for _, system := range rs.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&rock.BasicEntity, &rock.RenderComponent, &rock.SpaceComponent)
		case *common.CollisionSystem:
			sys.Add(&rock.BasicEntity, &rock.CollisionComponent, &rock.SpaceComponent)
		}
	}

	rs.Add(&rock.BasicEntity, &rock.SpaceComponent)
}

type ClearRocks struct{}

func (ClearRocks) Type() string {
	return "ClearRocks"
}

type SpawnRocks struct{}

func (SpawnRocks) Type() string {
	return "SpawnRocks"
}
