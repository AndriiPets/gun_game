package archetypes

import (
	"github.com/AndriiPets/FishGame/components"
	"github.com/AndriiPets/FishGame/layers"
	"github.com/AndriiPets/FishGame/tags"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

var (
	Player = NewArchetype(
		tags.Player,
		components.Player,
		components.Animation,
		components.Object,
		components.CollistionPlayer,
		components.Velocity,
		components.AttackVector,
		components.Shooter,
		components.Health,
	)

	Enemy = NewArchetype(
		tags.Enemy,
		components.Enemy,
		components.Animation,
		components.Object,
		components.CollistionPlayer,
		components.Velocity,
		components.AttackVector,
		components.Shooter,
		components.Health,
	)

	Bullet = NewArchetype(
		tags.Bullet,
		components.Bullet,
		components.Object,
		components.CollistionBullet,
		components.Velocity,
		components.Despawnable,
		components.Animation,
	)

	BouncerBullet = NewArchetype(
		tags.Bullet,
		components.Bullet,
		components.Object,
		components.Velocity,
		components.CollisionBouncer,
		components.Despawnable,
		components.Animation,
	)

	Wall = NewArchetype(
		tags.Wall,
		components.Block,
		components.Object,
		components.Animation,
	)

	Space = NewArchetype(
		components.Space,
	)

	Camera = NewArchetype(
		components.Camera,
	)

	WeaponSprite = NewArchetype(
		tags.WeaponSprite,
		components.Shooter,
		components.AttackVector,
		components.Animation,
		components.Object,
	)
)

type Archetype struct {
	Components []donburi.IComponentType
}

func NewArchetype(cs ...donburi.IComponentType) *Archetype {
	return &Archetype{
		Components: cs,
	}
}

func (a *Archetype) Spawn(ecs *ecs.ECS, cs ...donburi.IComponentType) *donburi.Entry {
	e := ecs.World.Entry(ecs.Create(
		layers.Default,
		append(a.Components, cs...)...,
	))

	return e
}
