package systems

import (
	"fmt"
	"time"

	"github.com/AndriiPets/FishGame/components"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

func UpdateHealth(ecs *ecs.ECS) {
	query := donburi.NewQuery(filter.Contains(components.Health))

	query.Each(ecs.World, func(e *donburi.Entry) {
		health := components.Health.Get(e)

		if health.Dead && !health.DeathLock {
			//trigger death event
			fmt.Println("just died!")
			health.DeathLock = true
		}

		if health.Hit {
			if time.Now().Sub(health.HitTime).Seconds() >= health.Cooldown {
				health.Hit = false
			}
		}
	})
}
