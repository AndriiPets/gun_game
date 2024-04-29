package systems

import (
	"github.com/AndriiPets/FishGame/components"
	dresolv "github.com/AndriiPets/FishGame/resolv"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/hierarchy"
	"github.com/yohamta/donburi/filter"
)

func UpdateDespawnable(ecs *ecs.ECS) {
	query := donburi.NewQuery(filter.Contains(components.Despawnable))

	query.Each(ecs.World, func(e *donburi.Entry) {
		despawn := components.Despawnable.Get(e)

		if despawn.DespawnRequest {
			if e.HasComponent(components.Object) {
				space := components.Space.MustFirst(ecs.World)

				dresolv.Remove(space, e)
			}
			hierarchy.RemoveRecursive(e)
		}
	})
}
