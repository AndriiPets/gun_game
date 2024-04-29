package components

import (
	//"github.com/solarlune/resolv"
	"time"

	"github.com/AndriiPets/FishGame/assets"
	"github.com/yohamta/donburi"
	"github.com/yohamta/ganim8/v2"
)

type PlayerData struct {
	FacingRight bool
	IsDashing   bool
	State       PlayerState
	DashTimer   time.Time
}

type PlayerState string

var (
	PlayerStateIdle PlayerState = "idle"
	PlayerStateRun  PlayerState = "run"
)

var Player = donburi.NewComponentType[PlayerData]()

func (p *PlayerData) Animation() *ganim8.Animation {
	return assets.GetAnimation("player_" + string(p.State))
}
