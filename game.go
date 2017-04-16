package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"

	"time"

	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
	"github.com/otraore/space-shooter/gui"
)

var guy Guy
var playing = true

type Guy struct {
	ecs.BasicEntity
	common.CollisionComponent
	common.RenderComponent
	common.SpaceComponent
}

type Rock struct {
	ecs.BasicEntity
	common.CollisionComponent
	common.RenderComponent
	common.SpaceComponent
}

type GameScene struct{}

func (GameScene) Preload() {
	err := engo.Files.Load("images/ui/playerLife3_red.png", "images/playerShip3_red.png", "images/rock.png", "fonts/kenvector_future.ttf")
	if err != nil {
		log.Println(err)
	}

	fmt.Println("Game Scene Preload")
}

func (GameScene) Setup(w *ecs.World) {
	fmt.Println("Game Scene Setup")

	common.SetBackground(color.Black)

	// Add all of the systems
	w.AddSystem(&common.RenderSystem{})
	w.AddSystem(&common.CollisionSystem{})
	w.AddSystem(&DeathSystem{})
	w.AddSystem(&FallingSystem{})
	w.AddSystem(&ControlSystem{})
	w.AddSystem(&RockSpawnSystem{})
	w.AddSystem(&GuySystem{})

	engo.Input.RegisterButton("quit", engo.Q, engo.Escape)

	texture, err := common.LoadedSprite("images/playerShip3_red.png")
	if err != nil {
		log.Println(err)
	}

	// Create an entity
	guy = Guy{BasicEntity: ecs.NewBasic()}

	// Initialize the components, set scale to 4x
	guy.RenderComponent = common.RenderComponent{
		Drawable: texture,
		Scale:    engo.Point{1, 1},
	}

	fmt.Println(engo.GameWidth())
	width := texture.Width() * guy.RenderComponent.Scale.X
	height := texture.Height() * guy.RenderComponent.Scale.Y
	guy.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{(engo.GameWidth() / 2) - width, engo.GameHeight() - height},
		Width:    width,
		Height:   height,
	}
	guy.CollisionComponent = common.CollisionComponent{
		Solid: true,
		Main:  true,
	}

	fnt := &common.Font{
		URL:  "fonts/kenvector_future.ttf",
		FG:   color.White,
		Size: 64,
	}

	err = fnt.CreatePreloaded()
	if err != nil {
		panic(err)
	}

	score := &gui.Label{
		World: w,
		Font:  fnt,
		Text:  "002600",
		Position: engo.Point{
			0,
			10,
		},
	}

	score.Init()

	score.SpaceComponent.Position.X = engo.GameWidth() - score.SpaceComponent.Width

	texture, err = common.LoadedSprite("images/ui/playerLife3_red.png")
	if err != nil {
		log.Println(err)
	}
	lifeImg := gui.Image{
		World:    w,
		Texture:  texture,
		Scale:    engo.Point{1, 1},
		Position: engo.Point{15, 15},
	}
	lifeImg.Init()

	lives := &gui.Label{
		World: w,
		Font:  fnt,
		Text:  "X 3",
		Position: engo.Point{
			60,
			10,
		},
	}

	lives.Init()
	// Add it to appropriate systems
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&guy.BasicEntity, &guy.RenderComponent, &guy.SpaceComponent)
		case *common.CollisionSystem:
			sys.Add(&guy.BasicEntity, &guy.CollisionComponent, &guy.SpaceComponent)
		case *ControlSystem:
			sys.Add(&guy.BasicEntity, &guy.SpaceComponent)
		case *GuySystem:
			sys.Add(&guy.BasicEntity, &guy.RenderComponent, &guy.SpaceComponent)
		}
	}
}

func (GameScene) Type() string { return "Game" }

type controlEntity struct {
	*ecs.BasicEntity
	*common.SpaceComponent
}

type ControlSystem struct {
	entities []controlEntity
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
	speed := 400 * dt
	// Out of bond
	if guy.SpaceComponent.Position.Y+guy.SpaceComponent.Height > engo.GameHeight() {
		if engo.Input.Axis(engo.DefaultVerticalAxis).Value() > 0 {
			return
		}
	}

	if guy.SpaceComponent.Position.X+guy.SpaceComponent.Width > engo.GameWidth() {
		if engo.Input.Axis(engo.DefaultHorizontalAxis).Value() > 0 {
			return
		}
	}

	for _, e := range c.entities {
		hori := engo.Input.Axis(engo.DefaultHorizontalAxis)
		e.SpaceComponent.Position.X += speed * hori.Value()

		vert := engo.Input.Axis(engo.DefaultVerticalAxis)
		e.SpaceComponent.Position.Y += speed * vert.Value()
	}

	if btn := engo.Input.Button("quit"); btn.JustPressed() {
		engo.SetScene(MenuScene{}, false)
	}
}

type RockSpawnSystem struct {
	world *ecs.World
}

func (rock *RockSpawnSystem) New(w *ecs.World) {
	rock.world = w
}

func (*RockSpawnSystem) Remove(ecs.BasicEntity) {}

func (rock *RockSpawnSystem) Update(dt float32) {
	if playing {
		if rand.Float32() < .96 {
			return
		}

		position := engo.Point{
			X: rand.Float32() * engo.GameWidth(),
			Y: -32,
		}
		NewRock(rock.world, position)
	}
}

func NewRock(world *ecs.World, position engo.Point) {
	texture, err := common.LoadedSprite("images/rock.png")
	if err != nil {
		log.Println(err)
	}

	rock := Rock{BasicEntity: ecs.NewBasic()}
	rock.RenderComponent = common.RenderComponent{
		Drawable: texture,
		Scale:    engo.Point{4, 4},
	}
	rock.SpaceComponent = common.SpaceComponent{
		Position: position,
		Width:    texture.Width() * rock.RenderComponent.Scale.X,
		Height:   texture.Height() * rock.RenderComponent.Scale.Y,
	}
	rock.CollisionComponent = common.CollisionComponent{Solid: true}

	for _, system := range world.Systems() {
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

type fallingEntity struct {
	*ecs.BasicEntity
	*common.SpaceComponent
}

type FallingSystem struct {
	entities []fallingEntity
}

type ClearRocks struct{}

func (ClearRocks) Type() string {
	return "ClearRocks"
}

func (f *FallingSystem) New(w *ecs.World) {
	engo.Mailbox.Listen("ClearRocks", func(message engo.Message) {
		for _, e := range f.entities {
			w.RemoveEntity(*e.BasicEntity)
		}
		playing = false
		go func() {
			time.Sleep(5 * time.Second)
			playing = true
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
	}
}

type DeathSystem struct{}

func (*DeathSystem) New(*ecs.World) {
	// Subscribe to ScoreMessage
	engo.Mailbox.Listen("CollisionMessage", func(message engo.Message) {
		_, isCollision := message.(common.CollisionMessage)

		if isCollision {
			log.Println("DEAD")
			engo.Mailbox.Dispatch(ClearRocks{})
			go func() {
				for i := 0; i < 5; i++ {
					time.Sleep(500 * time.Millisecond)
					guy.RenderComponent.Hidden = !guy.RenderComponent.Hidden
				}
				guy.RenderComponent.Hidden = false
			}()
			guy.SpaceComponent.Position.Y = engo.GameHeight() - guy.SpaceComponent.Height
			guy.SpaceComponent.Position.X = (engo.GameWidth() / 2) - guy.SpaceComponent.Width
		}
	})
}

type guyEntity struct {
	*ecs.BasicEntity
	*common.RenderComponent
	*common.SpaceComponent
}

type GuySystem struct {
	entities []guyEntity
}

func (g *GuySystem) New(w *ecs.World) {
	fmt.Println(engo.GameWidth())
	guy.SpaceComponent.Position.Y = engo.GameHeight() - guy.SpaceComponent.Height
	guy.SpaceComponent.Position.X = (engo.GameWidth() / 2) - guy.SpaceComponent.Width
}

func (s *GuySystem) Add(basic *ecs.BasicEntity, render *common.RenderComponent, space *common.SpaceComponent) {
	s.entities = append(s.entities, guyEntity{basic, render, space})
}

func (s *GuySystem) Remove(basic ecs.BasicEntity) {
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

func (s *GuySystem) Update(dt float32) {

}

func (*DeathSystem) Remove(ecs.BasicEntity) {}
func (*DeathSystem) Update(float32)         {}
