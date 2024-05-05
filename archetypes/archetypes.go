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
		layers.Player,
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
		layers.Actors,
		tags.Enemy,
		components.Enemy,
		components.AI,
		components.Animation,
		components.Object,
		components.CollistionPlayer,
		components.Velocity,
		components.AttackVector,
		components.Shooter,
		components.Health,
	)

	Bullet = NewArchetype(
		layers.Interactables,
		tags.Bullet,
		components.Bullet,
		components.Object,
		components.CollistionBullet,
		components.Velocity,
		components.Despawnable,
		components.Animation,
	)

	BouncerBullet = NewArchetype(
		layers.Interactables,
		tags.Bullet,
		components.Bullet,
		components.Object,
		components.Velocity,
		components.CollisionBouncer,
		components.Despawnable,
		components.Animation,
	)

	Wall = NewArchetype(
		layers.Architecture,
		tags.Wall,
		components.Block,
		components.Object,
		components.Animation,
	)

	Space = NewArchetype(
		layers.System,
		components.Space,
	)

	Camera = NewArchetype(
		layers.System,
		components.Camera,
	)

	WeaponSprite = NewArchetype(
		layers.Interactables,
		tags.WeaponSprite,
		components.Shooter,
		components.AttackVector,
		components.Animation,
		components.Object,
	)

	ParticleSprite = NewArchetype(
		layers.FX,
		tags.Particle,
		components.Animation,
		components.Object,
		components.Despawnable,
	)
)

type Archetype struct {
	Components []donburi.IComponentType
	Layer      ecs.LayerID
}

func NewArchetype(layer ecs.LayerID, cs ...donburi.IComponentType) *Archetype {
	return &Archetype{
		Components: cs,
		Layer:      layer,
	}
}

func (a *Archetype) Spawn(ecs *ecs.ECS, cs ...donburi.IComponentType) *donburi.Entry {
	e := ecs.World.Entry(ecs.Create(
		a.Layer,
		append(a.Components, cs...)...,
	))

	return e
}
