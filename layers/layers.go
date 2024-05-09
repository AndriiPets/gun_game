package layers

import "github.com/yohamta/donburi/ecs"

const (
	Default ecs.LayerID = iota
	Background
	Architecture
	Actors
	FX
	Player
	Interactables
	System
)
