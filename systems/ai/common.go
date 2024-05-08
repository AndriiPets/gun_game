package ai

import (
	"math"
	"math/rand"

	"github.com/AndriiPets/FishGame/components"
	"github.com/AndriiPets/FishGame/utils"
	"github.com/quasilyte/pathing"
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi/ecs"
)

func in_circle(centerX, centerY, radius, x, y float64) bool {
	square_dist := math.Pow(centerX-x, 2) + math.Pow(centerY-y, 2)
	return square_dist <= math.Pow(radius, 2)
}

func roll_attack_initiative(modifier, limit int) bool {
	roll := rand.Intn(limit)
	return roll <= modifier
}

// TODO: path seems to work need to build a movement system
func atempt_build_path(ecs *ecs.ECS, startObj, endObj *resolv.Object) pathing.BuildPathResult {
	pathfinder := components.PathFinder.MustFirst(ecs.World)
	pf := components.PathFinder.Get(pathfinder)

	steps := pf.MakePath(
		startObj.Position.X,
		startObj.Position.Y,
		endObj.Position.X,
		endObj.Position.Y,
		utils.BFS,
	)

	return steps
}

func line_of_sight_check(ecs *ecs.ECS, startObj, endObj *resolv.Object) bool {
	spaceEntry := components.Space.MustFirst(ecs.World)
	space := components.Space.Get(spaceEntry)

	cx, cy := startObj.CellPosition()
	mx, my := endObj.CellPosition()

	sightLine := space.CellsInLine(cx, cy, mx, my)
	//fmt.Println(startX, startY)

	for i, cell := range sightLine {

		//if not yourself
		if i == 0 {
			continue
		}

		if cell.ContainsTags("solid") {
			return false
		}
	}

	return true

}
