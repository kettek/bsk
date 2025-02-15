package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joonazan/vec2"
	"github.com/kettek/bsk/assets"
)

type Borb struct {
	BaseObject
	FaceRight     bool
	Up            bool
	Lifetime      int
	DestroyedCell bool
}

func NewBorb(p vec2.Vector, v vec2.Vector, dir bool) *Borb {
	return &Borb{
		BaseObject: BaseObject{
			I: assets.GetImageID("borb"),
			P: p,
			V: v,
		},
		FaceRight: dir,
		Up:        false,
		Lifetime:  200,
	}
}

func (b *Borb) Update(s *StatePlay) Request {
	if b.Up {
		b.V.Y = -0.5
	} else {
		b.V.Y = 0.5
	}
	if b.FaceRight {
		b.V.X = 1
	} else {
		b.V.X = -1
	}
	b.Lifetime--
	if b.Lifetime <= 0 || b.DestroyedCell {
		return RequestDelete{object: b}
	}
	return nil
}

func (b *Borb) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	if b.FaceRight {
		op.GeoM.Translate(-8, 0)
		op.GeoM.Scale(-1, 1)
	}

	op.GeoM.Translate(b.P.X, b.P.Y)
	screen.DrawImage(assets.GetImage(b.I), op)
}

func (b *Borb) TouchCell(coll CellCollision) Request {
	if coll.cell.Flags()&CellFlagSolid != 0 {
		if coll.top < b.P.Y {
			b.Up = false
		} else if coll.bottom > b.P.Y {
			b.Up = true
		}
	}
	if coll.cell.Flags()&CellFlagDestroyable != 0 {
		b.DestroyedCell = true
		return RequestSetCell{flag: 0, image: 0, x: coll.cellX, y: coll.cellY}
	}
	return nil
}

func (b *Borb) FallOut() Request {
	return RequestDelete{object: b}
}

func (b *Borb) Ground() {
}

func (b *Borb) Size() vec2.Vector {
	return vec2.Vector{X: 4, Y: 4}
}
