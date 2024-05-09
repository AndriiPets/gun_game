package components

import (
	"time"

	"github.com/AndriiPets/FishGame/assets"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/ganim8/v2"
)

//TODO position for weapon

type ShooterData struct {
	Type           string
	Fire           bool
	FireTime       time.Time
	Cooldown       float64
	CanFire        bool
	Position       math.Vec2
	HolderPosition math.Vec2
	WeaponFlash    bool
	HoldRange      float64
}

var Shooter = donburi.NewComponentType[ShooterData]()

func (sd *ShooterData) Animation() *ganim8.Animation {
	return assets.GetAnimation("weapon_" + sd.Type)
}
