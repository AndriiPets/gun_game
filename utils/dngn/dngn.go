/*
Package dngn is a simple random map generation library primarily made to be used for 2D games. It features a simple API,
and a couple of different means to generate maps. The easiest way to kick things off when using dngn is to simply create a Layout
to represent your overall game map, which can then be manipulated or have a Generate function run on it to actually generate the
content on the map.
*/
package dngn

import (
	"fmt"
	"image"
	"math"
	"math/rand"
	"sort"
	"time"
)

// Layout represents a dungeon map.
// Width and Height are the width and height of the Layout in the layout. This determines the size of the overall Data structure
// backing the Layout layout.
// Data is the core underlying data structure representing the dungeon. It's a 2D array of runes.
// Seed is the seed of the Layout to use when doing random generation using the Generate* functions below. By default, the seed is
// the lowest possible negative number (math.MinInt64), and so will use the time to set the seed.
type Layout struct {
	Width, Height int
	Data          [][]rune
	Seed          int64
}

// NewLayout returns a new Layout with the specified width and height.
func NewLayout(width, height int) *Layout {

	r := &Layout{Width: width, Height: height, Seed: math.MinInt64}
	r.Data = [][]rune{}
	for y := 0; y < height; y++ {
		r.Data = append(r.Data, []rune{})
		for x := 0; x < width; x++ {
			r.Data[y] = append(r.Data[y], ' ')
		}
	}

	return r

}

// NewLayoutFromRuneArrays creates a new Layout with the data contained in the provided rune arrays.
func NewLayoutFromRuneArrays(arrays [][]rune) *Layout {

	r := &Layout{Width: len(arrays[0]), Height: len(arrays)}
	r.Data = [][]rune{}
	for y := 0; y < len(arrays); y++ {
		r.Data = append(r.Data, []rune{})
		for x := 0; x < len(arrays[0]); x++ {
			r.Data[y] = append(r.Data[y], arrays[y][x])
		}
	}

	return r

}

// NewLayoutFromStringArray creates a new Map with the data contained in the provided string array.
func NewLayoutFromStringArray(array []string) *Layout {

	runes := [][]rune{}

	for _, str := range array {
		runes = append(runes, []rune(str))
	}

	return NewLayoutFromRuneArrays(runes)

}

type BSPOptions struct {
	WallValue       rune // Rune value to use for walls
	SplitCount      int  // How many times to split the layout
	DoorValue       rune // Rune value to use for doors / doorways
	MinimumRoomSize int  // Minimum allowed size of each room within the generated BSP layout
}

func NewDefaultBSPOptions() BSPOptions {

	return BSPOptions{
		WallValue:       'x',
		SplitCount:      10,
		DoorValue:       '#',
		MinimumRoomSize: 4,
	}

}

// BSPRoom represents a room generated through Layout.GenerateBSP().
type BSPRoom struct {
	X, Y, W, H  int        // X, Y, Width, and Height of the BSPRoom.
	Connected   []*BSPRoom // The BSPRooms this room is connected to.
	Traversible bool       // Whether the BSPRoom is traversible when using CountHopsTo().
}

func NewBSPRoom(x, y, w, h int) *BSPRoom {
	return &BSPRoom{
		X:           x,
		Y:           y,
		W:           w,
		H:           h,
		Connected:   []*BSPRoom{},
		Traversible: true,
	}
}

// Area returns the area of the BSPRoom (width * height).
func (bsp *BSPRoom) Area() int {
	return bsp.W * bsp.H
}

// MinSize returns the minimum size of the room.
func (bsp *BSPRoom) MinSize() int {
	if bsp.W < bsp.H {
		return bsp.W
	}
	return bsp.H
}

func (bsp *BSPRoom) Center() Position {
	return Position{bsp.X + bsp.W/2, bsp.Y + bsp.H/2}
}

// CountHopsTo will count the number of hops to go from one room to another, by hopping through connected neighbors. If no traversible link between the two rooms found, CountHopsTo will return -1.
func (bsp *BSPRoom) CountHopsTo(room *BSPRoom) int {

	toCheck := append([]*BSPRoom{}, bsp)
	perRoomHopCount := map[*BSPRoom]int{
		bsp: 0,
	}

	for len(toCheck) > 0 {

		next := toCheck[0]

		if next == room {
			return perRoomHopCount[next]
		}

		toCheck = toCheck[1:]

		if !next.Traversible {
			continue
		}

		for _, connected := range next.Connected {

			if _, exists := perRoomHopCount[connected]; !exists {
				toCheck = append(toCheck, connected)
				perRoomHopCount[connected] = perRoomHopCount[next] + 1
			}

		}

	}

	return -1

}

// Disconnect removes the BSPRoom from any of its neighbors' Connected lists, breaking the link between them.
func (bsp *BSPRoom) Disconnect() {

	for _, neighbor := range bsp.Connected {
		for i, me := range neighbor.Connected {
			if me == bsp {
				neighbor.Connected = append(neighbor.Connected[:i], neighbor.Connected[i+1:]...)
				break
			}
		}
	}

	bsp.Connected = []*BSPRoom{}

}

// Necessary returns if the BSPRoom is necessary to facilitate traversal from its neighbors to the rest of the BSP Layout.
func (bsp *BSPRoom) Necessary() bool {

	// If you only have one neighbor, then you're necessary
	if len(bsp.Connected) == 1 {
		return true
	}

	bsp.Traversible = false

	for _, neighbor := range bsp.Connected {

		// If your neighbor is only connected to you, then you're necessary.
		if len(neighbor.Connected) <= 1 {
			bsp.Traversible = true
			return true
		}

		for _, otherNeighbor := range bsp.Connected {

			if otherNeighbor == bsp || otherNeighbor == neighbor {
				continue
			}

			if neighbor.CountHopsTo(otherNeighbor) < 0 {
				bsp.Traversible = true
				return true
			}

		}

	}

	bsp.Traversible = true
	return false

}

// GenerateBSP generates a map in the given Layout using BSP (binary space partitioning) generation, drawing lines of WallValue runes horizontally and
// vertically across, partitioning the room into pieces. It also will place single cells of doorValue on the walls, creating
// doorways. Link: http://www.roguebasin.com/index.php?title=Basic_BSP_Dungeon_generation
// GenerateBSP works best with an empty Layout.
// GenerateBSP returns a list of rooms generated through this method; these are simply internal data structures used to tell
// where rooms start and stop, and can be used in accordance with Layout.Set() or Layout.Select() to modify these rooms. For BSP
// Generation, walls are always on the X and Y lines of the Layout only; not the right or bottom sides.
func (layout *Layout) GenerateBSP(bspOptions BSPOptions) []*BSPRoom {

	layout.Select().Fill(' ')

	subSplit := func(parent *BSPRoom) (*BSPRoom, *BSPRoom, bool) {

		vertical := rand.Float32() >= 0.5
		if parent.W > parent.H*2 {
			vertical = true
		} else if parent.H > parent.W*2 {
			vertical = false
		}

		splitPercentage := 0.2 + rand.Float32()*0.6

		if vertical {

			splitCX := int(float32(parent.W) * splitPercentage)

			a := NewBSPRoom(parent.X, parent.Y, splitCX, parent.H)
			b := NewBSPRoom(parent.X+splitCX, parent.Y, parent.W-splitCX, parent.H)

			if a.MinSize() <= bspOptions.MinimumRoomSize || b.MinSize() <= bspOptions.MinimumRoomSize {
				return a, b, false
			}

			// Line is attempting to start on a door
			if bspOptions.DoorValue != bspOptions.WallValue && bspOptions.DoorValue != 0 && (layout.Get(parent.X+splitCX, parent.Y) == bspOptions.DoorValue || layout.Get(parent.X+splitCX, parent.Y+parent.H) == bspOptions.DoorValue) {
				return a, b, false
			}

			layout.DrawLine(parent.X+splitCX, parent.Y, parent.X+splitCX, parent.Y+parent.H-1, bspOptions.WallValue, 1, false)

			return a, b, true
		}

		splitCY := int(float32(parent.H) * splitPercentage)

		a := NewBSPRoom(parent.X, parent.Y, parent.W, splitCY)
		b := NewBSPRoom(parent.X, parent.Y+splitCY, parent.W, parent.H-splitCY)

		// We can't split a room too small.
		if a.MinSize() <= bspOptions.MinimumRoomSize || b.MinSize() <= bspOptions.MinimumRoomSize {
			return a, b, false
		}

		// Line is attempting to start on a door
		if bspOptions.DoorValue != bspOptions.WallValue && bspOptions.DoorValue != 0 && (layout.Get(parent.X, parent.Y+splitCY) == bspOptions.DoorValue || layout.Get(parent.X+parent.W, parent.Y+splitCY) == bspOptions.DoorValue) {
			return a, b, false
		}

		layout.DrawLine(parent.X, parent.Y+splitCY, parent.X+parent.W-1, parent.Y+splitCY, bspOptions.WallValue, 1, false)

		return a, b, true

	}

	rooms := []*BSPRoom{
		NewBSPRoom(0, 0, layout.Width, layout.Height),
	}

	if layout.Seed > math.MinInt64 {
		rand.Seed(layout.Seed)
	} else {
		rand.Seed(time.Now().UnixNano())
	}

	splitCount := 0

	i := 0
	for true {

		// Sort the rooms so bigger ones can be prioritized sometimes
		sort.Slice(rooms, func(i, j int) bool {
			// return rooms[i].Area() > rooms[j].Area()
			return rooms[i].MinSize() > rooms[j].MinSize()
		})

		splitChoice := rooms[rand.Intn(len(rooms))]

		if rand.Float32() >= 0.2 {
			splitChoice = rooms[0] // Try to split the biggest rooms first
		}

		// Do the split

		a, b, success := subSplit(splitChoice)

		i++

		if i >= bspOptions.SplitCount*10 { // Assume it's done to avoid just hanging the system
			break
		}

		if !success {
			continue
		} else {

			rooms = append(rooms, a, b)

			for i, r := range rooms {
				if r == splitChoice {
					rooms = append(rooms[:i], rooms[i+1:]...)
					break
				}
			}

		}

		splitCount++

		if splitCount >= bspOptions.SplitCount {
			break
		}

	}

	// Generate Doors
	for _, subroom := range rooms {

		possibleExits := []Position{}

		spawnOptions := []int{0, 1, 2}

		spawnChoice := spawnOptions[rand.Intn(len(spawnOptions))]

		// Rooms on the border must generate a doorway that works for them
		if subroom.X == 0 || subroom.Y == 0 {
			spawnChoice = 2
		}

		// Spawn both directions; Y-axis
		if spawnChoice == 1 || spawnChoice == 2 {

			if subroom.Y > 0 {

				for x := subroom.X; x < subroom.X+subroom.W; x++ {

					up := layout.Get(x, subroom.Y-1)
					down := layout.Get(x, subroom.Y+1)

					if up == ' ' && down == ' ' {
						possibleExits = append(possibleExits, Position{x, subroom.Y})
					}

				}

				doorway := possibleExits[rand.Intn(len(possibleExits))]

				layout.Set(doorway.X, doorway.Y, bspOptions.DoorValue)

				doorRect := image.Rect(doorway.X, doorway.Y-1, doorway.X+1, doorway.Y)

				for _, other := range rooms {

					otherRect := image.Rect(other.X, other.Y, other.X+other.W, other.Y+other.H)

					if otherRect.Overlaps(doorRect) {
						other.Connected = append(other.Connected, subroom)
						subroom.Connected = append(subroom.Connected, other)
					}

				}

			}

		}

		possibleExits = []Position{}

		// X-axis doorways
		if spawnChoice == 0 || spawnChoice == 2 {

			if subroom.X > 0 {

				for y := subroom.Y; y < subroom.Y+subroom.H; y++ {

					left := layout.Get(subroom.X-1, y)
					right := layout.Get(subroom.X+1, y)

					if left == ' ' && right == ' ' {
						possibleExits = append(possibleExits, Position{subroom.X, y})
					}

				}

				doorway := possibleExits[rand.Intn(len(possibleExits))]

				layout.Set(doorway.X, doorway.Y, bspOptions.DoorValue)

				for _, other := range rooms {

					otherRect := image.Rect(other.X, other.Y, other.X+other.W, other.Y+other.H)
					if otherRect.Overlaps(image.Rect(doorway.X-1, doorway.Y, doorway.X, doorway.Y+1)) {
						other.Connected = append(other.Connected, subroom)
						subroom.Connected = append(subroom.Connected, other)
					}

				}

			}

		}

	}

	return rooms

}

// GenerateRandomRooms generates a map using random room creation. emptyRune is the rune to fill the rooms generated with, while wallRune is the rune to use as walls (unwalkable tiles).
// roomCount is how many rooms to place, roomMinWidth and Height are how small they can be, minimum, while roomMaxWidth and Height are how large
// they can be. connectRooms determines if the algorithm should also attempt to connect the rooms using pathways between each room. The
// function returns the positions of each room created.
func (layout *Layout) GenerateRandomRooms(emptyRune rune, wallRune rune, roomCount, roomMinWidth, roomMinHeight, roomMaxWidth, roomMaxHeight int, connectRooms bool) [][]int {

	layout.Select().Fill(wallRune)

	if layout.Seed > math.MinInt64 {
		rand.Seed(layout.Seed)
	} else {
		rand.Seed(time.Now().UnixNano())
	}

	roomPositions := make([][]int, 0)

	for i := 0; i < roomCount; i++ {

		// roomSize := float64(2 + rand.Intn(2))

		sx := rand.Intn(layout.Width)
		sy := rand.Intn(layout.Height)

		roomPositions = append(roomPositions, []int{sx, sy})

		roomW := roomMinWidth + rand.Intn(roomMaxWidth-roomMinWidth)
		roomH := roomMinHeight + rand.Intn(roomMaxHeight-roomMinHeight)

		drawRoom := func(x, y int) bool {
			dx := int(math.Abs(float64(sx) - float64(x)))
			dy := int(math.Abs(float64(sy) - float64(y)))
			if dx < roomW && dy < roomH {
				layout.Set(x, y, emptyRune)
			}
			return true
		}

		layout.Select().FilterBy(drawRoom)

	}

	if connectRooms {

		for p := 0; p < len(roomPositions); p++ {

			if p < len(roomPositions)-1 {

				x := roomPositions[p][0]
				y := roomPositions[p][1]

				x2 := roomPositions[p+1][0]
				y2 := roomPositions[p+1][1]

				layout.DrawLine(x, y, x2, y2, emptyRune, 1, true)

			}

		}

	}

	return roomPositions

}

// GenerateDrunkWalk generates a map in the bounds of the Layout specified using drunk walking. It will pick a random point in the
// Layout and begin walking around at random, placing fillRune in the Layout, until at least percentageFilled (0.0 - 1.0) of the Layout
// is filled. Note that it only counts values placed in the cell, not instances where it moves over a cell that already has the
// value being placed. This can be used to generate maps more similar to simple natural cave systems, as an imaginary example.
// Link: http://www.roguebasin.com/index.php?title=Random_Walk_Cave_Generation
func (layout *Layout) GenerateDrunkWalk(emptyRune rune, wallRune rune, percentageFilled float32) (X, Y int) {

	selection := layout.Select()
	selection.Fill(wallRune)

	if layout.Seed > math.MinInt64 {
		rand.Seed(layout.Seed)
	} else {
		rand.Seed(time.Now().UnixNano())
	}

	sx := rand.Intn(layout.Width)
	sy := rand.Intn(layout.Height)

	startX := sx
	startY := sy

	fillCount := float32(0)

	totalArea := float32(layout.Area())

	for true {

		cell := selection.FilterByArea(sx, sy, 2, 2)

		if cell.ContainsChar(wallRune) {
			cell.Fill(emptyRune)
			fillCount += 4
		}

		dir := rand.Intn(4)

		if dir == 0 {
			sx++
		} else if dir == 1 {
			sx--
		} else if dir == 2 {
			sy++
		} else if dir == 3 {
			sy--
		}

		if sx < 0 {
			sx = 0
		} else if sx >= layout.Width {
			sx = layout.Width - 1
		}

		if sy < 0 {
			sy = 0
		} else if sy >= layout.Height {
			sy = layout.Height - 1
		}

		if fillCount/totalArea >= percentageFilled {
			break
		}

	}

	return startX, startY

}

// Rotate rotates the entire room 90 degrees clockwise.
func (layout *Layout) Rotate() {

	newData := make([][]rune, 0)

	for y := 0; y < len(layout.Data[0]); y++ {
		newData = append(newData, []rune{})
		for x := 0; x < len(layout.Data); x++ {
			nx := (layout.Height - x) - 1
			newData[y] = append(newData[y], layout.Data[nx][y])
		}
	}

	layout.Data = newData
	layout.Height = len(layout.Data)
	layout.Width = len(layout.Data[0])

}

// CopyFrom copies the data from the other Layout into this Layout's data. x and y are the position of the other Layout's data in the
// destination (calling) Layout.
func (layout *Layout) CopyFrom(other *Layout, x, y int) {

	for cy := 0; cy < layout.Height; cy++ {
		for cx := 0; cx < layout.Width; cx++ {
			if cx >= x && cy >= y && cx-x < other.Width && cy-y < other.Height {
				layout.Set(cx, cy, other.Get(cx-x, cy-y))
			}
		}
	}

}

// DrawLine is used to draw a line from x, y, to x2, y2, placing the rune specified by fillRune in the cells between those points (including)
// in those points themselves, as well. thickness controls how thick the line is. If stagger is on, then the line will stagger it's
// vertical movement, allowing a 1-thickness line to actually be pass-able if an object was only able to move in cardinal directions
// and the line had a diagonal slope.
func (layout *Layout) DrawLine(x, y, x2, y2 int, fillRune rune, thickness int, stagger bool) {

	dx := int(math.Abs(float64(x2 - x)))
	dy := int(math.Abs(float64(y2 - y)))
	slope := float32(0)
	xAxis := true

	if dx != 0 {
		slope = float32(dy) / float32(dx)
	}
	length := dx

	if dy > dx {
		xAxis = false
		if dy != 0 {
			slope = float32(dx) / float32(dy)
		}
		length = dy
	}

	sx := float32(x)
	sy := float32(y)

	set := func(x, y int) {
		for fx := 0; fx < thickness; fx++ {
			for fy := 0; fy < thickness; fy++ {
				layout.Set(x+fx-thickness/2, y+fy-thickness/2, fillRune)
			}
		}
	}

	for c := 0; c < length+1; c++ {

		set(int(math.Round(float64(sx))), int(math.Round(float64(sy))))

		mx := int(math.Round(float64(sx)))

		if xAxis {
			if x2 > x {
				sx++
			} else {
				sx--
			}
			if y2 > y {
				sy += slope
			} else {
				sy -= slope
			}
		} else {
			if y2 > y {
				sy++
			} else {
				sy--
			}
			if x2 > x {
				sx += slope
			} else {
				sx -= slope
			}
		}

		if stagger {
			set(mx, int(math.Round(float64(sy))))
		}

	}

}

// Set sets the rune provided in the Layout's Data. If the position given is outside of the Layout, no rune will be set.
func (layout *Layout) Set(x, y int, char rune) {

	if x < 0 || x >= layout.Width || y < 0 || y >= layout.Height {
		return
	}

	layout.Data[y][x] = char

}

// Get returns the rune in the specified position in the Layout's Data array. If the position given goes outside of the bounds of the Layout, it will return a null rune (0).
func (layout *Layout) Get(x, y int) rune {

	if x < 0 || x >= layout.Width || y < 0 || y >= layout.Height {
		return 0
	}

	return layout.Data[y][x]
}

// ClosestChar returns the position of the closest rune to the given x and y position with the value of char. If no character is found
// in the Layout, the position of -1, -1 is returned.
func (layout *Layout) ClosestChar(x, y int, char rune) Position {

	points := []Position{}

	closest := Position{-1, -1}

	for y := 0; y < layout.Height; y++ {
		for x := 0; x < layout.Width; x++ {
			if layout.Get(x, y) == char {
				points = append(points, Position{x, y})
			}
		}
	}

	sort.Slice(points, func(i, j int) bool {
		return Position{x, y}.DistanceTo(points[i]) < Position{x, y}.DistanceTo(points[j])
	})

	if len(points) > 0 {
		closest = points[0]
	}

	return closest

}

// Center returns the center position of the Layout.
func (layout *Layout) Center() Position {
	return Position{layout.Width / 2, layout.Height / 2}
}

// Resize resizes the room to be of the width and height provided. Note that resizing to a smaller Layout is destructive (and so,
// data will be lost if resizing to a smaller Layout).
func (layout *Layout) Resize(width, height int) *Layout {

	layout.Width = width
	layout.Height = height

	data := make([][]rune, 0)

	for y := 0; y < height; y++ {

		data = append(data, []rune{})

		for x := 0; x < width; x++ {

			if len(layout.Data) > y && len(layout.Data[y]) > x {
				data[y] = append(data[y], layout.Get(x, y))
			} else {
				data[y] = append(data[y], 0)
			}

		}

	}

	layout.Data = data

	return layout

}

// Area returns the overall size of the Layout by multiplying the width by the height.
func (layout *Layout) Area() int {
	return layout.Width * layout.Height
}

// MinimumSize returns the minimum distance (width or height) for the Layout.
func (layout *Layout) MinimumSize() int {
	if layout.Width < layout.Height {
		return layout.Width
	}
	return layout.Height
}

// DataToString returns the underlying data of the overall Layout layout in an easily understood visual format.
// 0's turn into blank spaces when using DataToString, and the column is shown at the left of the map.
func (layout *Layout) DataToString() string {

	s := fmt.Sprintf("  W:%d H:%d\n\n       ", layout.Width, layout.Height)

	for y := 1; y < len(layout.Data[0]); y += 2 {
		s += fmt.Sprintf("%2d  ", y)
	}
	s += "\n     "
	for y := 0; y < len(layout.Data[0]); y += 2 {
		s += fmt.Sprintf("%2d  ", y)
	}

	s += "\n"

	for y := 0; y < len(layout.Data); y++ {
		s += fmt.Sprintf("%3d  |", y)
		for x := 0; x < len(layout.Data[y]); x++ {
			// s += " " + string(room.Data[y][x])
			s += fmt.Sprintf("%v ", string(layout.Data[y][x]))
		}
		s += "|\n"
	}

	return s

}

// Select returns a filled Selection of a Layout.
func (layout *Layout) Select() Selection {
	newSelection := Selection{
		Layout: layout,
		Cells:  map[Position]bool{},
	}
	return newSelection.All()
}

// SelectContiguous creates a Selection from all cells contiguous to the cell in the (x,y) position provided.
func (layout *Layout) SelectContiguous(x, y int, diagonal bool) Selection {

	toAdd := []Position{
		{x, y},
	}

	added := map[Position]bool{
		{x, y}: true,
	}

	startingValue := layout.Get(x, y)

	for len(toAdd) > 0 {

		position := toAdd[0]

		added[position] = true

		sides := []Position{
			{position.X - 1, position.Y},
			{position.X + 1, position.Y},
			{position.X, position.Y - 1},
			{position.X, position.Y + 1},
		}

		if diagonal {
			sides = append(sides,
				Position{position.X - 1, position.Y - 1},
				Position{position.X + 1, position.Y + 1},
				Position{position.X + 1, position.Y - 1},
				Position{position.X - 1, position.Y + 1},
			)
		}

		for _, side := range sides {

			if layout.Get(side.X, side.Y) == startingValue && !added[side] {
				toAdd = append(toAdd, side)
				added[side] = true
			}

		}

		toAdd = toAdd[1:]

	}

	newSelection := layout.Select()

	newSelection.Cells = added

	return newSelection

}

// Position represents a cell within a Selection.
type Position struct {
	X, Y int
}

// DistanceTo returns the distance from one Position to another.
func (position Position) DistanceTo(other Position) float64 {
	return math.Sqrt(float64(math.Pow(float64(other.X-position.X), 2) + math.Pow(float64(other.Y-position.Y), 2)))
}

// A Selection represents a selection of cell positions in the Layout's data array, and can be filtered down and manipulated
// using the functions on the Selection struct. You can use Selections to manipulate a
type Selection struct {
	Layout *Layout
	Cells  map[Position]bool
}

func (selection *Selection) Clone() Selection {
	newSelection := Selection{
		Layout: selection.Layout,
		Cells:  map[Position]bool{},
	}
	for key := range selection.Cells {
		newSelection.Cells[key] = true
	}
	return newSelection
}

// FilterByRune filters the Selection down to the cells that have the character (rune) provided.
func (selection Selection) FilterByRune(value rune) Selection {
	return selection.FilterBy(func(x, y int) bool {
		return selection.Layout.Get(x, y) == value
	})
}

// All returns a selection with all cells from the Layout selected.
func (selection Selection) All() Selection {
	newSelection := selection.Clone()
	for y := 0; y < len(newSelection.Layout.Data); y++ {
		for x := 0; x < len(newSelection.Layout.Data[y]); x++ {
			newSelection.Cells[Position{x, y}] = true
		}
	}
	return newSelection
}

// None returns a selection with no selected cells from the Layout.
func (selection Selection) None() Selection {
	newSelection := Selection{
		Layout: selection.Layout,
		Cells:  map[Position]bool{},
	}
	return newSelection
}

// FilterByPercentage selects the provided percentage (from 0 - 1) of the cells curently in the Selection.
func (selection Selection) FilterByPercentage(percentage float32) Selection {

	return selection.FilterBy(func(x, y int) bool {

		if rand.Float32() <= percentage {
			return true
		}
		return false
	})

}

// FilterByArea filters down a selection by only selecting the cells that have X and Y values between X, Y, and X+W and Y+H.
// It crops the selection, basically.
func (selection Selection) FilterByArea(x, y, w, h int) Selection {

	return selection.FilterBy(func(cx, cy int) bool {
		return cx >= x && cy >= y && cx <= x+w-1 && cy <= y+h-1
	})

}

// Add returns a clone of the current Selection with the cells in the other Selection.
func (selection Selection) Add(other Selection) Selection {

	newSelection := selection.Clone()

	for position := range other.Cells {
		newSelection.Cells[position] = true
	}

	return newSelection

}

// Remove returns a clone of the current Selection without the cells in the other Selection.
func (selection Selection) Remove(other Selection) Selection {

	newSelection := selection.Clone()

	for position := range other.Cells {
		delete(newSelection.Cells, position)
	}

	return newSelection

}

// FilterByNeighbor returns a filtered Selection of the cells that are surrounded at least by minNeighborCount neighbors with a value of
// neighborValue. If diagonals is true, then diagonals are also checked. If atMost is true, then FilterByNeighbor will only
// work if there's at MOST that many neighbors.
func (selection Selection) FilterByNeighbor(neighborValue rune, minNeighborCount int, diagonals bool, atMost bool) Selection {

	return selection.FilterBy(func(x, y int) bool {

		n := 0

		if selection.Layout.Get(x-1, y) == neighborValue {
			n++
		}
		if selection.Layout.Get(x+1, y) == neighborValue {
			n++
		}
		if selection.Layout.Get(x, y-1) == neighborValue {
			n++
		}
		if selection.Layout.Get(x, y+1) == neighborValue {
			n++
		}

		if diagonals {
			if selection.Layout.Get(x-1, y-1) == neighborValue {
				n++
			}
			if selection.Layout.Get(x+1, y-1) == neighborValue {
				n++
			}
			if selection.Layout.Get(x-1, y+1) == neighborValue {
				n++
			}
			if selection.Layout.Get(x+1, y+1) == neighborValue {
				n++
			}
		}

		if atMost {
			return n <= minNeighborCount
		}

		return n >= minNeighborCount

	})

}

// FilterBy takes a function that takes the X and Y values of each cell position contained in the Selection, and returns a
// boolean to indicate whether to include that cell in the Selection or not. If the result is true, the cell is included in the Selection;
// Otherwise, it is filtered out. This allows you to easily make custom filtering functions to filter down the cells in a Selection.
func (selection Selection) FilterBy(filterFunc func(x, y int) bool) Selection {

	// Note that while we're assigning the Cells variable of selection here directly,
	// because this function doesn't take a pointer notation, we're operating on a copy
	// of the selection, not the original.

	newSelection := selection.Clone()

	cells := map[Position]bool{}

	for c := range newSelection.Cells {
		if filterFunc(c.X, c.Y) {
			cells[c] = true
		}
	}

	newSelection.Cells = cells

	return newSelection

}

// Select attempts to select a number of the cells contained within the selection and returns them. If there's fewer cells in the selection,
// then it will simply return the entirety of the selection.
func (selection Selection) Select(num int) []Position {

	cells := []Position{}

	for cell := range selection.Cells {
		cells = append(cells, cell)
		if len(cells) >= num {
			return cells
		}
	}

	return cells

}

// Expand expands the selection outwards by the distance value provided. Diagonal indicates if the expansion should happen
// diagonally as well, or just on the cardinal 4 directions. If a negative value is given for distance, it shrinks the selection.
func (selection Selection) Expand(distance int, diagonal bool) Selection {

	newSelection := selection.Clone()

	if distance == 0 {
		return newSelection
	}

	shrinking := false
	if distance < 0 {
		shrinking = true
		distance *= -1
	}

	for i := 0; i < distance; i++ {

		// We can't loop through the cells while modifying them, so we'll make a copy after each iteration.
		cells := map[Position]bool{}

		for c := range newSelection.Cells {
			cells[c] = true
		}

		toRemove := []Position{}

		for cp := range cells {

			if shrinking {

				if !newSelection.Contains(cp.X-1, cp.Y) || !newSelection.Contains(cp.X+1, cp.Y) || !newSelection.Contains(cp.X, cp.Y-1) || !newSelection.Contains(cp.X, cp.Y+1) {

					if !diagonal || !newSelection.Contains(cp.X-1, cp.Y-1) || !newSelection.Contains(cp.X+1, cp.Y-1) || !newSelection.Contains(cp.X-1, cp.Y+1) || !newSelection.Contains(cp.X+1, cp.Y+1) {

						toRemove = append(toRemove, cp)

					}

				}

			} else {

				newSelection.AddPosition(cp.X-1, cp.Y)
				newSelection.AddPosition(cp.X+1, cp.Y)
				newSelection.AddPosition(cp.X, cp.Y-1)
				newSelection.AddPosition(cp.X, cp.Y+1)

				if diagonal {
					newSelection.AddPosition(cp.X-1, cp.Y-1)
					newSelection.AddPosition(cp.X-1, cp.Y+1)
					newSelection.AddPosition(cp.X+1, cp.Y-1)
					newSelection.AddPosition(cp.X+1, cp.Y+1)
				}

			}

		}

		for _, c := range toRemove {
			newSelection.RemovePosition(c.X, c.Y)
		}

	}

	return newSelection

}

// Invert inverts the selection (selects all non-selected cells from the Selection's source Map).
func (selection Selection) Invert() Selection {

	inverted := selection.Layout.Select()

	return inverted.FilterBy(func(x, y int) bool {
		return !selection.Contains(x, y)
	})

}

// Contains returns a boolean indicating if the specified cell is in the list of cells contained in the selection.
func (selection *Selection) Contains(x, y int) bool {
	for c := range selection.Cells {
		if c.X == x && c.Y == y {
			return true
		}
	}
	return false
}

func (selection *Selection) ContainsChar(char rune) bool {
	for c := range selection.Cells {
		if selection.Layout.Get(c.X, c.Y) == char {
			return true
		}
	}
	return false
}

// Fill fills the cells in the Selection with the rune provided.
func (selection Selection) Fill(char rune) Selection {
	return selection.FilterBy(func(x, y int) bool {
		selection.Layout.Set(x, y, char)
		return true
	})
}

// AddPosition adds a specific position to the Selection. If the position lies outside of the layout's area, then it's removed.
func (selection *Selection) AddPosition(x, y int) {

	if x < 0 || y < 0 || x >= selection.Layout.Width || y >= selection.Layout.Height {
		return
	}

	selection.Cells[Position{x, y}] = true

}

// RemovePosition removes a specific position from the Selection.
func (selection *Selection) RemovePosition(x, y int) {
	delete(selection.Cells, Position{x, y})
}
