package components

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
)

type VelocityData struct {
	Vel   math.Vec2
	Speed float64
}

var Velocity = donburi.NewComponentType[VelocityData]()
