package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	State State
}

func (g *Game) Update() error {
	if state := g.State.Update(); state != nil {
		g.State = state
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.State.Draw(screen)
}

func (g *Game) Layout(ow, oh int) (int, int) {
	return 160, 144
}

type State interface {
	Update() State
	Draw(screen *ebiten.Image)
}
