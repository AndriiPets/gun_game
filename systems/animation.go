package systems

import (
	"github.com/AndriiPets/FishGame/components"
	dresolv "github.com/AndriiPets/FishGame/resolv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
	"github.com/yohamta/ganim8/v2"
)

func UpdateAnimations(ecs *ecs.ECS) {
	components.Animation.Each(ecs.World, func(e *donburi.Entry) {
		a := components.Animation.Get(e)
		a.Animation.Update()
	})
}

func DrawAnimation(layer ecs.LayerID) func(ecs *ecs.ECS, screen *ebiten.Image) {
	animations := ecs.NewQuery(layer, filter.Contains(
		components.Animation,
		components.Object,
	))

	return func(ecs *ecs.ECS, screen *ebiten.Image) {
		animations.Each(ecs.World, func(e *donburi.Entry) {
			a := components.Animation.Get(e)
			o := dresolv.GetObject(e)

			a.Animation.Sprite().SetFlipH(a.FlipH)
			a.Animation.Sprite().SetFlipV(a.FlipV)

			middleX := o.Position.X
			origin_offset := 0.5

			if a.Type == components.AnimationActor {
				middleX = o.Position.X + (o.Size.X / 2)
			}

			if a.Type == components.AnimationStatic {
				origin_offset = 0
			}

			ganim8.DrawAnime(screen, a.Animation, middleX, o.Position.Y, a.Rotation, 1, 1, origin_offset, origin_offset)
		})
	}
}
