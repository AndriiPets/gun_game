package ai

import (
	"fmt"

	"github.com/AndriiPets/FishGame/components"
	"github.com/quasilyte/pathing"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	dmath "github.com/yohamta/donburi/features/math"
)

func UpdateGruntAI(ecs *ecs.ECS, enemy *donburi.Entry, player *donburi.Entry) {
	ai := components.AI.Get(enemy)
	obj := components.Object.Get(enemy)
	attVec := components.AttackVector.Get(enemy)
	shooter := components.Shooter.Get(enemy)

	playerObj := components.Object.Get(player)
	playerObjPos := playerObj.Position

	ai.PathCurrent = pathing.DirNone

	//check if player in vision circle
	if in_circle(obj.Position.X, obj.Position.Y, ai.VisionRadius, playerObjPos.X, playerObjPos.Y) {
		//change attack vector to follow the player
		playerVec := dmath.NewVec2(playerObjPos.X, playerObjPos.Y)
		enemyVec := dmath.NewVec2(obj.Position.X, obj.Position.Y)
		attVec.Vec = playerVec.Sub(enemyVec).Normalized()
		ai.Path = atempt_build_path(ecs, obj, playerObj)

		//before shooting make shure actor has line of sight
		if line_of_sight_check(ecs, obj, playerObj) {

			if roll_attack_initiative(ai.AgressionModifier, 100) && shooter.CanFire {
				shooter.Fire = true
			}

		}

		ai.PathCurrent = ai.Path.Steps.Next()
		fmt.Println(ai.Path.Steps.String(), ai.PathCurrent)
	}
}
