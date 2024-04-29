package factory

import (
	"github.com/AndriiPets/FishGame/archetypes"
	"github.com/AndriiPets/FishGame/components"
	dresolv "github.com/AndriiPets/FishGame/resolv"
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func CreateEnemy(ecs *ecs.ECS, posX, posY float64, enemyType components.EnemyType) *donburi.Entry {
	enemyEntry := archetypes.Enemy.Spawn(ecs)

	//setup initial state
	enemy := components.Enemy.Get(enemyEntry)
	enemy.State = components.EnemyStateIdle
	enemy.Type = enemyType

	//setup enemy stats
	health := components.Health.Get(enemyEntry)
	health.Ammount = 3
	health.Cooldown = 0.4

	//setup animation
	animation := components.Animation.Get(enemyEntry)
	animation.Animation = enemy.Animation()
	animation.FlipH = false
	animation.Type = components.AnimationActor

	//setup enemy object
	obj := resolv.NewObject(posX, posY, 16, 16)
	obj.AddTags("damageable")
	dresolv.SetObject(enemyEntry, obj)

	return enemyEntry

}
