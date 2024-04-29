package systems

import (
	"image/color"

	dresolv "github.com/AndriiPets/FishGame/resolv"

	"github.com/AndriiPets/FishGame/tags"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func DrawBullet(ecs *ecs.ECS, image *ebiten.Image) {
	tags.Bullet.Each(ecs.World, func(e *donburi.Entry) {
		o := dresolv.GetObject(e)

		vector.DrawFilledCircle(image, float32(o.Position.X), float32(o.Position.Y), 4, color.RGBA{225, 225, 225, 1}, false)
	})
}
