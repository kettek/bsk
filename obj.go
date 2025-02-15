package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joonazan/vec2"
	"github.com/kettek/bsk/assets"
)

type Object interface {
	Update(*StatePlay) Request
	Draw(*ebiten.Image)
	Velocity() *vec2.Vector
	Position() *vec2.Vector
	Ground()
	Size() vec2.Vector
	TouchCell(CellCollision) Request
	FallOut() Request
}

type BaseObject struct {
	I assets.ImageID
	P vec2.Vector
	V vec2.Vector
}

func (b *BaseObject) Velocity() *vec2.Vector {
	return &b.V
}

func (b *BaseObject) Position() *vec2.Vector {
	return &b.P
}
