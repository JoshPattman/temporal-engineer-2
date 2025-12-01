package entities

import (
	"ent"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

func NewMiningBeam(startUUID, endUUID ent.EntityUUID) *MiningBeam {
	return &MiningBeam{
		sprite:  GlobalSpriteManager.FullSprite("tether.png"),
		startID: startUUID,
		endID:   endUUID,
	}
}

type MiningBeam struct {
	ent.CoreEntity
	ent.WithUpdate
	ent.WithDraw
	sprite   *pixel.Sprite
	startID  ent.EntityUUID
	endID    ent.EntityUUID
	startPos pixel.Vec
	endPos   pixel.Vec
	inverted bool
	timer    float64
	destroy  bool
}

func (e *MiningBeam) Update(win *pixelgl.Window, world *ent.World, dt float64) {
	start, okStart := ent.OneOfType[ent.Transform](world.WithUUID(e.startID))
	if okStart {
		e.startPos = start.Position()
	}
	end, okEnd := ent.OneOfType[ent.Transform](world.WithUUID(e.endID))
	if okEnd {
		e.endPos = end.Position()
	}
	e.timer += dt
	if e.timer > 0.2 {
		e.timer = 0
		e.inverted = !e.inverted
	}
	if e.destroy {
		world.Destroy(e)
	}
}

func (e *MiningBeam) Draw(win *pixelgl.Window, _ *ent.World, worldToScreen pixel.Matrix) {
	dist := e.startPos.To(e.endPos).Len()
	if dist == 0 {
		return
	}
	yScale := 1.0
	if e.inverted {
		yScale = -1.0
	}
	e.sprite.Draw(
		win,
		pixel.IM.
			Scaled(pixel.ZV, 1.0/16.0).
			Moved(pixel.V(0.5, 0)).
			ScaledXY(pixel.ZV, pixel.V(dist, yScale)).
			Rotated(pixel.ZV, e.startPos.To(e.endPos).Angle()).
			Moved(e.startPos).
			Chained(worldToScreen),
	)
}

func (e *MiningBeam) Destroy(struct{}) {
	e.destroy = true
}

type MiningBeamOff struct{}

func (e *MiningBeam) HandleMessage(world *ent.World, msg any) {
	switch msg.(type) {
	case MiningBeamOff:
		world.Destroy(e)
	}
}
