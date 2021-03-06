package systems

import (
	"strconv"
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/otraore/space-shooter/gui"
	"github.com/otraore/space-shooter/scenes/config"
)

type Ship struct {
	LivesLeft, Score int
	// Config are the settings that change the look and functionality of the ship
	Config config.ShipConfig
	// Font is the font used to display the lives and score for the ship
	Font        *common.Font
	LaserPlayer *common.Player
	ecs.BasicEntity
	common.CollisionComponent
	common.RenderComponent
	common.SpaceComponent
	common.AudioComponent
}

type shipEntity struct {
	*ecs.BasicEntity
	*common.RenderComponent
	*common.SpaceComponent
}

type ShipSystem struct {
	entities   []shipEntity
	Ship       *Ship
	LivesLabel *gui.Label
	ScoreLabel *gui.Label
}

func (s *ShipSystem) New(w *ecs.World) {
	s.Ship.SpaceComponent.Position.Y = engo.GameHeight() - s.Ship.SpaceComponent.Height
	s.Ship.SpaceComponent.Position.X = (engo.GameWidth() / 2) - s.Ship.SpaceComponent.Width

	s.Ship.CollisionComponent = common.CollisionComponent{
		Main: 1,
	}

	s.LivesLabel = gui.NewLabel(gui.Label{
		World: w,
		Font:  s.Ship.Font,
		Text:  "X " + strconv.Itoa(s.Ship.LivesLeft),
		Position: engo.Point{
			X: 60,
			Y: 10,
		},
	})

	s.ScoreLabel = gui.NewLabel(gui.Label{
		World: w,
		Font:  s.Ship.Font,
		Text:  "0",
		Position: engo.Point{
			X: 0,
			Y: 10,
		},
	})

	s.ScoreLabel.SpaceComponent.Position.X = engo.GameWidth() - s.ScoreLabel.SpaceComponent.Width - 7.5
	engo.Mailbox.Listen(ClearRocks{}.Type(), func(_ engo.Message) {
		s.Ship.LivesLeft--
		s.LivesLabel.SetText("X " + strconv.Itoa(s.Ship.LivesLeft))
	})

	engo.Mailbox.Listen(config.ScoreChanged{}.Type(), func(_ engo.Message) {
		s.Ship.Score += 20
		s.ScoreLabel.SetText(strconv.Itoa(s.Ship.Score))
		s.ScoreLabel.SpaceComponent.Position.X = engo.GameWidth() - s.ScoreLabel.SpaceComponent.Width - 7.5
	})
}

func (s *ShipSystem) Add(basic *ecs.BasicEntity, render *common.RenderComponent, space *common.SpaceComponent) {
	s.entities = append(s.entities, shipEntity{basic, render, space})
}

func (s *ShipSystem) Remove(basic ecs.BasicEntity) {
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

func (s *ShipSystem) Update(dt float32) {

}

// AssetURL returns the location of the ship image based on it's type and color
func (s Ship) AssetURL() string {
	return "ships/" + s.Config.Color.String() + s.Config.Type + ".png"
}

// ResetPos puts the ship back to it's starting position
func (s *Ship) ResetPos() {
	s.SpaceComponent.Position = engo.Point{X: (engo.GameWidth() / 2) - s.SpaceComponent.Width, Y: engo.GameHeight() - s.SpaceComponent.Height}
}

func (s *Ship) OnCollision() {
	engo.Mailbox.Dispatch(ClearRocks{})
	s.ResetPos()

	go func() {
		for i := 0; i < 5; i++ {
			time.Sleep(500 * time.Millisecond)
			s.RenderComponent.Hidden = !s.RenderComponent.Hidden
		}
		s.RenderComponent.Hidden = false
	}()
}
