package components

import (
	ecs "github.com/x-hgg-x/goecs"
)

// Components contains references to all game components
type Components struct {
	Player       *ecs.Component
	PlayerBullet *ecs.Component
	Deleted      *ecs.Component
}

// Player component
type Player struct {
	Width float64
}

// PlayerBullet component
type PlayerBullet struct {
	Height   float64
	Velocity float64
}

// Deleted component
type Deleted struct{}
