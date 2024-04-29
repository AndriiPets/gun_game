package systems

import (
	"image/color"

	dresolv "github.com/AndriiPets/FishGame/resolv"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"

	"github.com/AndriiPets/FishGame/tags"
)

func DrawWall(ecs *ecs.ECS, image *ebiten.Image) {
	tags.Wall.Each(ecs.World, func(e *donburi.Entry) {
		o := dresolv.GetObject(e)
		drawColor := color.RGBA{60, 60, 60, 255}
		vector.DrawFilledRect(image, float32(o.Position.X), float32(o.Position.Y), float32(o.Size.X), float32(o.Size.Y), drawColor, false)
	})
}
