package components

import (
	ecs "github.com/x-hgg-x/goecs"
)

// Components contains references to all game components
type Components struct {
	Player *ecs.Component
}

// Player component
type Player struct {
	Width float64
}
