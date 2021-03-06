package scenes

import (
	"image/color"
	"log"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/otraore/space-shooter/gui"
	"github.com/otraore/space-shooter/scenes/config"
	"github.com/otraore/space-shooter/systems"
)

var (
	ship systems.Ship
)

type GameScene struct{}

func (GameScene) Preload() {
	err := engo.Files.Load(gameSheet, "sounds/sfx_laser.mp3")
	if err != nil {
		log.Println(err)
	}
}

func (GameScene) Setup(u engo.Updater) {
	w, _ := u.(*ecs.World)

	fnt := &common.Font{
		URL:  uiFont,
		FG:   color.White,
		Size: 30,
	}

	err := fnt.CreatePreloaded()
	handleErr(err)

	// Create the ship object
	ship = systems.Ship{
		BasicEntity: ecs.NewBasic(),
		LivesLeft:   3,
		Score:       0,
		Config: config.ShipConfig{
			Color: config.ShipColorOrange,
			Type:  "2",
		},
		Font: fnt,
	}
	rockSys := &systems.RockSpawnSystem{SpawnRocks: true}
	projSys := &systems.ProjectileSystem{Ship: &ship, Color: "green"}
	// Add all of the systems
	w.AddSystem(&common.RenderSystem{})
	w.AddSystem(&common.CollisionSystem{})
	w.AddSystem(&common.AudioSystem{})
	w.AddSystem(projSys)
	w.AddSystem(rockSys)
	w.AddSystem(&systems.CollisionSystem{Ship: &ship.BasicEntity, ProjectileMaster: &projSys.Master, RockSys: rockSys, ProjSys: projSys})
	w.AddSystem(&systems.ShipSystem{Ship: &ship})
	w.AddSystem(&systems.ControlSystem{Ship: &ship, MenuScene: &MenuScene{}})

	engo.Input.RegisterButton("quit", engo.KeyQ, engo.KeyEscape)
	engo.Input.RegisterButton("fire", engo.KeySpace)

	texture, err := common.LoadedSprite(ship.AssetURL())
	handleErr(err)

	ship.RenderComponent = common.RenderComponent{
		Drawable:    texture,
		Scale:       config.GlobalScale,
		StartZIndex: 999999,
	}

	width := texture.Width() * ship.RenderComponent.Scale.X
	height := texture.Height() * ship.RenderComponent.Scale.Y
	ship.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{X: (engo.GameWidth() / 2) - width, Y: engo.GameHeight() - height},
		Width:    width,
		Height:   height,
	}

	player, err := common.LoadedPlayer("sounds/sfx_laser.mp3")
	if err != nil {
		handleErr(err)
	}
	ship.LaserPlayer = player
	ship.AudioComponent = common.AudioComponent{Player: player}

	texture, err = common.LoadedSprite("lives/" + ship.Config.Color.String() + ship.Config.Type + ".png")
	handleErr(err)

	lifeImg := gui.Image{
		World:    w,
		Texture:  texture,
		Scale:    engo.Point{X: 1, Y: 1},
		Position: engo.Point{X: 15, Y: 15},
	}
	lifeImg.Init()

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&ship.BasicEntity, &ship.RenderComponent, &ship.SpaceComponent)
		case *common.CollisionSystem:
			sys.Add(&ship.BasicEntity, &ship.CollisionComponent, &ship.SpaceComponent)
		case *common.AudioSystem:
			sys.Add(&ship.BasicEntity, &ship.AudioComponent)
		case *systems.ControlSystem:
			sys.Add(&ship.BasicEntity, &ship.SpaceComponent)
		case *systems.ShipSystem:
			sys.Add(&ship.BasicEntity, &ship.RenderComponent, &ship.SpaceComponent)
		}
	}

	err = gui.SetBackgroundImage(w, "backgrounds/blue.png")
	handleErr(err)

	log.Println("Start game, lives remaining: ", ship.LivesLeft)
	// Laser style 1 & 2
}

func (GameScene) Type() string { return "Game" }
