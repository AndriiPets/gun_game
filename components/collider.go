package components

import "github.com/yohamta/donburi"

type CollisionType struct {
	Type string
}

var (
	CollistionPlayer = donburi.NewComponentType[CollisionType]()
	CollistionBullet = donburi.NewComponentType[CollisionType]()
	CollisionBouncer = donburi.NewComponentType[CollisionType]()
)
