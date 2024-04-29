package factory

import (
	"github.com/AndriiPets/FishGame/archetypes"
	"github.com/AndriiPets/FishGame/components"
	dresolv "github.com/AndriiPets/FishGame/resolv"

	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func CreateWall(ecs *ecs.ECS, obj *resolv.Object, blockType components.BlockType) *donburi.Entry {
	wall := archetypes.Wall.Spawn(ecs)
	block := components.Block.Get(wall)
	animation := components.Animation.Get(wall)

	//setup type and sprite
	switch blockType {
	case components.BlockWall:
		block.Type = components.BlockWall
		obj.AddTags("solid")
	case components.BlockFloor:
		block.Type = components.BlockFloor
	default:
		block.Type = components.BlockFloor

	}

	anim := block.Animation()
	animation.Animation = anim
	animation.Type = components.AnimationStatic

	dresolv.SetObject(wall, obj)

	return wall
}
