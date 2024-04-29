package components

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
)

type CameraData struct {
	ViewPort math.Vec2
	Position math.Vec2
	Rotation int
	Zoom     int
	CursorX  float64
	CursorY  float64
	Recoil   math.Vec2
	Flash    bool
}

var Camera = donburi.NewComponentType[CameraData]()
