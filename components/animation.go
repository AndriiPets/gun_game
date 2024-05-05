package components

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/ganim8/v2"
)

type AnimationType string

var (
	AnimationActor  AnimationType = "actor"
	AnimationStatic AnimationType = "static"
	AnimationFollow AnimationType = "follow"
)

type AnimationData struct {
	Animation     *ganim8.Animation
	Position      math.Vec2
	Rotation      float64
	FlipH         bool
	FlipV         bool
	PlaybackSpeed float64
	Type          AnimationType
}

var Animation = donburi.NewComponentType[AnimationData]()
