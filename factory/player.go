package factory

import (
	"github.com/AndriiPets/FishGame/archetypes"
	"github.com/AndriiPets/FishGame/components"
	dresolv "github.com/AndriiPets/FishGame/resolv"

	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/math"
)

func CreatePlayer(ecs *ecs.ECS, posX, PosY float64) *donburi.Entry {
	player := archetypes.Player.Spawn(ecs)

	//setup player initial state
	pl := components.Player.Get(player)
	pl.State = components.PlayerStateIdle

	//setup animation
	animation := components.Animation.Get(player)
	animation.Animation = pl.Animation()
	animation.FlipH = false
	animation.Type = components.AnimationActor

	obj := resolv.NewObject(posX, PosY, 16, 16)
	//obj.AddTags("damageble")
	dresolv.SetObject(player, obj)
	components.Player.SetValue(player, components.PlayerData{
		FacingRight: true,
		IsDashing:   false,
	})
	components.Shooter.SetValue(player, components.ShooterData{
		Type:    "default", //bouncer, //default
		Fire:    false,
		CanFire: true,
	})
	components.Velocity.SetValue(player, components.VelocityData{
		Vel: math.NewVec2(0, 0),
	})
	components.Health.SetValue(player, components.HealthData{
		Cooldown: 0.2,
	})

	//setup weapon sprite
	//space := components.Space.MustFirst(ecs.World)

	wSprite := archetypes.WeaponSprite.Spawn(ecs)
	shooter := components.Shooter.Get(player)
	wAnimation := components.Animation.Get(wSprite)

	wSprite.SetComponent(components.Shooter, player.Component(components.Shooter))
	wSprite.SetComponent(components.AttackVector, player.Component(components.AttackVector))

	dresolv.SetObject(wSprite, resolv.NewObject(posX, PosY, 16, 16))
	wAnimation.Animation = shooter.Animation()
	wAnimation.Type = components.AnimationFollow

	//dresolv.Add(space, wSprite)

	return player
}
