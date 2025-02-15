package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joonazan/vec2"
	"github.com/kettek/bsk/assets"
)

type Birb struct {
	BaseObject
	FaceRight bool
	Jumping   bool
	JumpJuice float64
	Firing    bool
	FireJuice float64
}

func NewBirb() *Birb {
	return &Birb{
		BaseObject: BaseObject{
			I: assets.GetImageID("birb"),
			P: vec2.Vector{X: 100, Y: 100},
		},
	}
}

func (b *Birb) Update(s *StatePlay) Request {
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		if b.Jumping {
			b.V.X -= 0.05
		} else {
			b.V.X -= 0.1
		}

		b.FaceRight = false
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		if b.Jumping {
			b.V.X += 0.05
		} else {
			b.V.X += 0.1
		}
		b.FaceRight = true
	}
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		if b.JumpJuice > 0 {
			b.JumpJuice--
			b.V.Y -= 0.4
			if b.V.Y < -2 {
				b.V.Y = -2
			}
		}
		b.Jumping = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		if b.FireJuice > 0 {
			b.FireJuice = -20
			if b.FaceRight {
				return RequestAdd{object: NewBorb(b.P, vec2.Vector{X: 1, Y: 0}, true)}
			} else {
				return RequestAdd{object: NewBorb(b.P, vec2.Vector{X: -1, Y: 0}, false)}
			}
		}
	}
	b.FireJuice++
	return nil
}

func (b *Birb) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	if b.FaceRight {
		op.GeoM.Translate(-8, 0)
		op.GeoM.Scale(-1, 1)
	}

	// Draw from center.
	op.GeoM.Translate(-1, -4)

	op.GeoM.Translate(b.P.X, b.P.Y)
	screen.DrawImage(assets.GetImage(b.I), op)
}

func (b *Birb) TouchCell(coll CellCollision) Request {
	if coll.cell.Flags()&CellFlagExit != 0 {
		return RequestLevel{}
	} else if coll.cell.Flags()&CellFlagClimbable != 0 {
		b.JumpJuice = 1
	} else if coll.cell.Flags()&CellFlagHurts != 0 {
		return RequestReset{}
	} else if coll.cell.Flags()&CellFlagPushUp != 0 {
		b.V.Y = -2
		b.Ground()
	}
	return nil
}

func (b *Birb) FallOut() Request {
	return RequestReset{}
}

func (b *Birb) Velocity() *vec2.Vector {
	return &b.V
}

func (b *Birb) Position() *vec2.Vector {
	return &b.P
}

func (b *Birb) Ground() {
	b.JumpJuice = 5
	b.Jumping = false
}

func (b *Birb) Size() vec2.Vector {
	return vec2.Vector{X: 4, Y: 4}
}
