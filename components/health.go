package components

import (
	"fmt"
	"time"

	"github.com/yohamta/donburi"
)

type HealthData struct {
	Ammount   int
	Dead      bool
	DeathLock bool

	Hit      bool
	HitTime  time.Time
	Cooldown float64
}

var Health = donburi.NewComponentType[HealthData]()

func (h *HealthData) DamageHealth(ammount int) {
	fmt.Println("access health!", h.Ammount)
	if !h.DeathLock {

		h.Ammount -= ammount
		fmt.Println("health damaged! remain:", h.Ammount)
		if h.Ammount <= 0 {
			h.Dead = true
		}
	}
}
