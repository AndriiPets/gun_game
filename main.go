package main

import (
	//"fmt"
	"image"
	"log"

	//"time"

	"github.com/AndriiPets/FishGame/config"
	"github.com/AndriiPets/FishGame/scenes"
	"github.com/hajimehoshi/ebiten/v2"
)

type Scene interface {
	Update()
	Draw(screen *ebiten.Image)
}

type Game struct {
	bounds image.Rectangle
	scene  Scene
}

func NewGame() *Game {
	g := &Game{
		bounds: image.Rectangle{},
		scene:  &scenes.MainScene{},
	}

	//go func() {

	//	for {

	//		fmt.Println("FPS: ", ebiten.ActualFPS())
	//		fmt.Println("Ticks: ", ebiten.ActualTPS())
	//		time.Sleep(time.Second)

	//	}

	//}()

	return g
}

func (g *Game) Update() error {
	g.scene.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	g.scene.Draw(screen)
}

func (g *Game) Layout(width, height int) (int, int) {
	g.bounds = image.Rect(0, 0, width, height)
	return width, height
}

func main() {
	ebiten.SetWindowSize(config.C.ScreenWidth, config.C.ScreenHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
