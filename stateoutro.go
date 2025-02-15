package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/kettek/bsk/assets"
)

type StateOutro struct {
	level       *Level
	messages    []string
	messageLife int
}

func NewStateOutro(level *Level) *StateOutro {
	return &StateOutro{
		level: level,
		messages: []string{
			"*smoch!!!*",
			"*kis*!!!",
			"let's fly outta here!",
			"",
		},
	}
}

func (s *StateOutro) Update() State {
	if len(s.messages) > 0 {
		if s.messageLife > 130 {
			s.messages = s.messages[1:]
			s.messageLife = 0
		}
		s.messageLife++
	}
	return nil
}

func (s *StateOutro) Draw(screen *ebiten.Image) {
	if len(s.messages) == 0 {
		ebitenutil.DebugPrintAt(screen, "DA END", 60, 45)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(62, 60)
		screen.DrawImage(assets.GetImage(assets.GetImageID("love")), op)
		return
	}
	s.level.Draw(screen)
	for _, obj := range s.level.Objects {
		obj.Draw(screen)
	}
	s.level.DrawForeground(screen)

	if len(s.messages) > 0 {
		if len(s.messages) == 2 {
			ebitenutil.DebugPrintAt(screen, s.messages[0], 30, 40)
		} else {
			ebitenutil.DebugPrintAt(screen, s.messages[0], 50, 40)
		}
	}
}
