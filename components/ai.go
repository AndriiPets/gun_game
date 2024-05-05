package components

import "github.com/yohamta/donburi"

type AIType int

const (
	AITypeShooter AIType = iota
	AITypeBrawler
)

type AIData struct {
	AIType            AIType
	VisionRadius      float64
	AgressionModifier int
}

var AI = donburi.NewComponentType[AIData]()
