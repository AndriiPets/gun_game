package layers

import "github.com/yohamta/donburi/ecs"

const (
	Default ecs.LayerID = iota
	Background
	Architecture
	Actors
	Player
	Interactables
	System
	FX
)
