package scenes

import (
	//"fmt"
	"fmt"
	"image/color"
	"sync"

	"github.com/AndriiPets/FishGame/assets"
	"github.com/AndriiPets/FishGame/components"
	"github.com/AndriiPets/FishGame/config"
	"github.com/AndriiPets/FishGame/events"
	"github.com/AndriiPets/FishGame/factory"
	"github.com/AndriiPets/FishGame/layers"
	dresolv "github.com/AndriiPets/FishGame/resolv"
	"github.com/AndriiPets/FishGame/systems"
	"github.com/AndriiPets/FishGame/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type MainScene struct {
	ecs         *ecs.ECS
	once        sync.Once
	WorldScreen *ebiten.Image
	Time        *ecs.Time
}

func (ms *MainScene) Update() {
	ms.once.Do(ms.configure)
	ms.ecs.Update()
	ms.Time.Update()
}

func (ms *MainScene) Draw(screen *ebiten.Image) {
	ms.WorldScreen.Fill(color.RGBA{194, 178, 128, 255})
	ms.ecs.Draw(ms.WorldScreen)

	//fmt.Println(systems.CameraString(ms.ecs))
	//fmt.Println(systems.PlayerString(ms.ecs))
	systems.CameraRender(ms.WorldScreen, screen)
}

func (ms *MainScene) configure() {

	ms.Time = ecs.NewTime()

	ecs := ecs.NewECS(donburi.NewWorld())

	loadAssets()

	factory.CreateCamera(ecs)
	ms.WorldScreen = ebiten.NewImage(config.C.WorldWidth, config.C.WorldHeigth)

	events.SetupEvents(ecs)

	ecs.AddSystem(systems.UpdateObjects)
	ecs.AddSystem(systems.UpdatePlayer)
	ecs.AddSystem(systems.UpdateAttackVector)
	ecs.AddSystem(systems.UpdateCollisions)
	ecs.AddSystem(systems.CameraUpdate)
	ecs.AddSystem(systems.UpdateShooters)
	ecs.AddSystem(systems.UpdateDespawnable)
	ecs.AddSystem(systems.UpdateAnimations)
	ecs.AddSystem(systems.UpdateWeaponSprite)
	ecs.AddSystem(systems.UpdateHealth)
	ecs.AddSystem(systems.UpdateEnemies)
	ecs.AddSystem(systems.UpdateAI)
	ecs.AddSystem(systems.UpdateParticles)

	ecs.AddSystem(events.UpdateEvents)

	ecs.AddSystem(systems.UpdateSettings)

	//ecs.AddRenderer(layers.Default, systems.DrawWall)
	ecs.AddRenderer(layers.Default, systems.DrawPlayer)
	//ecs.AddRenderer(layers.Default, systems.DrawWeaponFlash)
	ecs.AddRenderer(layers.Default, systems.DrawEnemyes)
	ecs.AddRenderer(layers.Default, systems.DrawDebugAi)
	//ecs.AddRenderer(layers.Default, systems.DrawBullet)

	//Draw animations for each layer
	ecs.AddRenderer(layers.Player, systems.DrawAnimation(layers.Player))
	ecs.AddRenderer(layers.Actors, systems.DrawAnimation(layers.Actors))
	ecs.AddRenderer(layers.Architecture, systems.DrawAnimation(layers.Architecture))
	ecs.AddRenderer(layers.Interactables, systems.DrawAnimation(layers.Interactables))
	ecs.AddRenderer(layers.FX, systems.DrawAnimation(layers.FX))
	ecs.AddRenderer(layers.System, systems.DrawDebug)
	//

	ms.ecs = ecs

	world := utils.NewWorldMap()
	world.GenerateMap(utils.BSP)

	//gw, gh := float64(config.C.WorldWidth), float64(config.C.WorldHeigth)

	space := factory.CreateSpace(ms.ecs)

	for y, row := range world.Map.Data {
		for x, val := range row {
			posX, posY := x*config.BlockSize, y*config.BlockSize //works
			//fmt.Println(posX, posY)
			//var block *donburi.Entry
			if val == 'x' {
				dresolv.Add(space, factory.CreateWall(ms.ecs, resolv.NewObject(float64(posX), float64(posY), float64(config.BlockSize), float64(config.BlockSize)), components.BlockWall))
			}
			if val == 'e' {
				dresolv.Add(space, factory.CreateEnemy(ms.ecs, float64(posX), float64(posY), components.EnemyTypeGrunt))
			}
			if val == 'P' {
				fmt.Println(posX, posY)
				dresolv.Add(space, factory.CreatePlayer(ms.ecs, float64(posX), float64(posY)))
			}

		}
	}

	//dresolv.Add(space,
	//	factory.CreateWall(ms.ecs, resolv.NewObject(0, 0, 16, gh), components.BlockWall),
	//	factory.CreateWall(ms.ecs, resolv.NewObject(gw-16, 0, 16, gh)),
	//	factory.CreateWall(ms.ecs, resolv.NewObject(0, 0, gw, 16)),
	//	factory.CreateWall(ms.ecs, resolv.NewObject(0, gh-24, gw, 32)),

	//	factory.CreatePlayer(ms.ecs),
	//)
}

func loadAssets() {
	for _, fn := range []func() error{
		assets.Load,
	} {
		if err := fn(); err != nil {
			panic(err)
		}
	}
}
