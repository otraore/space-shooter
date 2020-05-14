package systems

import (
	"log"
	"math/rand"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type RockSpawnSystem struct {
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

func (rock *RockSpawnSystem) New(w *ecs.World) {
	rock.world = w
	rock.SpawnRocks = false

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
}

func (*RockSpawnSystem) Remove(ecs.BasicEntity) {}

func (rock *RockSpawnSystem) Update(dt float32) {
	if rock.SpawnRocks {
		if rand.Float32() < .96 {
			return
		}

		position := engo.Point{
			X: rand.Float32() * engo.GameWidth(),
			Y: -32,
		}
		rock.NewRock(position)
	}
}

func (rs *RockSpawnSystem) NewRock(position engo.Point) {
	rock := Rock{BasicEntity: ecs.NewBasic()}
	rock.RenderComponent = common.RenderComponent{
		Drawable: rs.texture,
		Scale:    engo.Point{X: 1, Y: 1},
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
		case *FallingSystem:
			sys.Add(&rock.BasicEntity, &rock.SpaceComponent)
		}
	}
}

type ClearRocks struct{}

func (ClearRocks) Type() string {
	return "ClearRocks"
}

type SpawnRocks struct{}

func (SpawnRocks) Type() string {
	return "SpawnRocks"
}
