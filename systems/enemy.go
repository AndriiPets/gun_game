package systems

import (
	"fmt"
	"image/color"
	mmath "math"

	"github.com/AndriiPets/FishGame/components"
	"github.com/AndriiPets/FishGame/tags"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/filter"

	dresolv "github.com/AndriiPets/FishGame/resolv"
)

func UpdateEnemies(ecs *ecs.ECS) {
	query := donburi.NewQuery(filter.Contains(tags.Enemy))

	query.Each(ecs.World, func(e *donburi.Entry) {
		health := components.Health.Get(e)
		enemy := components.Enemy.Get(e)
		enemyVelocity := components.Velocity.Get(e)
		obj := components.Object.Get(e)
		weapon := components.Shooter.Get(e)
		attVec := components.AttackVector.Get(e).Vec

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

		//update enemy facing direction
		enemyVec := math.NewVec2(obj.Position.X, 0).Normalized()
		dot := attVec.Dot(&enemyVec)
		//fmt.Println(dot)

		if dot > 0 {
			flip(e, false)
			//fmt.Println("flip false")
		}
		if dot < 0 {
			flip(e, true)
			//fmt.Println("flip true")
		}
		//

		//update weapon sprite position
		centerX, centerY := obj.Position.X+(obj.Size.X/2), obj.Position.Y+(obj.Size.Y/2)
		weapon.Position = math.NewVec2(centerX+attVec.X, centerY+attVec.Y)
		weapon.HolderPosition = math.NewVec2(centerX, centerY)

	})
}

func flip(entry *donburi.Entry, flipH bool) {
	anim := components.Animation.Get(entry)

	if anim.FlipH == flipH {
		return
	}

	anim.FlipH = flipH
	//fmt.Println(anim.FlipH, entry.Id())
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

func DrawEnemyes(ecs *ecs.ECS, screen *ebiten.Image) {
	tags.Enemy.Each(ecs.World, func(e *donburi.Entry) {

		o := dresolv.GetObject(e)
		playerColor := color.RGBA{0, 255, 60, 255}

		vector.DrawFilledRect(screen, float32(o.Position.X), float32(o.Position.Y), float32(o.Size.X), float32(o.Size.Y), playerColor, false)
		//vector.DrawFilledCircle(screen, float32(weaponPosX), float32(weaponPosY), 7, color.RGBA{255, 255, 255, 255}, false)
	})
}
