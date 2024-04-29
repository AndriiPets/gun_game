package components

import (
	"github.com/AndriiPets/FishGame/assets"
	"github.com/yohamta/donburi"
	"github.com/yohamta/ganim8/v2"
)

type EnemyData struct {
	Type  EnemyType
	State EnemyState
}

type EnemyType string
type EnemyState string

const (
	EnemyStateIdle EnemyState = "idle"
	EnemyStateRun  EnemyState = "run"
	EnemyStateDead EnemyState = "dead"
	EnemyStateHit  EnemyState = "hit"

	EnemyTypeGrunt EnemyType = "orc"
)

var Enemy = donburi.NewComponentType[EnemyData]()

func (e *EnemyData) Animation() *ganim8.Animation {
	return assets.GetAnimation(string(e.Type) + "_" + string(e.State))
}
