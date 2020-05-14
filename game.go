package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/otraore/space-shooter/gui"
	"github.com/otraore/space-shooter/systems"
)

var (
	ship      systems.Ship
	playing   = true
	livesLeft = 3
)

type GameScene struct{}

func (GameScene) Preload() {
	err := engo.Files.Load("spritesheets/game.xml", "images/rock.png", "images/playerLife3_red.png")
	if err != nil {
		log.Println(err)
	}

	fmt.Println("Game Scene Preload")
}

func (GameScene) Setup(u engo.Updater) {
	w, _ := u.(*ecs.World)

	fnt := &common.Font{
		URL:  "fonts/kenvector_future.ttf",
		FG:   color.White,
		Size: 30,
	}

	err := fnt.CreatePreloaded()
	if err != nil {
		panic(err)
	}

	// Create the ship object
	ship = systems.Ship{
		BasicEntity: ecs.NewBasic(),
		LivesLeft:   10,
		Color:       "red",
		Type:        "1",
		Font:        fnt,
	}
	// Add all of the systems
	w.AddSystem(&common.RenderSystem{})
	w.AddSystem(&common.CollisionSystem{})
	w.AddSystem(&systems.DeathSystem{Ship: &ship})
	w.AddSystem(&systems.FallingSystem{})
	w.AddSystem(&systems.ControlSystem{Ship: &ship, MenuScene: &MenuScene{}})
	w.AddSystem(&systems.RockSpawnSystem{SpawnRocks: false})
	w.AddSystem(&systems.ShipSystem{Ship: &ship})

	engo.Input.RegisterButton("quit", engo.KeyQ, engo.KeyEscape)

	texture, err := common.LoadedSprite(ship.AssetURL())
	if err != nil {
		log.Println(err)
	}

	ship.RenderComponent = common.RenderComponent{
		Drawable: texture,
		Scale:    engo.Point{X: 1, Y: 1},
	}

	width := texture.Width() * ship.RenderComponent.Scale.X
	height := texture.Height() * ship.RenderComponent.Scale.Y
	ship.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: (engo.GameWidth() / 2) - width, Y: engo.GameHeight() - height},
		Width:    width,
		Height:   height,
	}
	ship.CollisionComponent = common.CollisionComponent{
		Main:  1,
		Group: 1,
	}

	score := gui.NewLabel(gui.Label{
		World: w,
		Font:  fnt,
		Text:  "002600",
		Position: engo.Point{
			X: 0,
			Y: 10,
		},
	})

	score.SpaceComponent.Position.X = engo.GameWidth() - score.SpaceComponent.Width

	texture, err = common.LoadedSprite("images/playerLife3_red.png")
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

	// Add it to appropriate systems
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&ship.BasicEntity, &ship.RenderComponent, &ship.SpaceComponent)
		case *common.CollisionSystem:
			sys.Add(&ship.BasicEntity, &ship.CollisionComponent, &ship.SpaceComponent)
		case *systems.ControlSystem:
			sys.Add(&ship.BasicEntity, &ship.SpaceComponent)
		case *systems.ShipSystem:
			sys.Add(&ship.BasicEntity, &ship.RenderComponent, &ship.SpaceComponent)
		}
	}

	err = gui.SetBackgroundImage(w, "backgrounds/blue.png")
	if err != nil {
		log.Println(err)
	}

	engo.Mailbox.Dispatch(systems.SpawnRocks{})
	log.Println("Start game, lives remaining: ", ship.LivesLeft)
}

func (GameScene) Type() string { return "Game" }
