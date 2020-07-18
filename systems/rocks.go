package systems

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/otraore/space-shooter/config"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

type RockSpawnSystem struct {
	entities   []rockEntity
	world      *ecs.World
	textures   map[string]*common.Texture
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

func (rs *RockSpawnSystem) New(w *ecs.World) {
	rs.world = w
	rs.textures = make(map[string]*common.Texture)

	for _, color := range config.RockColors {
		for _, size := range config.RockSizes {
			count := 2
			if size == config.RockSizeBig {
				count = 4
			}

			for i := 1; i <= count; i++ {
				url := rs.rockURL(size, color, config.RockType(i))
				texture, err := common.LoadedSprite(url)
				if err != nil {
					log.Println(err)
				}
				rs.textures[url] = texture
			}
		}
	}

	engo.Mailbox.Listen(ClearRocks{}.Type(), func(_ engo.Message) {
		rs.SpawnRocks = false
	})

	engo.Mailbox.Listen(SpawnRocks{}.Type(), func(_ engo.Message) {
		rs.SpawnRocks = true
	})

	engo.Mailbox.Listen(ClearRocks{}.Type(), func(message engo.Message) {
		for _, e := range rs.entities {
			rs.Remove(*e.BasicEntity)
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
		if rand.Float32() > .95 {
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

	color := config.RockColors[rand.Intn(len(config.RockColors))]
	size := config.RockSizes[rand.Intn(len(config.RockSizes))]
	numTypes := 2
	if size == config.RockSizeBig {
		numTypes = 4
	}
	rtype := config.RockTypes[rand.Intn(numTypes)]

	textureURL := rs.rockURL(size, color, rtype)
	texture := rs.textures[textureURL]
	rock.RenderComponent = common.RenderComponent{
		Drawable: texture,
		Scale:    config.GlobalScale,
	}
	rock.SpaceComponent = common.SpaceComponent{
		Position: position,
		Width:    texture.Width() * rock.RenderComponent.Scale.X,
		Height:   texture.Height() * rock.RenderComponent.Scale.Y,
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

func (RockSpawnSystem) rockURL(size config.RockSize, color config.RockColor, num config.RockType) string {
	return fmt.Sprintf("rocks/%s_%s_%d.png", size.String(), color.String(), num)
}

type ClearRocks struct{}

func (ClearRocks) Type() string {
	return "ClearRocks"
}

type SpawnRocks struct{}

func (SpawnRocks) Type() string {
	return "SpawnRocks"
}
