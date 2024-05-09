package systems

import (
	"image/color"
	"math"
	"time"

	"github.com/AndriiPets/FishGame/archetypes"
	"github.com/AndriiPets/FishGame/components"
	"github.com/AndriiPets/FishGame/events"
	"github.com/AndriiPets/FishGame/factory"
	dresolv "github.com/AndriiPets/FishGame/resolv"
	"github.com/AndriiPets/FishGame/resources"
	"github.com/AndriiPets/FishGame/tags"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/solarlune/resolv"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"

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
			//fmt.Println("Fire shooter\nCooldown:", weaponData.Cooldown)
			//spawn bullet
			spawnBullet(e, ecs)

			//recoil screen shake if fired by player
			if e.HasComponent(components.Player) {
				events.ScreenShakeEvent.Publish(ecs.World, events.ScreenShake{Type: "recoil"})
			}

			//weapon sprite recoil
			events.WeaponRecoilEvent.Publish(ecs.World, events.WeaponRecoil{Entry: e})
			shooter.WeaponFlash = true

			shooter.FireTime = time.Now()
			shooter.CanFire = false
			shooter.Fire = false
		}

		if !shooter.CanFire {
			if time.Now().Sub(shooter.FireTime).Seconds() >= weaponData.Cooldown {
				shooter.CanFire = true
				//fmt.Println("Cooldown over, can fire")
			}
		}

	})
}

func DrawWeaponFlash(ecs *ecs.ECS, screen *ebiten.Image) {
	query := donburi.NewQuery(filter.Contains(components.Shooter, components.AttackVector))

	query.Each(ecs.World, func(e *donburi.Entry) {
		shooter := components.Shooter.Get(e)
		weaponData := resources.WeaponMap[shooter.Type]
		attackVec := components.AttackVector.Get(e).Vec
		cooldown := weaponData.Cooldown / 2
		spawnPosition := shooter.HolderPosition.Add(attackVec.MulScalar(22))

		if !shooter.CanFire {

			if time.Now().Sub(shooter.FireTime).Seconds() <= cooldown {
				vector.DrawFilledCircle(screen, float32(spawnPosition.X), float32(spawnPosition.Y), 7, color.RGBA{225, 225, 225, 255}, false)
				//vector.DrawFilledRect(screen, float32(spawnPosition.X), float32(spawnPosition.Y), 32, 32, color.RGBA{225, 225, 225, 255}, false)

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
	//fmt.Println(weaponData.Bullet)

	bullet := bulletSpawnMap[weaponData.Bullet].Spawn(ecs)

	//setup animation sprite
	animation := components.Animation.Get(bullet)
	bulletComp := components.Bullet.Get(bullet)

	//rotate image
	angle := math.Atan2(attackVec.Y, attackVec.X)
	animation.Rotation = angle
	animation.Animation = bulletComp.Animation()

	//bullet spawn position
	spawnPosition := shooter.HolderPosition.Add(attackVec.MulScalar(24))

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
		weaponData := resources.WeaponMap[shooter.Type]

		//update animation obj position based on shooter position
		var pos dmath.Vec2
		ran := shooter.HoldRange

		//if on cooldown recoil sprite backwards
		if shooter.WeaponFlash {
			//weapon fierd apply new easing function for recoil
			anim.Ease = gween.New(1, float32(shooter.HoldRange), float32(weaponData.Cooldown), ease.Linear)
			ran = 1

			//Spawn weapon flash animation
			spawnPosition := shooter.HolderPosition.Add(attVec.Vec.MulScalar(37))
			factory.CreateParticle(
				ecs,
				spawnPosition.X,
				spawnPosition.Y,
				factory.ParticleGunFlash,
				anim.Rotation,
				anim.FlipH,
				anim.FlipV,
			)
			shooter.WeaponFlash = false

		}
		if !shooter.CanFire {
			//update weapon position based on easing function
			curr, finish := anim.Ease.Update(float32(ecs.Time.DeltaTime().Seconds()))
			if !finish {
				ran = float64(curr)
			}
		}

		pos = shooter.HolderPosition.Add(attVec.Vec.MulScalar(ran))

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
