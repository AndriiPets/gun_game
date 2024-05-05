package systems

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/AndriiPets/FishGame/components"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	dmath "github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/filter"
)

func UpdateAI(ecs *ecs.ECS) {
	query := donburi.NewQuery(filter.Contains(components.AI, components.Object, components.AttackVector, components.Health, components.Shooter))
	playerEntity := components.Player.MustFirst(ecs.World)
	playerObjPos := components.Object.Get(playerEntity).Position

	query.Each(ecs.World, func(e *donburi.Entry) {
		health := components.Health.Get(e)

		//if enemy is dead stop the ai
		if !health.Dead {
			ai := components.AI.Get(e)
			obj := components.Object.Get(e)
			attVec := components.AttackVector.Get(e)
			shooter := components.Shooter.Get(e)

			//check if player in vision circle
			if in_circle(obj.Position.X, obj.Position.Y, ai.VisionRadius, playerObjPos.X, playerObjPos.Y) {
				//change attack vector to follow the player
				playerVec := dmath.NewVec2(playerObjPos.X, playerObjPos.Y)
				enemyVec := dmath.NewVec2(obj.Position.X, obj.Position.Y)
				attVec.Vec = playerVec.Sub(enemyVec).Normalized()

				if roll_attack_initiative(ai.AgressionModifier, 100) && shooter.CanFire {
					shooter.Fire = true
				}
			}
		}

	})
}

func in_circle(centerX, centerY, radius, x, y float64) bool {
	square_dist := math.Pow(centerX-x, 2) + math.Pow(centerY-y, 2)
	return square_dist <= math.Pow(radius, 2)
}

func roll_attack_initiative(modifier, limit int) bool {
	roll := rand.Intn(limit)
	return roll <= modifier
}

func DrawDebugAi(ecs *ecs.ECS, screen *ebiten.Image) {
	query := donburi.NewQuery(filter.Contains(components.AI))

	query.Each(ecs.World, func(e *donburi.Entry) {
		obj := components.Object.Get(e)
		ai := components.AI.Get(e)
		vector.StrokeCircle(screen, float32(obj.Position.X), float32(obj.Position.Y), float32(ai.VisionRadius), 1, color.RGBA{225, 225, 225, 255}, false)
	})
}
