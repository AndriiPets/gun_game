package tags

import "github.com/yohamta/donburi"

var (
	Player       = donburi.NewTag().SetName("player")
	Wall         = donburi.NewTag().SetName("wall")
	Floor        = donburi.NewTag().SetName("floor")
	Bullet       = donburi.NewTag().SetName("bullet")
	WeaponSprite = donburi.NewTag().SetName("WeaponSprite")
	Enemy        = donburi.NewTag().SetName("enemy")
)
