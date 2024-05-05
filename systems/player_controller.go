package systems

import (
	"fmt"
	"image/color"
	mmath "math"
	"time"

	"github.com/AndriiPets/FishGame/components"
	"github.com/AndriiPets/FishGame/factory"
	"github.com/AndriiPets/FishGame/tags"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/math"

	dresolv "github.com/AndriiPets/FishGame/resolv"
)

var dashVec math.Vec2

func UpdatePlayer(ecs *ecs.ECS) {
	playerEntity, _ := components.Player.First(ecs.World)
	playerVelocity := components.Velocity.Get(playerEntity)
	player := components.Player.Get(playerEntity)
	anim := components.Animation.Get(playerEntity)
	attackVec := components.AttackVector.Get(playerEntity).Vec
	playerObj := dresolv.GetObject(playerEntity)

	shooter := components.Shooter.Get(playerEntity)

	//MOVEMENT
	//dx, dy := 0.0, 0.0 //direction vector
	friction := 0.9
	accel := 0.4
	maxSpeed := 3.0

	dashCooldown := 0.3
	particleCooldown := 0.5

	if !player.IsDashing {

		//update direction
		if ebiten.IsKeyPressed(ebiten.KeyW) {
			playerVelocity.Vel = math.NewVec2(0.0, -1.0)
			playerVelocity.Speed += accel
		}

		if ebiten.IsKeyPressed(ebiten.KeyS) {
			playerVelocity.Vel = math.NewVec2(0.0, 1.0)
			playerVelocity.Speed += accel
		}

		if ebiten.IsKeyPressed(ebiten.KeyA) {
			playerVelocity.Vel = math.NewVec2(-1.0, 0.0)
			playerVelocity.Speed += accel
		}

		if ebiten.IsKeyPressed(ebiten.KeyD) {
			playerVelocity.Vel = math.NewVec2(1.0, 0.0)
			playerVelocity.Speed += accel
		}

		//diagonal movement
		if ebiten.IsKeyPressed(ebiten.KeyW) && ebiten.IsKeyPressed(ebiten.KeyA) {
			playerVelocity.Vel = math.NewVec2(-1.0, -1.0)
			playerVelocity.Speed += accel
		}
		if ebiten.IsKeyPressed(ebiten.KeyW) && ebiten.IsKeyPressed(ebiten.KeyD) {
			playerVelocity.Vel = math.NewVec2(1.0, -1.0)
			playerVelocity.Speed += accel
		}
		if ebiten.IsKeyPressed(ebiten.KeyS) && ebiten.IsKeyPressed(ebiten.KeyA) {
			playerVelocity.Vel = math.NewVec2(-1.0, 1.0)
			playerVelocity.Speed += accel
		}
		if ebiten.IsKeyPressed(ebiten.KeyS) && ebiten.IsKeyPressed(ebiten.KeyD) {
			playerVelocity.Vel = math.NewVec2(1.0, 1.0)
			playerVelocity.Speed += accel
		}
		//

		//TODO: player velocity resets to 0 when run move button is not pushed

		// Apply friction and horizontal speed limiting.
		playerVelocity.Speed *= friction
		if mmath.Abs(playerVelocity.Speed) < 0.1 {
			playerVelocity.Speed = 0
			playerVelocity.Vel = math.NewVec2(0, 0)
		}

		if playerVelocity.Speed > maxSpeed {
			playerVelocity.Speed = maxSpeed
		}
		//

		//fmt.Println(playerVelocity.Speed)

		//dash controls
		if ebiten.IsKeyPressed(ebiten.KeyShiftLeft) {
			player.IsDashing = true
			fmt.Println(playerVelocity.Speed)

			player.DashTimer = time.Now()

			if playerVelocity.Vel.IsZero() {
				dashVec = attackVec
			} else {
				dashVec = playerVelocity.Vel
			}

		}

	} else {

		playerVelocity.Vel = dashVec

		playerVelocity.Speed = maxSpeed * 2

	}

	if time.Now().Sub(player.DashTimer).Seconds() >= dashCooldown {
		player.IsDashing = false
	}

	//playerVelocity.Speed = maxSpeed
	//playerVelocity.Vel = math.NewVec2(dx, dy)

	//update plater facing direction
	playerVec := math.NewVec2(playerObj.Position.X, 0).Normalized()
	dot := attackVec.Dot(&playerVec)

	if dot > 0 {
		anim.FlipH = false
		//fmt.Println("flip false")
	}
	if dot < 0 {
		anim.FlipH = true
		//fmt.Println("flip true")
	}
	//

	//update animation
	if !playerVelocity.Vel.IsZero() {
		updatePlayerState(playerEntity, components.PlayerStateRun)

		//dust particle spawn
		playerSpawnDustParticle(ecs, playerEntity, particleCooldown)
		//

	} else {
		updatePlayerState(playerEntity, components.PlayerStateIdle)
	}

	//updatePlayerDir(playerEntity, flip)

	//Shooting
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if shooter.CanFire && !player.IsDashing {
			shooter.Fire = true
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		if shooter.Type == "default" {
			shooter.Type = "bouncer"
		} else {
			shooter.Type = "default"
		}
	}

}

func PlayerString(ecs *ecs.ECS) string {
	playerEntity, _ := components.Player.First(ecs.World)
	p := dresolv.GetObject(playerEntity)
	return fmt.Sprintf(
		"Player T: %.1f",
		math.NewVec2(p.Position.X, p.Position.Y),
	)
}

func playerSpawnDustParticle(ecs *ecs.ECS, playerEntry *donburi.Entry, cooldown float64) {
	player := components.Player.Get(playerEntry)
	anim := components.Animation.Get(playerEntry)
	playerObj := dresolv.GetObject(playerEntry)

	//dust particle spawn
	if !player.ParticleSpawn {
		player.ParticleTimer = time.Now()
		player.ParticleSpawn = true
	}

	if time.Now().Sub(player.ParticleTimer).Seconds() >= cooldown {
		factory.CreateParticle(
			ecs,
			playerObj.Position.X,
			playerObj.Position.Y,
			factory.ParticleDust,
			0,
			anim.FlipH,
			false,
		)
		player.ParticleSpawn = false
	}
	//
}

func updatePlayerState(entry *donburi.Entry, state components.PlayerState) {
	player := components.Player.Get(entry)
	if player.State == state {
		//fmt.Println(flip)
		return
	}

	player.State = state

	//update animation
	animation := components.Animation.Get(entry)
	animation.Animation = player.Animation()
}

func DrawPlayer(ecs *ecs.ECS, screen *ebiten.Image) {
	tags.Player.Each(ecs.World, func(e *donburi.Entry) {

		o := dresolv.GetObject(e)
		playerColor := color.RGBA{0, 255, 60, 255}
		//player := components.Player.Get(e)
		attackVec := components.AttackVector.Get(e).Vec.MulScalar(15)
		playerWeapon := components.Shooter.Get(e)
		//fmt.Println(attackVec)

		//player weapon position
		centerX, centerY := o.Position.X+(o.Size.X/2), o.Position.Y+(o.Size.Y/2)
		weaponPosX, weaponPosY := centerX+attackVec.X, centerY+attackVec.Y

		playerWeapon.Position = math.NewVec2(weaponPosX, weaponPosY)
		playerWeapon.HolderPosition = math.NewVec2(centerX, centerY)

		vector.DrawFilledRect(screen, float32(o.Position.X), float32(o.Position.Y), float32(o.Size.X), float32(o.Size.Y), playerColor, false)
		//vector.DrawFilledCircle(screen, float32(weaponPosX), float32(weaponPosY), 7, color.RGBA{255, 255, 255, 255}, false)
	})
}
