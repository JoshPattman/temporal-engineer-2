package main

import (
	"ent"
	"math"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

func NewPlayer() *Player {
	shipSprite := GlobalSpriteManager.FullSprite("ship.png")
	bubbleSprite := GlobalSpriteManager.FullSprite("bubble.png")
	return &Player{
		sprite:           shipSprite,
		radius:           1,
		boosterForce:     50,
		linearDragCoeff:  0.3,
		angularDragCoeff: 0.6,
		boosterTorue:     12,
		bubbleSprite:     bubbleSprite,
		sheilds:          3,
	}
}

var _ ent.ActivePhysicsBody = &Player{}

type Player struct {
	ent.EntityBase
	Transform
	velocity         pixel.Vec
	angularSpeed     float64
	sprite           *pixel.Sprite
	bubbleSprite     *pixel.Sprite
	radius           float64
	boosterForce     float64
	linearDragCoeff  float64
	angularDragCoeff float64
	boosterTorue     float64
	lastDamageTimer  float64
	bubbleTimer      float64
	sheilds          int
	dead             bool
}

// Elasticity implements ent.ActivePhysicsBody.
func (p *Player) Elasticity() float64 {
	return 0.3
}

// IsPhysicsActive implements ent.ActivePhysicsBody.
func (p *Player) IsPhysicsActive() bool {
	return true
}

// Mass implements ent.ActivePhysicsBody.
func (p *Player) Mass() float64 {
	return 1
}

// SetState implements ent.ActivePhysicsBody.
func (p *Player) SetState(state ent.BodyState) {
	p.pos = state.Position
	p.velocity = state.Velocity
	p.angularSpeed = state.AngularVelocity
	p.rot = state.Angle
}

// Shape implements ent.ActivePhysicsBody.
func (p *Player) Shape() ent.Shape {
	return ent.Circle{
		Center: p.pos,
		Radius: p.radius,
	}
}

// State implements ent.ActivePhysicsBody.
func (p *Player) State() ent.BodyState {
	return ent.BodyState{
		Position:        p.pos,
		Velocity:        p.velocity,
		Angle:           p.rot,
		AngularVelocity: p.angularSpeed,
	}
}

func (p *Player) Radius() float64 {
	return p.radius
}

func (p *Player) Update(win *pixelgl.Window, entities *ent.World, dt float64) ([]ent.Entity, []ent.Entity) {
	if p.dead {
		return []ent.Entity{
				NewExplosion(p.pos),
				NewPlayer(),
			}, []ent.Entity{
				p,
			}
	}
	fx := ent.BodyEffects{}

	if win.Pressed(pixelgl.KeyW) {
		fx.Force = fx.Force.Add(p.Forward().Scaled(p.boosterForce))
	}
	if win.Pressed(pixelgl.KeyA) {
		fx.Torque += p.boosterTorue
	}
	if win.Pressed(pixelgl.KeyD) {
		fx.Torque -= p.boosterTorue
	}

	fx.Force = fx.Force.Add(ent.CalculateDragForce(p.velocity, p.linearDragCoeff, 0.5))
	fx.Torque += ent.CalculateDragTorque(p.angularSpeed, p.angularDragCoeff, 0.8)

	ent.EulerStateUpdate(p, fx, dt)

	p.lastDamageTimer += dt
	p.bubbleTimer -= dt
	return nil, nil
}

func (p *Player) Draw(win *pixelgl.Window, worldToScreen pixel.Matrix) {
	drawMat := pixel.IM.Scaled(
		pixel.ZV,
		p.radius*2.0/p.sprite.Frame().W(),
	).Rotated(
		pixel.ZV,
		-math.Pi/2,
	).Chained(
		p.Mat(),
	).Chained(
		worldToScreen,
	)
	p.sprite.Draw(win, drawMat)
	if p.bubbleTimer > 0 {
		p.bubbleSprite.DrawColorMask(
			win,
			pixel.IM.Scaled(pixel.ZV, 0.6).Chained(drawMat),
			pixel.Alpha(p.bubbleTimer/0.5),
		)
	}
}

func (p *Player) Tags() []string {
	return []string{
		"player",
		"player_camera_target",
	}
}

func (p *Player) OnCollision(col ent.Collision) {
	if p.sheilds <= 0 {
		p.dead = true
		return
	}
	if p.lastDamageTimer > 0.5 {
		p.lastDamageTimer = 0
		p.velocity = p.velocity.Add(col.Normal.Scaled(10))
		p.bubbleTimer = 0.5
		p.sheilds--
	}
}

func NewExplosion(at pixel.Vec) *Explosion {
	sprites := GlobalSpriteManager.TiledSprites(
		"boom.png",
		36,
		[]TilePos{
			{0, 1},
			{1, 1},
			{2, 1},
			{3, 1},
			{4, 1},
			{0, 0},
			{1, 0},
			{2, 0},
			{3, 0},
			{4, 0},
		},
	)
	return &Explosion{
		pos:     at,
		timer:   0,
		sprites: sprites,
	}
}

var _ ent.Entity = &Explosion{}

type Explosion struct {
	ent.EntityBase
	pos     pixel.Vec
	timer   float64
	sprites []*pixel.Sprite
}

// Draw implements ent.Entity.
func (e *Explosion) Draw(win *pixelgl.Window, worldToScreen pixel.Matrix) {
	idx := int(e.timer / 0.5 * float64(len(e.sprites)))
	s := e.sprites[idx]
	s.Draw(
		win,
		pixel.IM.Scaled(pixel.ZV, 0.1).Moved(e.pos).Chained(worldToScreen),
	)
}

// Tags implements ent.Entity.
func (e *Explosion) Tags() []string {
	return []string{"player_camera_target"}
}

// Update implements ent.Entity.
func (e *Explosion) Update(win *pixelgl.Window, all *ent.World, dt float64) (toCreate []ent.Entity, toDestroy []ent.Entity) {
	e.timer += dt
	if e.timer >= 0.5 {
		return nil, []ent.Entity{e}
	}
	return nil, nil
}

func (e *Explosion) Position() pixel.Vec {
	return e.pos
}
