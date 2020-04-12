package systems

import (
	gc "github.com/x-hgg-x/space-invaders-go/lib/components"

	ecs "github.com/x-hgg-x/goecs/v2"
	ec "github.com/x-hgg-x/goecsengine/components"
	w "github.com/x-hgg-x/goecsengine/world"
)

// DeleteSystem removes deleted entities after animation end
func DeleteSystem(world w.World) {
	gameComponents := world.Components.Game.(*gc.Components)

	world.Manager.Join(gameComponents.Deleted, world.Components.Engine.AnimationControl).Visit(ecs.Visit(func(entity ecs.Entity) {
		animationControl := world.Components.Engine.AnimationControl.Get(entity).(*ec.AnimationControl)
		if animationControl.GetState().Type == ec.ControlStateDone {
			world.Manager.DeleteEntity(entity)
		}
	}))
}
