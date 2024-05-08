package components

import (
	"github.com/quasilyte/pathing"
	"github.com/yohamta/donburi"
)

type AIType int

const (
	AITypeShooter AIType = iota
	AITypeBrawler
)

type AIData struct {
	AIType            AIType
	VisionRadius      float64
	AgressionModifier int
	Path              pathing.BuildPathResult
	PathCurrent       pathing.Direction
}

var AI = donburi.NewComponentType[AIData]()
