package systems

import (
	"log"
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

type DeathSystem struct {
	Ship *Ship
}

func (s *DeathSystem) New(*ecs.World) {
	// Subscribe to ScoreMessage
	engo.Mailbox.Listen("CollisionMessage", func(message engo.Message) {
		_, isCollision := message.(common.CollisionMessage)

		if isCollision {
			log.Println("DEAD")
			engo.Mailbox.Dispatch(ClearRocks{})
			go func() {
				for i := 0; i < 5; i++ {
					time.Sleep(500 * time.Millisecond)
					s.Ship.RenderComponent.Hidden = !s.Ship.RenderComponent.Hidden
				}
				s.Ship.RenderComponent.Hidden = false
			}()
			s.Ship.SpaceComponent.Position.Y = engo.GameHeight() - s.Ship.SpaceComponent.Height
			s.Ship.SpaceComponent.Position.X = (engo.GameWidth() / 2) - s.Ship.SpaceComponent.Width
		}
	})
}

func (*DeathSystem) Remove(ecs.BasicEntity) {}
func (*DeathSystem) Update(float32)         {}
