package entities

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
		toMiningBeams:    ent.NewBus(),
	}
}

type Player struct {
	ent.CoreEntity
	ent.WithDraw
	ent.WithUpdate
	ent.WithActivePhysics
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
	minerals         int
	miningTimer      float64
	miningAsteroid   ent.EntityUUID
	lastfx           ent.BodyEffects
	toMiningBeams    *ent.Bus
}

func (p *Player) AfterAdd(w *ent.World) {
	w.AddTags(p, "player", "player_camera_target")
}

// Shape implements ent.ActivePhysicsBody.
func (p *Player) Shape() ent.Shape {
	return ent.Circle{
		Center: p.Position(),
		Radius: p.radius,
	}
}

func (p *Player) Radius() float64 {
	return p.radius
}

type MiningBeamOff struct{}

func (p *Player) Update(win *pixelgl.Window, entities *ent.World, dt float64) {
	// Deal with dead player
	if p.dead {
		entities.Instantiate(
			NewExplosion(p.Position(), 1),
			NewPlayer(),
		)
		entities.Destroy(p)
		return
	}
	// Handle state changes of mining
	if win.JustPressed(pixelgl.KeySpace) {
		asteroid, ok := p.selectClosestAsteroid(entities)
		if ok {
			beam := NewMiningBeam(p.UUID(), asteroid.UUID())
			ent.Subscribe(p.toMiningBeams, beam)
			p.miningTimer = 0
			entities.Instantiate(beam)
			p.miningAsteroid = asteroid.UUID()
		}
	} else if win.JustReleased(pixelgl.KeySpace) {
		p.miningAsteroid = ""
		ent.Emit(entities, p.toMiningBeams, MiningBeamOff{})
	}

	// Handle mining
	if p.miningAsteroid != "" {
		asteroid, ok := p.getTargetAsteroid(entities)
		if !ok || asteroid.Position().To(p.Position()).Len() > 10 {
			p.miningAsteroid = ""
			ent.Emit(entities, p.toMiningBeams, MiningBeamOff{})
		} else {
			p.miningTimer += dt
			if p.miningTimer > 1 {
				p.miningTimer = 0
				p.minerals++
				asteroid.resources--
				if asteroid.resources <= 0 {
					entities.Destroy(asteroid)
					entities.Instantiate(NewExplosion(asteroid.Position(), asteroid.Radius()))
				} else {
					edgePos := asteroid.Position().To(p.Position()).Unit().Scaled(asteroid.radius).Add(asteroid.Position())
					entities.Instantiate(NewExplosion(edgePos, 0.3))
				}
			}
		}
	}

	// Handle ship movement
	fx := ent.BodyEffects{}
	if win.Pressed(pixelgl.KeyW) {
		fx.Force = fx.Force.Add(ent.Forward(p).Scaled(p.boosterForce))
	}
	if win.Pressed(pixelgl.KeyA) {
		fx.Torque += p.boosterTorue
	}
	if win.Pressed(pixelgl.KeyD) {
		fx.Torque -= p.boosterTorue
	}
	fx.Force = fx.Force.Add(ent.CalculateDragForce(p.Velocity(), p.linearDragCoeff, 0.5))
	fx.Torque += ent.CalculateDragTorque(p.AngularVelocity(), p.angularDragCoeff, 0.8)
	p.lastfx = fx

	// Handle timers
	p.lastDamageTimer += dt
	p.bubbleTimer -= dt
}

func (p *Player) PysicsUpdate(dt float64) {
	ent.EulerStateUpdate(p, p.lastfx, dt)
}

func (p *Player) getTargetAsteroid(entities *ent.World) (*Asteroid, bool) {
	return ent.OneOfType[*Asteroid](entities.WithUUID(p.miningAsteroid))
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
	drawMat := pixel.IM.Scaled(
		pixel.ZV,
		p.radius*2.0/p.sprite.Frame().W(),
	).Rotated(
		pixel.ZV,
		-math.Pi/2,
	).Chained(
		ent.TransMat(p),
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

func (p *Player) OnCollision(col ent.Collision) {
	if p.sheilds <= 0 {
		p.dead = true
		return
	}
	if p.lastDamageTimer > 0.5 {
		p.lastDamageTimer = 0
		p.SetVelocity(p.Velocity().Add(col.Normal.Scaled(10)))
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
	ent.CoreEntity
	ent.WithDraw
	ent.WithUpdate
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
func (e *Explosion) Update(win *pixelgl.Window, all *ent.World, dt float64) {
	e.timer += dt
	if e.timer >= 0.5 {
		all.Destroy(e)
	}
}

func (e *Explosion) Position() pixel.Vec {
	return e.pos
}
