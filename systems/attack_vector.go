package systems

import (
	"github.com/AndriiPets/FishGame/components"
	dresolv "github.com/AndriiPets/FishGame/resolv"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/math"
)

func UpdateAttackVector(ecs *ecs.ECS) {
	//get camera cursor position
	cameraEntity, _ := components.Camera.First(ecs.World)
	camera := components.Camera.Get(cameraEntity)

	playerEntity, _ := components.Player.First(ecs.World)
	playerObj := dresolv.GetObject(playerEntity)

	//player attack vector calc
	mouseVec := math.NewVec2(camera.CursorX, camera.CursorY)
	playerVec := math.NewVec2(playerObj.Position.X, playerObj.Position.Y)

	playerAttackVec := mouseVec.Sub(playerVec)
	playerAttackUnit := playerAttackVec.Normalized()

	components.AttackVector.Each(ecs.World, func(e *donburi.Entry) {

		if e.HasComponent(components.Player) {
			components.AttackVector.SetValue(e, components.AttackVectorData{
				Vec: playerAttackUnit,
			})
		}
	})
}
