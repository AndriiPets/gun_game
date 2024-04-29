package assets

import (
	"bytes"
	"fmt"
	"image"
	_ "image/png"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/ganim8/v2"
)

type sprite struct {
	File string `json:"file"`
	W    int    `json:"w"`
	H    int    `json:"h"`
}

type animation struct {
	File   string        `json:"file"`
	Name   string        `json:"name"`
	Frames []interface{} `json:"frames"`
	Flip   bool
}

type spriteConfig struct {
	Sprites    []sprite    `json:"sprites"`
	Animations []animation `json:"animations"`
}

var (
	grids      = make(map[string]*ganim8.Grid)
	images     = make(map[string]*ebiten.Image)
	sprites    = make(map[string]*ganim8.Sprite)
	animations = make(map[string]*ganim8.Animation)
)

func GetSprite(name string) *ganim8.Sprite {
	if _, ok := sprites[name]; !ok {
		panic(fmt.Sprintf("sprite not found: %s", name))
	}

	return sprites[name]
}

func GetAnimation(name string) *ganim8.Animation {
	if _, ok := animations[name]; !ok {
		panic(fmt.Sprintf("animation not found: %s", name))
	}

	return animations[name].Clone()
}

func loadSprites(cfg *spriteConfig) {
	for _, s := range cfg.Sprites {
		b := mustRead(s.File)                             //load from file
		img := ebiten.NewImageFromImage(*decodeImage(&b)) //convert to ebiten image

		images[s.File] = img
		size := img.Bounds().Size()
		g := ganim8.NewGrid(s.W, s.H, size.X, size.Y) //create grid from image
		grids[s.File] = g

		spr := ganim8.NewSprite(img, g.Frames()) //create sprites from grid
		sprites[s.File] = spr
	}
}

func loadAnimations(cfg *spriteConfig) {
	for _, a := range cfg.Animations {
		g, ok := grids[a.File]
		if !ok {
			panic(fmt.Sprintf("grid not found: %s", a.File))
		}

		img, ok := images[a.File]
		if !ok {
			panic(fmt.Sprintf("image not found: %s", a.File))
		}

		// create sprite for the specified frames
		spr := ganim8.NewSprite(img, g.GetFrames(a.Frames...))

		//create animation
		anim := ganim8.NewAnimation(spr, time.Millisecond*60)
		animations[a.Name] = anim
	}
}

func decodeImage(b *[]byte) *image.Image {
	img, _, _ := image.Decode(bytes.NewReader(*b))
	return &img
}
