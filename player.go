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
	miningPos        pixel.Vec
	mining           bool
	lastInvertTimer  float64
	invert           bool
	minerals         int
	miningTimer      float64
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
	// return ent.Circle{
	// 	Center: p.pos,
	// 	Radius: p.radius,
	// }
	return ent.MultiShape{
		Shapes: []ent.Shape{
			ent.Circle{
				Center: p.pos,
				Radius: p.radius,
			},
			ent.Line{
				A: p.pos,
				B: p.pos.Add(p.Forward().Scaled(6)),
			},
		},
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
	p.lastInvertTimer += dt
	if p.lastInvertTimer > 0.2 {
		p.lastInvertTimer = 0
		p.invert = !p.invert
	}
	// Deal with dead player
	if p.dead {
		return []ent.Entity{
				NewExplosion(p.pos, 1),
				NewPlayer(),
			}, []ent.Entity{
				p,
			}
	}
	// Handle state changes of mining
	if win.JustPressed(pixelgl.KeySpace) {
		_, ok := p.getTargetAsteroid(entities)
		if !ok {
			asteroid, ok := p.selectClosestAsteroid(entities)
			if ok {
				entities.AddTags(asteroid, "player_target")
			}
		}
		p.mining = true
		p.miningTimer = 0
	} else if win.JustReleased(pixelgl.KeySpace) {
		asteroid, ok := p.getTargetAsteroid(entities)
		if ok {
			entities.RemoveTags(asteroid, "player_target")
		}
		p.mining = false
	}

	// Handle mining
	var destroyAsteroids []ent.Entity
	var addExplosions []ent.Entity
	if p.mining {
		asteroid, ok := p.getTargetAsteroid(entities)
		if !ok {
			p.mining = false
		} else if asteroid.Position().To(p.Position()).Len() > 10 {
			p.mining = false
			entities.RemoveTags(asteroid, "player_target")
		} else {
			p.miningPos = asteroid.Position()
		}
		p.miningTimer += dt
		if p.miningTimer > 1 {
			p.miningTimer = 0
			p.minerals++
			asteroid.resources--
			if asteroid.resources <= 0 {
				destroyAsteroids = append(destroyAsteroids, asteroid)
				addExplosions = append(addExplosions, NewExplosion(asteroid.pos, asteroid.Radius()))
			} else {
				edgePos := asteroid.pos.To(p.pos).Unit().Scaled(asteroid.radius).Add(asteroid.pos)
				addExplosions = append(addExplosions, NewExplosion(edgePos, 0.3))
			}
		}
	}

	// Handle ship movement
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

	// Handle timers
	p.lastDamageTimer += dt
	p.bubbleTimer -= dt
	return addExplosions, destroyAsteroids
}

func (p *Player) getTargetAsteroid(entities *ent.World) (*Asteroid, bool) {
	return ent.First(
		ent.OfType[*Asteroid](
			entities.ForTag("player_target"),
		),
	)
}

func (p *Player) selectClosestAsteroid(entities *ent.World) (*Asteroid, bool) {
	return ent.Closest(
		p.Position(),
		ent.OfType[*Asteroid](
			entities.ForTag("mineable_asteroid"),
		),
	)
}

func (p *Player) Draw(win *pixelgl.Window, _ *ent.World, worldToScreen pixel.Matrix) {
	if p.mining {
		t := NewTether()
		t.end = p.pos
		t.start = p.miningPos
		t.inverted = p.invert
		t.Draw(win, worldToScreen)
	}
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

func (p *Player) AfterAdd(w *ent.World) {
	w.AddTags(p, "player", "player_camera_target")
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

func (p *Player) Shields() int {
	return p.sheilds
}

func (p *Player) Minerals() int {
	return p.minerals
}

func NewExplosion(at pixel.Vec, scale float64) *Explosion {
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
		scale:   scale,
	}
}

var _ ent.Entity = &Explosion{}

type Explosion struct {
	ent.EntityBase
	pos     pixel.Vec
	timer   float64
	sprites []*pixel.Sprite
	scale   float64
}

// Draw implements ent.Entity.
func (e *Explosion) Draw(win *pixelgl.Window, _ *ent.World, worldToScreen pixel.Matrix) {
	idx := int(e.timer / 0.5 * float64(len(e.sprites)))
	s := e.sprites[idx]
	s.Draw(
		win,
		pixel.IM.Scaled(pixel.ZV, 0.1*e.scale).Moved(e.pos).Chained(worldToScreen),
	)
}

func (e *Explosion) DrawLayer() int { return -1 }

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

func NewTether() *Tether {
	return &Tether{
		sprite: GlobalSpriteManager.FullSprite("tether.png"),
	}
}

type Tether struct {
	start    pixel.Vec
	end      pixel.Vec
	sprite   *pixel.Sprite
	inverted bool
}

func (e *Tether) Draw(win *pixelgl.Window, worldToScreen pixel.Matrix) {
	dist := e.start.To(e.end).Len()
	if dist == 0 {
		return
	}
	yScale := 1.0
	if e.inverted {
		yScale = -1.0
	}
	e.sprite.Draw(
		win,
		pixel.IM.Scaled(pixel.ZV, 1.0/16.0).Moved(pixel.V(0.5, 0)).ScaledXY(pixel.ZV, pixel.V(dist, yScale)).Rotated(pixel.ZV, e.start.To(e.end).Angle()).Moved(e.start).Chained(worldToScreen),
	)
}
