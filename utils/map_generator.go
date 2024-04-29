package utils

import (
	"fmt"

	"github.com/AndriiPets/FishGame/config"
	"github.com/AndriiPets/FishGame/utils/dngn"
)

type GenerationType string

const (
	BSP         GenerationType = "bsp"
	DrunkWalk   GenerationType = "drunk"
	RandomRooms GenerationType = "random"
)

type World struct {
	Map *dngn.Layout
}

func NewWorldMap() *World {
	world := &World{
		Map: dngn.NewLayout(config.MapWidth, config.MapHeigth),
	}

	return world
}

func (w *World) GenerateMap(genType GenerationType) {
	mapSelection := w.Map.Select()

	switch genType {
	case BSP:
		bspOptions := dngn.NewDefaultBSPOptions()
		bspOptions.SplitCount = 100
		bspOptions.MinimumRoomSize = 5

		bspRooms := w.Map.GenerateBSP(bspOptions)

		start := bspRooms[0]

		for _, subroom := range bspRooms {

			subroomCenter := subroom.Center()
			center := w.Map.Center()

			margin := 10

			if subroomCenter.X > center.X-margin &&
				subroomCenter.X < center.X+margin &&
				subroomCenter.Y > center.Y-margin &&
				subroomCenter.Y < center.Y+margin {
				start = subroom
				break
			}

		}

		for _, room := range bspRooms {

			hops := room.CountHopsTo(start)
			//mapSelection.FilterByArea(room.X, room.Y, room.W+1, room.H+1).Fill('Z')

			if hops < 0 || hops > 4 {
				// We're filtering out a little bit more on the width and height because the walls and doorways in GenerateBSP() are always on the top and left sides of each room.
				// By adding the right and bottom as well, we can nuke any doors that led into rooms we're deleting.
				mapSelection.FilterByArea(room.X, room.Y, room.W+1, room.H+1).Fill('x')
				room.Disconnect()
			}

		}

		player_pos := start.Center()
		fmt.Println(w.Map.Get(player_pos.X, player_pos.Y))
		w.Map.Set(player_pos.X, player_pos.Y, 'P')

	case DrunkWalk:
		startX, startY := w.Map.GenerateDrunkWalk(' ', 'x', 0.8)
		w.Map.Set(startX, startY, 'P')

	case RandomRooms:
		rooms := w.Map.GenerateRandomRooms(' ', 'x', 10, 3, 3, 5, 5, true)
		start := rooms[0]
		w.Map.Set(start[0], start[1], 'P')

		// This selects the ground tiles that are between walls to place doors randomly. This isn't really good, but it at least
		// gets the idea across.
		mapSelection.FilterByRune(' ').FilterBy(func(x, y int) bool {
			return (w.Map.Get(x+1, y) == 'x' && w.Map.Get(x-1, y) == 'x') || (w.Map.Get(x, y-1) == 'x' && w.Map.Get(x, y+1) == 'x')
		}).FilterByPercentage(0.25).Fill('#')
	}

	// Fill the outer walls
	mapSelection.Remove(mapSelection.FilterByArea(1, 1, w.Map.Width-2, w.Map.Height-2)).Fill('x')

	// Add a different tile for an alternate floor
	mapSelection.FilterByRune(' ').FilterByPercentage(0.1).Fill('.')
	mapSelection.FilterByRune(' ').FilterByPercentage(0.02).Fill('e')

	fmt.Println(w.Map.DataToString())
}
