package components

import "github.com/yohamta/donburi"

type DespawnableData struct {
	DespawnRequest bool
}

var Despawnable = donburi.NewComponentType[DespawnableData]()
