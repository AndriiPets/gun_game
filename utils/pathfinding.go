package utils

import (
	"image/color"
	"math"

	"github.com/AndriiPets/FishGame/config"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	path "github.com/quasilyte/pathing"
)

type PathFinder struct {
	Grid   *path.Grid
	Layers path.GridLayer
	BFS    *path.GreedyBFS
	AStar  *path.AStar
}

type PathAlgo string

const (
	TileWall = iota
	TileFloor
)

const (
	BFS   PathAlgo = "bfs"
	AStar PathAlgo = "astar"
)

func NewPathFinder() *PathFinder {
	grid := path.NewGrid(path.GridConfig{
		WorldWidth:  uint(config.MapWidth) * uint(config.BlockSize),
		WorldHeight: uint(config.MapHeigth) * uint(config.BlockSize),
		CellWidth:   uint(config.BlockSize),
		CellHeight:  uint(config.BlockSize),
	})

	bfs := path.NewGreedyBFS(path.GreedyBFSConfig{
		NumCols: uint(grid.NumCols()),
		NumRows: uint(grid.NumRows()),
	})

	aStar := path.NewAStar(path.AStarConfig{
		NumCols: uint(grid.NumCols()),
		NumRows: uint(grid.NumRows()),
	})

	return &PathFinder{Grid: grid, BFS: bfs, AStar: aStar}
}

func (p *PathFinder) GenerateLayout(data [][]rune, wall rune) {

	for y, row := range data {
		for x, val := range row {

			if val == wall {

				p.Grid.SetCellTile(path.GridCoord{X: x, Y: y}, TileWall)

			} else {

				p.Grid.SetCellTile(path.GridCoord{X: x, Y: y}, TileFloor)

			}
		}
	}

	groundltLayer := path.MakeGridLayer([4]uint8{
		TileFloor: 1,
		TileWall:  0,
	})

	//airLayer := path.MakeGridLayer([4]uint8{
	//	TileFloor: 1,
	//	TileWall:  1,
	//})

	p.Layers = groundltLayer
}

func (p *PathFinder) MakePath(startX, startY, endX, endY float64, algo PathAlgo) path.BuildPathResult {
	startPos := p.Grid.PosToCoord(startX, startY)
	endPos := p.Grid.PosToCoord(endX, endY)

	var steps path.BuildPathResult

	switch algo {
	case "bfs":
		steps = p.BFS.BuildPath(p.Grid, startPos, endPos, p.Layers)
	case "astar":
		steps = p.AStar.BuildPath(p.Grid, startPos, endPos, p.Layers)
	}

	return steps
}

func (p *PathFinder) CoordToGrid(x, y float64) path.GridCoord {
	fx := int(math.Floor(x / float64(config.BlockSize)))
	fy := int(math.Floor(y / float64(config.BlockSize)))
	return path.GridCoord{X: fx, Y: fy}
}

func (p *PathFinder) DrawDebugGrid(screen *ebiten.Image) {
	for y := 0; y < p.Grid.NumRows(); y++ {
		for x := 0; x < p.Grid.NumCols(); x++ {
			if p.Grid.GetCellTile(path.GridCoord{X: x, Y: y}) == TileWall {
				posX, posY := p.Grid.CoordToPos(path.GridCoord{X: x, Y: y})
				vector.StrokeRect(screen, float32(posX)-float32(config.BlockSize/2), float32(posY)-float32(config.BlockSize/2), float32(config.BlockSize), float32(config.BlockSize), 2, color.RGBA{225, 225, 225, 255}, false)
			}
		}
	}
}
