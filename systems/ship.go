package systems

import (
	"strconv"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/otraore/space-shooter/gui"
)

type Ship struct {
	LivesLeft int
	// Color of the ship (red, blue, orange)
	Color string
	// Type of the ship (1, 2, 3)
	Type string
	// Font is the font used to display the lives and score for the ship
	Font *common.Font
	ecs.BasicEntity
	common.CollisionComponent
	common.RenderComponent
	common.SpaceComponent
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
}

func (s *ShipSystem) New(w *ecs.World) {
	s.Ship.SpaceComponent.Position.Y = engo.GameHeight() - s.Ship.SpaceComponent.Height
	s.Ship.SpaceComponent.Position.X = (engo.GameWidth() / 2) - s.Ship.SpaceComponent.Width

	s.LivesLabel = gui.NewLabel(gui.Label{
		World: w,
		Font:  s.Ship.Font,
		Text:  "X " + strconv.Itoa(s.Ship.LivesLeft),
		Position: engo.Point{
			X: 60,
			Y: 10,
		},
	})

	engo.Mailbox.Listen(ClearRocks{}.Type(), func(_ engo.Message) {
		s.Ship.LivesLeft--
		s.LivesLabel.SetText("X " + strconv.Itoa(s.Ship.LivesLeft))
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
	return "ships/" + s.Color + s.Type + ".png"
}
