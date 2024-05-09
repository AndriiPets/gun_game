package ai

import (
	"image/color"

	"github.com/AndriiPets/FishGame/components"
	"github.com/AndriiPets/FishGame/config"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/quasilyte/pathing"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/math"
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
		path := ai.PathCurrent

		pathfinder := components.PathFinder.MustFirst(ecs.World)
		pf := components.PathFinder.Get(pathfinder)

		x, y := obj.Position.X, obj.Position.Y
		cellSize := float64(config.BlockSize)
		pos := math.NewVec2(0, 0)
		switch path {
		case pathing.DirRight:
			pos = math.NewVec2(x+cellSize, y)
		case pathing.DirLeft:
			pos = math.NewVec2(x-cellSize, y)
		case pathing.DirUp:
			pos = math.NewVec2(x, y-cellSize)
		case pathing.DirDown:
			pos = math.NewVec2(x, y+cellSize)
		case pathing.DirNone:
			pos = math.NewVec2(x, y)
		}
		vector.StrokeCircle(screen, float32(obj.Position.X), float32(obj.Position.Y), float32(ai.VisionRadius), 1, color.RGBA{225, 225, 225, 255}, false)
		vector.StrokeRect(screen, float32(pos.X), float32(pos.Y), float32(cellSize), float32(cellSize), 1, color.RGBA{225, 225, 225, 255}, false)
		vector.StrokeRect(screen, float32(obj.Position.X), float32(obj.Position.Y), float32(cellSize), float32(cellSize), 1, color.RGBA{225, 225, 225, 255}, false)
		pf.DrawDebugGrid(screen)

	})
}
