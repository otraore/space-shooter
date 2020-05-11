package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"strconv"

	"time"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/otraore/space-shooter/gui"
)

var (
	guy        Guy
	playing    = true
	livesLeft  = 3
	livesLabel *gui.Label
)

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
	err := engo.Files.Load("spritesheets/sheet.xml", "images/rock.png", "fonts/kenvector_future.ttf")
	if err != nil {
		log.Println(err)
	}

	fmt.Println("Game Scene Preload")
}

func (GameScene) Setup(u engo.Updater) {
	w, _ := u.(*ecs.World)

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

	engo.Input.RegisterButton("quit", engo.KeyQ, engo.KeyEscape)

	texture, err := common.LoadedSprite("playerShip3_orange.png")
	if err != nil {
		log.Println(err)
	}

	// Create an entity
	guy = Guy{BasicEntity: ecs.NewBasic()}

	// Initialize the components, set scale to 4x
	guy.RenderComponent = common.RenderComponent{
		Drawable: texture,
		Scale:    engo.Point{X: 1, Y: 1},
	}

	width := texture.Width() * guy.RenderComponent.Scale.X
	height := texture.Height() * guy.RenderComponent.Scale.Y
	guy.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: (engo.GameWidth() / 2) - width, Y: engo.GameHeight() - height},
		Width:    width,
		Height:   height,
	}
	guy.CollisionComponent = common.CollisionComponent{
		Main:  1,
		Group: 1,
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
			X: 0,
			Y: 10,
		},
	}

	score.Init()

	score.SpaceComponent.Position.X = engo.GameWidth() - score.SpaceComponent.Width

	texture, err = common.LoadedSprite("playerLife3_red.png")
	if err != nil {
		log.Println(err)
	}
	lifeImg := gui.Image{
		World:    w,
		Texture:  texture,
		Scale:    engo.Point{X: 1, Y: 1},
		Position: engo.Point{X: 15, Y: 15},
	}
	lifeImg.Init()

	livesLabel = &gui.Label{
		World: w,
		Font:  fnt,
		Text:  "X " + strconv.Itoa(livesLeft),
		Position: engo.Point{
			X: 60,
			Y: 10,
		},
	}

	livesLabel.Init()
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
	world   *ecs.World
	texture *common.Texture
}

func (rock *RockSpawnSystem) New(w *ecs.World) {
	rock.world = w
	texture, err := common.LoadedSprite("images/rock.png")
	if err != nil {
		log.Println(err)
	}

	rock.texture = texture
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
		rock.NewRock(position)
	}
}

func (rs *RockSpawnSystem) NewRock(position engo.Point) {
	rock := Rock{BasicEntity: ecs.NewBasic()}
	rock.RenderComponent = common.RenderComponent{
		Drawable: rs.texture,
		Scale:    engo.Point{X: 4, Y: 4},
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
		livesLeft--
		livesLabel.SetText("X " + strconv.Itoa(livesLeft))
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

		if e.SpaceComponent.Position.Y > engo.GameHeight() {
			f.Remove(*e.BasicEntity)
		}
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
