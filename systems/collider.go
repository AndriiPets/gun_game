package systems

import (
	//"fmt"

	"time"

	"github.com/AndriiPets/FishGame/components"
	dresolv "github.com/AndriiPets/FishGame/resolv"
	"github.com/solarlune/resolv"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

func UpdateCollisions(ecs *ecs.ECS) {
	updatePlayerCollisions(ecs)
	updateBulletCollisions(ecs)
	UpdateBoucerCollisions(ecs)
}

func updatePlayerCollisions(ecs *ecs.ECS) {
	query := donburi.NewQuery(filter.Contains(components.CollistionPlayer, components.Object, components.Velocity, components.Health))
	bQuery := donburi.NewQuery(filter.Contains(components.CollistionBullet))

	query.Each(ecs.World, func(e *donburi.Entry) {
		object := dresolv.GetObject(e)
		velocity := components.Velocity.Get(e)
		UnitVector := velocity.Vel.Normalized().MulScalar(velocity.Speed)
		health := components.Health.Get(e)
		var damage int = 0
		//fmt.Println(velocity.Vel)

		dx := UnitVector.X
		dy := UnitVector.Y

		if col := object.Check(dx, 0); col != nil {
			if col.HasTags("solid") {
				dx = col.ContactWithCell(col.Cells[len(col.Cells)-1]).X
			}

			if bullets := col.ObjectsByTags("bullet"); len(bullets) > 0 {
				bullet := bullets[0]
				if object.Overlaps(bullet) {
					if !health.Dead {
						damage = apply_colision_with_bullet(ecs, e, col, bQuery)
					}
				}
			}

			//if col.HasTags("bullet") {
			//	if !health.Dead {
			//		damage = apply_colision_with_bullet(ecs, e, col, bQuery)
			//	}
			//}
		}

		object.Position.X += dx

		if col := object.Check(0, dy); col != nil {
			if col.HasTags("solid") {
				dy = col.ContactWithCell(col.Cells[len(col.Cells)-1]).Y
			}

			if bullets := col.ObjectsByTags("bullet"); len(bullets) > 0 {
				bullet := bullets[0]
				if object.Overlaps(bullet) {
					if !health.Dead {
						damage = apply_colision_with_bullet(ecs, e, col, bQuery)
					}
				}
			}

			//if col.HasTags("bullet") {
			//	if !health.Dead {
			//		damage = apply_colision_with_bullet(ecs, e, col, bQuery)
			//	}
			//}
		}

		object.Position.Y += dy

		if damage > 0 {
			health.DamageHealth(damage)
		}

	})
}

func check_bullet_collision(ecs *ecs.ECS, col *resolv.Collision, query *donburi.Query) *donburi.Entry {
	var bulletID interface{}
	var bulletEntity *donburi.Entry

	for _, c := range col.Objects {
		if c.Data != nil {
			bulletID = c.Data
		}
	}

	query.Each(ecs.World, func(e *donburi.Entry) {
		if bulletID == e.Id() {
			despawn := components.Despawnable.Get(e)
			despawn.DespawnRequest = true
			bulletEntity = e
		}
	})

	return bulletEntity
}

func apply_colision_with_bullet(ecs *ecs.ECS, e *donburi.Entry, col *resolv.Collision, query *donburi.Query) int {
	velocity := components.Velocity.Get(e)
	health := components.Health.Get(e)

	bullet := check_bullet_collision(ecs, col, query)
	bulletVec := components.Velocity.Get(bullet).Vel.Normalized().MulScalar(5)
	//centerY := object.Position.Y + (object.Size.Y / 2)
	velocity.Vel = bulletVec
	velocity.Speed = 2
	health.Hit = true
	health.HitTime = time.Now()
	damage := 1

	return damage
}

func updateBulletCollisions(ecs *ecs.ECS) {
	query := donburi.NewQuery(filter.Contains(components.CollistionBullet, components.Object, components.Velocity))

	query.Each(ecs.World, func(e *donburi.Entry) {
		object := dresolv.GetObject(e)
		velocity := components.Velocity.Get(e)
		UnitVector := velocity.Vel.Normalized().MulScalar(velocity.Speed)

		despawn := components.Despawnable.Get(e)

		dx := UnitVector.X
		dy := UnitVector.Y

		if col := object.Check(dx, 0); col != nil {
			if col.HasTags("solid") {
				despawn.DespawnRequest = true
			}

		}

		object.Position.X += dx

		if col := object.Check(0, dy); col != nil {
			if col.HasTags("solid") {
				despawn.DespawnRequest = true
			}
		}

		object.Position.Y += dy

	})
}

func UpdateBoucerCollisions(ecs *ecs.ECS) {
	query := donburi.NewQuery(filter.Contains(components.CollisionBouncer, components.Object, components.Velocity))

	query.Each(ecs.World, func(e *donburi.Entry) {
		object := dresolv.GetObject(e)
		velocity := components.Velocity.Get(e)
		UnitVector := velocity.Vel.Normalized().MulScalar(velocity.Speed)

		//despawn := components.Despawnable.Get(e)

		dx := UnitVector.X
		dy := UnitVector.Y

		if col := object.Check(dx, 0); col != nil {
			if col.HasTags("solid") {
				dx = col.ContactWithCell(col.Cells[0]).X
				velocity.Vel.X *= -1
			}
		}

		object.Position.X += dx

		if col := object.Check(0, dy); col != nil {
			if col.HasTags("solid") {
				dy = col.ContactWithCell(col.Cells[0]).Y
				velocity.Vel.Y *= -1
			}
		}

		object.Position.Y += dy

	})
}
