package factory

import (
	"github.com/AndriiPets/FishGame/archetypes"
	"github.com/AndriiPets/FishGame/components"
	"github.com/AndriiPets/FishGame/config"

	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func CreateSpace(ecs *ecs.ECS) *donburi.Entry {
	space := archetypes.Space.Spawn(ecs)

	cfg := config.C
	spaceData := resolv.NewSpace(cfg.WorldWidth, cfg.WorldHeigth, 16, 16)
	components.Space.Set(space, spaceData)

	return space
}
