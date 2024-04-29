package factory

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/math"

	"github.com/AndriiPets/FishGame/archetypes"
	"github.com/AndriiPets/FishGame/components"
	"github.com/AndriiPets/FishGame/config"
)

func CreateCamera(ecs *ecs.ECS) *donburi.Entry {
	camera := archetypes.Camera.Spawn(ecs)

	components.Camera.SetValue(camera, components.CameraData{
		ViewPort: math.NewVec2(float64(config.C.ScreenWidth), float64(config.C.ScreenHeight)),
		Position: math.NewVec2(32, 128),
	})

	return camera
}
