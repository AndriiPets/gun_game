package systems

import (
	"fmt"
	"image/color"
	"math"

	"github.com/AndriiPets/FishGame/components"
	dresolv "github.com/AndriiPets/FishGame/resolv"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi/ecs"
	dmath "github.com/yohamta/donburi/features/math"
)

// var GeoM ebiten.GeoM
var (
	cam       *components.CameraData
	playerPos dmath.Vec2
	delta     float64
)

func CameraUpdate(ecs *ecs.ECS) {

	delta = ecs.Time.DeltaTime().Seconds()

	cameraEntity, ok := components.Camera.First(ecs.World)

	if ok != true {
		panic("no camera!")
	}

	playerEntity, _ := components.Player.First(ecs.World)
	playerObj := dresolv.GetObject(playerEntity)
	camera := components.Camera.Get(cameraEntity)

	cam = camera

	playerPos = dmath.NewVec2(playerObj.Position.X, playerObj.Position.Y)

	//translate cursor position on the screen to position in the world
	mouseX, mouseY := ScreenToWorld(ebiten.CursorPosition())

	camera.CursorX = mouseX
	camera.CursorY = mouseY

	reset()

}

func CameraString(ecs *ecs.ECS) string {
	cameraEntity, _ := components.Camera.First(ecs.World)
	c := components.Camera.Get(cameraEntity)
	return fmt.Sprintf(
		"Camera T: %.1f, R: %d, S: %.5f",
		c.Position, c.Rotation, delta,
	)
}

func CameraRender(world, screen *ebiten.Image) {
	screen.DrawImage(world, &ebiten.DrawImageOptions{
		GeoM: WorldMatrix(cam),
	})

	//draw aim circle
	mouseX, mouseY := ebiten.CursorPosition()
	mx, my := float32(mouseX), float32(mouseY)

	aimColor := color.RGBA{0, 225, 0, 225}

	vector.StrokeCircle(screen, mx, my, 12, 2, aimColor, false)
	vector.DrawFilledCircle(screen, mx, my, 2, aimColor, false)

	if cam.Flash {
		//screen.Fill(color.RGBA{225, 225, 225, 55})
		cam.Flash = false
	}
}

func WorldMatrix(c *components.CameraData) ebiten.GeoM {

	ViewPortCenter := dmath.NewVec2(c.ViewPort.X*0.5, c.ViewPort.Y*0.5)

	//camera smooth follow calculation
	minSpeed := 20.0
	minEffectLen := 3.0
	fractionSpeed := 2.3

	diff := playerPos.Sub(c.Position)
	len := diff.Magnitude()

	if len > minEffectLen {
		speed := math.Max(fractionSpeed*len, minSpeed)
		c.Position = c.Position.Add(diff.MulScalar(speed * delta / len))
	}

	m := ebiten.GeoM{}
	m.Translate(-c.Position.X-c.Recoil.X, -c.Position.Y-c.Recoil.Y) //target
	//m.Translate(-c.Recoil.X, -c.Recoil.Y)     //recoil

	//zoom
	m.Scale(
		math.Pow(1.01, float64(c.Zoom)),
		math.Pow(1.01, float64(c.Zoom)),
	)

	m.Translate(ViewPortCenter.X, ViewPortCenter.Y) //offset

	return m
}

func reset() {
	if cam.Zoom < 0 {
		cam.Zoom += 1
	}

	if !cam.Recoil.IsZero() {
		cam.Recoil = cam.Recoil.DivScalar(2)
	}
}

func ScreenToWorld(posX, posY int) (float64, float64) {

	invMatrix := WorldMatrix(cam)
	//GeoM = invMatrix
	if invMatrix.IsInvertible() {

		invMatrix.Invert()
		return invMatrix.Apply(float64(posX), float64(posY))

	} else {

		return math.NaN(), math.NaN()

	}
}
