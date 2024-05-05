package systems

import (
	"github.com/AndriiPets/FishGame/components"
	"github.com/AndriiPets/FishGame/tags"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

func UpdateParticles(ecs *ecs.ECS) {
	query := donburi.NewQuery(filter.Contains(tags.Particle))

	query.Each(ecs.World, func(e *donburi.Entry) {
		animation := components.Animation.Get(e)
		despawn := components.Despawnable.Get(e)

		if animation.Animation.IsEnd() {
			despawn.DespawnRequest = true
		}
	})
}
