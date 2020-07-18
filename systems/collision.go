package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/otraore/space-shooter/scenes/config"
)

type CollisionSystem struct {
	Ship             *ecs.BasicEntity
	ProjectileMaster *ecs.BasicEntity
	RockSys          *RockSpawnSystem
	ProjSys          *ProjectileSystem
}

func (c *CollisionSystem) New(*ecs.World) {
	// Subscribe to ScoreMessage
	engo.Mailbox.Listen("CollisionMessage", func(message engo.Message) {
		msg, isCollision := message.(common.CollisionMessage)
		if !isCollision {
			return
		}

		if msg.Entity.ID() == c.Ship.ID() || msg.To.ID() == c.Ship.ID() {
			// c.Ship.OnCollision()
		} else {
			if msg.Entity.Parent() != nil && msg.To.Parent() != nil {
				if msg.Entity.Parent().ID() == c.ProjectileMaster.ID() && msg.To.Parent().ID() == c.ProjectileMaster.ID() {
					// Both of the entities are projectiles
					return
				}
			}

			// projectile has hit a rock
			c.RockSys.Remove(*msg.To.BasicEntity)
			c.ProjSys.Remove(*msg.Entity.BasicEntity)
			engo.Mailbox.Dispatch(config.ScoreChanged{})
		}
	})
}

func (*CollisionSystem) Remove(ecs.BasicEntity) {}
func (*CollisionSystem) Update(float32)         {}
