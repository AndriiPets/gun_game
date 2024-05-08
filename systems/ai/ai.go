package ai

import (
	"image/color"

	"github.com/AndriiPets/FishGame/components"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

func UpdateAI(ecs *ecs.ECS) {
	query := donburi.NewQuery(filter.Contains(components.AI, components.Object, components.AttackVector, components.Health, components.Shooter))
	playerEntity := components.Player.MustFirst(ecs.World)

	query.Each(ecs.World, func(e *donburi.Entry) {
		health := components.Health.Get(e)

		//if enemy is dead stop the ai
		if !health.Dead {
			UpdateGruntAI(ecs, e, playerEntity)
		}

	})
}

func DrawDebugAi(ecs *ecs.ECS, screen *ebiten.Image) {
	query := donburi.NewQuery(filter.Contains(components.AI))

	query.Each(ecs.World, func(e *donburi.Entry) {
		obj := components.Object.Get(e)
		ai := components.AI.Get(e)
		vector.StrokeCircle(screen, float32(obj.Position.X), float32(obj.Position.Y), float32(ai.VisionRadius), 1, color.RGBA{225, 225, 225, 255}, false)
	})
}
