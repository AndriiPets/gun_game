package components

import (
	"github.com/AndriiPets/FishGame/assets"
	"github.com/yohamta/donburi"
	"github.com/yohamta/ganim8/v2"
)

type BlockData struct {
	Type         BlockType
	Destructable bool
}

type BlockType string

var (
	BlockWall  BlockType = "wall"
	BlockFloor BlockType = "floor"
)

var Block = donburi.NewComponentType[BlockData]()

func (b *BlockData) String() string {
	return string(b.Type)
}

func (b *BlockData) Animation() *ganim8.Animation {
	return assets.GetAnimation(b.String() + "_default")
}
