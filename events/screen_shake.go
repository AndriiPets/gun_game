package events

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/events"

	//"github.com/yohamta/donburi/features/math"

	"github.com/AndriiPets/FishGame/components"
	//dresolv "github.com/AndriiPets/FishGame/resolv"
)

type ScreenShake struct {
	Type string
}

var ScreenShakeEvent = events.NewEventType[ScreenShake]()

func OnRecoilScreenShake(w donburi.World, event ScreenShake) {

	cameraEntity, ok := components.Camera.First(w)

	if ok != true {
		panic("no camera!")
	}

	camera := components.Camera.Get(cameraEntity)

	if event.Type == "recoil" {

		playerEntity, _ := components.Player.First(w)
		//playerObj := dresolv.GetObject(playerEntity)

		//playerPos = math.NewVec2(playerObj.Position.X, playerObj.Position.Y)

		//Calculate recoil vector
		attackVec := components.AttackVector.Get(playerEntity).Vec.MulScalar(-7)

		//centerX, centerY := playerObj.Position.X+(playerObj.Size.X/2), playerObj.Position.Y+(playerObj.Size.Y/2)
		//weaponPosX, weaponPosY := centerX+attackVec.X, centerY+attackVec.Y

		camera.Recoil = attackVec
		//camera.Zoom = -5

		camera.Flash = true

	}
}

func WeaponFlash(w donburi.World, event ScreenShake) {

	cameraEntity, _ := components.Camera.First(w)

	camera := components.Camera.Get(cameraEntity)
	camera.Flash = true

}

func SetupEvents(ecs *ecs.ECS) {
	ScreenShakeEvent.Subscribe(ecs.World, OnRecoilScreenShake)
}

func UpdateEvents(ecs *ecs.ECS) {
	events.ProcessAllEvents(ecs.World)
}
