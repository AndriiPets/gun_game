package ai

import (
	"fmt"

	"github.com/AndriiPets/FishGame/components"
	"github.com/hajimehoshi/ebiten/v2"
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
		if blockNum, canSee := line_of_sight_check(ecs, obj, playerObj); canSee {

			if roll_attack_initiative(ai.AgressionModifier, 100) && shooter.CanFire {
				shooter.Fire = true
			}

			//if player is too close attempt to flee
			if blockNum <= 10 {
				ai.PathCurrent = ai.Path.Steps.Next().Reversed()
			}

		} else {
			ai.PathCurrent = ai.Path.Steps.Next()
		}

		if ebiten.IsKeyPressed(ebiten.KeyV) {
			ai.Path.Steps.Rewind()
		}
		fmt.Println(ai.Path.Steps.String(), ai.PathCurrent)
	}
}
