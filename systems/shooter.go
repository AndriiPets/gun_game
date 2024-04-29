package systems

import (
	"fmt"
	"math"
	"time"

	"github.com/AndriiPets/FishGame/archetypes"
	"github.com/AndriiPets/FishGame/components"
	"github.com/AndriiPets/FishGame/events"
	dresolv "github.com/AndriiPets/FishGame/resolv"
	"github.com/AndriiPets/FishGame/resources"
	"github.com/AndriiPets/FishGame/tags"
	"github.com/solarlune/resolv"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	dmath "github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/filter"
)

func UpdateShooters(ecs *ecs.ECS) {
	query := donburi.NewQuery(filter.Contains(components.Shooter, components.AttackVector))

	query.Each(ecs.World, func(e *donburi.Entry) {
		shooter := components.Shooter.Get(e)
		weaponData := resources.WeaponMap[shooter.Type]

		if shooter.Fire && shooter.CanFire {
			fmt.Println("Fire shooter\nCooldown:", weaponData.Cooldown)
			//spawn bullet
			spawnBullet(e, ecs)

			//recoil screen shake if fired by player
			if e.HasComponent(components.Player) {
				events.ScreenShakeEvent.Publish(ecs.World, events.ScreenShake{Type: "recoil"})
			}

			shooter.FireTime = time.Now()
			shooter.CanFire = false
			shooter.Fire = false
		}

		if !shooter.CanFire {
			if time.Now().Sub(shooter.FireTime).Seconds() >= weaponData.Cooldown {
				shooter.CanFire = true
				fmt.Println("Cooldown over, can fire")
			}
		}

	})
}

var bulletSpawnMap = map[string]*archetypes.Archetype{
	"normal": archetypes.Bullet,
	"bounce": archetypes.BouncerBullet,
}

func spawnBullet(e *donburi.Entry, ecs *ecs.ECS) {

	shooter := components.Shooter.Get(e)
	attackVec := components.AttackVector.Get(e).Vec
	space := components.Space.MustFirst(ecs.World)

	weaponData := resources.WeaponMap[shooter.Type]
	bulletData := resources.ProjectileMap[weaponData.Bullet]
	fmt.Println(weaponData.Bullet)

	bullet := bulletSpawnMap[weaponData.Bullet].Spawn(ecs)

	//setup animation sprite
	animation := components.Animation.Get(bullet)
	bulletComp := components.Bullet.Get(bullet)
	//rotate image
	angle := math.Atan2(attackVec.Y, attackVec.X)
	animation.Rotation = angle
	animation.Animation = bulletComp.Animation()

	//bullet spawn position
	spawnPosition := shooter.HolderPosition.Add(attackVec.MulScalar(22))

	obj := resolv.NewObject(spawnPosition.X, spawnPosition.Y, bulletData.Size, bulletData.Size)
	obj.AddTags("bullet")
	obj.Data = bullet.Id()
	dresolv.SetObject(bullet, obj)

	components.Velocity.SetValue(bullet, components.VelocityData{
		Vel:   attackVec,
		Speed: bulletData.Speed,
	})

	dresolv.Add(space, bullet)
}

func UpdateWeaponSprite(ecs *ecs.ECS) {
	query := donburi.NewQuery(filter.Contains(tags.WeaponSprite))

	query.Each(ecs.World, func(e *donburi.Entry) {
		shooter := components.Shooter.Get(e)
		anim := components.Animation.Get(e)
		attVec := components.AttackVector.Get(e)

		//update animation obj position based on shooter position
		pos := dmath.Vec2(shooter.Position)
		obj := dresolv.GetObject(e)

		obj.Position.X = pos.X
		obj.Position.Y = pos.Y

		//calculate weapon sprite rotation
		angle := math.Atan2(attVec.Vec.Y, attVec.Vec.X)
		anim.Rotation = angle

		//flip weapon sprite
		playerVec := dmath.NewVec2(pos.X, 0).Normalized()
		dot := attVec.Vec.Dot(&playerVec)

		if dot > 0 {
			anim.FlipV = false
		}
		if dot < 0 {
			anim.FlipV = true
		}

	})
}
