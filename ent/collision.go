package ent

import "github.com/gopxl/pixel"

// All the information about a collision that was detected by the physics engine.
type Collision struct {
	Self   PhysicsBody
	Other  PhysicsBody
	Normal pixel.Vec
	Point  pixel.Vec
}

// An interface that is able to listen to collisions.
type CollisionListener interface {
	// Additional logic to be run after a collision is detected and resolved.
	OnCollision(Collision)
}

// Flip the observer and self of the collision, and also the normal.
func (c Collision) ForOther() Collision {
	return Collision{
		Self:   c.Other,
		Other:  c.Self,
		Normal: c.Normal.Scaled(-1),
		Point:  c.Point,
	}
}

func calculateFinalVelocity(m1, m2, e, u1, u2 float64) float64 {
	return (((m1 - e*m2) * u1) + (((1 + e) * m2) * u2)) / (m1 + m2)
}

func calculateFinalVelocityLimInf(e, u1, u2 float64) float64 {
	return (-e * u1) + ((1 + e) * u2)
}

func decomposeAxis(u pixel.Vec, ax pixel.Vec) (float64, pixel.Vec) {
	dist := u.Dot(ax)
	leftover := u.Sub(ax.Scaled(dist))
	return dist, leftover
}

func recomposeAxis(leftover pixel.Vec, ax pixel.Vec, dist float64) pixel.Vec {
	return leftover.Add(ax.Scaled(dist))
}

// TODO: Both the below functions may exhibit buggy behaviour when colliding with multiple objects

// Checks for a collision and updates the two active bodies
func checkActiveBodies(a, b ActivePhysicsBody) (Collision, bool) {
	col := collideShapes(a.Shape(), b.Shape())
	if !col.collided {
		return Collision{}, false
	}
	aState := a.State()
	bState := b.State()

	combinedElasticity := 0.5 * (a.Elasticity() + b.Elasticity())
	aSpeed, aLeftover := decomposeAxis(aState.Velocity, col.normal)
	bSpeed, bLeftover := decomposeAxis(bState.Velocity, col.normal)
	newASpeed := calculateFinalVelocity(a.Mass(), b.Mass(), combinedElasticity, aSpeed, bSpeed)
	newBSpeed := calculateFinalVelocity(b.Mass(), a.Mass(), combinedElasticity, bSpeed, aSpeed)
	newAVel := recomposeAxis(aLeftover, col.normal, newASpeed)
	newBVel := recomposeAxis(bLeftover, col.normal, newBSpeed)

	correction := col.normal.Scaled(col.overlap / 2)

	a.SetState(BodyState{
		Position:        aState.Position.Sub(correction),
		Velocity:        newAVel,
		Angle:           aState.Angle,
		AngularVelocity: aState.AngularVelocity,
	})

	b.SetState(BodyState{
		Position:        bState.Position.Add(correction),
		Velocity:        newBVel,
		Angle:           bState.Angle,
		AngularVelocity: bState.AngularVelocity,
	})

	return Collision{
		Self:   a,
		Other:  b,
		Normal: col.normal.Scaled(-1),
	}, true
}

// Checks for a collision and updates the active body
func checkActiveAndKinematicBodies(a ActivePhysicsBody, b PhysicsBody) (Collision, bool) {
	col := collideShapes(a.Shape(), b.Shape())
	if !col.collided {
		return Collision{}, false
	}
	aState := a.State()
	bState := b.State()

	combinedElasticity := 0.5 * (a.Elasticity() + b.Elasticity())
	aSpeed, aLeftover := decomposeAxis(aState.Velocity, col.normal)
	bSpeed, _ := decomposeAxis(bState.Velocity, col.normal)
	newASpeed := calculateFinalVelocityLimInf(combinedElasticity, aSpeed, bSpeed)
	newAVel := recomposeAxis(aLeftover, col.normal, newASpeed)

	correction := col.normal.Scaled(col.overlap)

	a.SetState(BodyState{
		Position:        aState.Position.Sub(correction),
		Velocity:        newAVel,
		Angle:           aState.Angle,
		AngularVelocity: aState.AngularVelocity,
	})

	return Collision{
		Self:   a,
		Other:  b,
		Normal: col.normal.Scaled(-1),
	}, true
}

// Perform a collision physics update on the set of bodies.
// Perform corrections to overlapping objects, and returns collisions to be passed to handlers.
func StatelessCollisionPhysics(bodies []PhysicsBody) []Collision {
	// Sort bodies
	kinematicBodies := make([]PhysicsBody, 0)
	activeBodies := make([]ActivePhysicsBody, 0)
	for _, b := range bodies {
		ab, ok := b.(ActivePhysicsBody)
		if !ok || !ab.IsPhysicsActive() {
			kinematicBodies = append(kinematicBodies, b)

		} else {
			activeBodies = append(activeBodies, ab)
		}
	}

	collisions := make([]Collision, 0)

	// Collide static and active
	for _, a := range activeBodies {
		for _, b := range kinematicBodies {
			col, ok := checkActiveAndKinematicBodies(a, b)
			if ok {
				collisions = append(collisions, col)
			}
		}
	}

	// Collide active and active
	for i, a := range activeBodies {
		for j, b := range activeBodies {
			if j <= i {
				continue
			}
			col, ok := checkActiveBodies(a, b)
			if ok {
				collisions = append(collisions, col)
			}
		}
	}
	return collisions
}
