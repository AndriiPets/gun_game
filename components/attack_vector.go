package components

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
)

type AttackVectorData struct {
	Vec math.Vec2
}

var AttackVector = donburi.NewComponentType[AttackVectorData]()
