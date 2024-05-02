package systems

import (
	"fmt"
	mmath "math"

	"github.com/AndriiPets/FishGame/components"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/filter"
)

func UpdateEnemies(ecs *ecs.ECS) {
	query := donburi.NewQuery(filter.Contains(components.Enemy))

	query.Each(ecs.World, func(e *donburi.Entry) {
		health := components.Health.Get(e)
		enemy := components.Enemy.Get(e)
		enemyVelocity := components.Velocity.Get(e)
		obj := components.Object.Get(e)
		weapon := components.Shooter.Get(e)

		//MOVEMENT

		friction := 0.9
		//accel := 0.6
		maxSpeed := 4.0

		// Apply friction and horizontal speed limiting.
		enemyVelocity.Speed *= friction
		if mmath.Abs(enemyVelocity.Speed) < 0.1 {
			enemyVelocity.Speed = 0
			enemyVelocity.Vel = math.NewVec2(0, 0)
		}

		if enemyVelocity.Speed > maxSpeed {
			enemyVelocity.Speed = maxSpeed
		}
		//

		if health.Dead && enemy.State != components.EnemyStateDead {
			fmt.Println("enemy felled!")
			updateEnemyState(e, components.EnemyStateDead)
		}

		if health.Hit && !health.Dead {
			updateEnemyState(e, components.EnemyStateHit)
		}

		if !health.Hit && !health.Dead {
			updateEnemyState(e, components.EnemyStateIdle)
		}

		//update weapon sprite position
		centerX, centerY := obj.Position.X+(obj.Size.X/2), obj.Position.Y+(obj.Size.Y/2)
		weapon.Position = math.NewVec2(centerX, centerY)

	})
}

func updateEnemyState(entry *donburi.Entry, state components.EnemyState) {
	enemy := components.Enemy.Get(entry)
	if enemy.State == state {
		//fmt.Println(flip)
		return
	}

	enemy.State = state

	anim := enemy.Animation()
	//fmt.Println("change state!")

	//update animation
	animation := components.Animation.Get(entry)
	animation.Animation = anim
}
