package components

import (
	"github.com/AndriiPets/FishGame/assets"
	"github.com/yohamta/donburi"
	"github.com/yohamta/ganim8/v2"
)

type BulletData struct {
	IsDead bool
}

var Bullet = donburi.NewComponentType[BulletData]()

func (db *BulletData) Animation() *ganim8.Animation {
	return assets.GetAnimation("bullet_default")
}
