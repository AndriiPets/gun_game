package resources

type Weapon struct {
	Type     string
	Cooldown float64
	Bullet   string
}

type Projectile struct {
	Size  float64
	Speed float64
}

var WeaponMap = map[string]Weapon{
	"default": {
		Type:     "default",
		Cooldown: 0.2,
		Bullet:   "normal",
	},
	"bouncer": {
		Type:     "bouncer",
		Cooldown: 0.5,
		Bullet:   "bounce",
	},
	"enemy_default": {
		Type:     "default",
		Cooldown: 0.5,
		Bullet:   "normal",
	},
}

var ProjectileMap = map[string]Projectile{
	"normal": {
		Size:  4,
		Speed: 15.0,
	},
	"bounce": {
		Size:  4,
		Speed: 15.0,
	},
}
