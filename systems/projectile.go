package systems

import (
	"log"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	"github.com/otraore/space-shooter/config"
)

type Projectile struct {
	ecs.BasicEntity
	common.CollisionComponent
	common.RenderComponent
	common.SpaceComponent
}

type projectileEntity struct {
	*ecs.BasicEntity
	*common.SpaceComponent
}

type ProjectileSystem struct {
	world    *ecs.World
	Master   ecs.BasicEntity
	entities []projectileEntity
	texture  *common.Texture
	Color    string
	Ship     *Ship
}

func (p *ProjectileSystem) New(w *ecs.World) {
	if p.Color == "" {
		p.Color = "blue"
	}

	texture, err := common.LoadedSprite("lasers/" + p.Color + ".png")
	if err != nil {
		log.Println(err)
	}

	p.texture = texture
	p.world = w
	p.Master = ecs.NewBasic()
}

func (p *ProjectileSystem) Add(basic *ecs.BasicEntity, space *common.SpaceComponent) {
	p.entities = append(p.entities, projectileEntity{basic, space})
}

func (p *ProjectileSystem) Remove(basic ecs.BasicEntity) {
	for i, e := range p.entities {
		if e.BasicEntity.ID() == basic.ID() {
			for _, system := range p.world.Systems() {
				switch system.(type) {
				case *common.CollisionSystem:
					system.Remove(*e.BasicEntity)
				case *common.RenderSystem:
					system.Remove(*e.BasicEntity)
				}
			}

			p.entities = append(p.entities[:i], p.entities[i+1:]...)
			break
		}
	}
}

func (p *ProjectileSystem) Update(dt float32) {
	if btn := engo.Input.Button("fire"); btn.JustPressed() {
		posX := p.Ship.SpaceComponent.Position.X + (p.Ship.Width / 2) - p.texture.Width() + 5
		p.Fire(engo.Point{X: posX, Y: p.Ship.SpaceComponent.Position.Y})
	}

	const speed = 12

	for _, e := range p.entities {
		e.SpaceComponent.Position.Y -= speed
		if e.SpaceComponent.Position.Y < 0 {
			p.Remove(*e.BasicEntity)
		}
	}
}

func (p *ProjectileSystem) Fire(pos engo.Point) {
	projectile := Projectile{BasicEntity: ecs.NewBasic()}
	p.Master.AppendChild(&projectile.BasicEntity)

	projectile.RenderComponent = common.RenderComponent{
		Drawable: p.texture,
		Scale:    config.GlobalScale,
	}

	projectile.SpaceComponent = common.SpaceComponent{
		Position: pos,
		Width:    p.texture.Width() * projectile.RenderComponent.Scale.X,
		Height:   p.texture.Height() * projectile.RenderComponent.Scale.Y,
	}

	projectile.CollisionComponent = common.CollisionComponent{Main: 1, Group: 1}

	for _, system := range p.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&projectile.BasicEntity, &projectile.RenderComponent, &projectile.SpaceComponent)
		case *common.CollisionSystem:
			sys.Add(&projectile.BasicEntity, &projectile.CollisionComponent, &projectile.SpaceComponent)
		}
	}

	p.Add(&projectile.BasicEntity, &projectile.SpaceComponent)
}
