package factory

import (
	"github.com/AndriiPets/FishGame/archetypes"
	"github.com/AndriiPets/FishGame/assets"
	"github.com/AndriiPets/FishGame/components"
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/ganim8/v2"

	dresolv "github.com/AndriiPets/FishGame/resolv"
)

type ParticleType string

const (
	ParticleDust     ParticleType = "dust"
	ParticleGunFlash ParticleType = "gun_flash"
)

func CreateParticle(ecs *ecs.ECS, posX, posY float64, pType ParticleType, rotation float64, flipH, flipV bool) *donburi.Entry {
	particle := archetypes.ParticleSprite.Spawn(ecs)

	//set animation
	animation := components.Animation.Get(particle)
	anim := assets.GetAnimation("particle_" + string(pType))
	anim.SetOnLoop(ganim8.PauseAtEnd)

	animation.Animation = anim
	animation.Rotation = rotation
	animation.FlipH = flipH
	animation.FlipV = flipV

	//setup particle object
	obj := resolv.NewObject(posX, posY, 8, 8)
	dresolv.SetObject(particle, obj)

	return particle

}
