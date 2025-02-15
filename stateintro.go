package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/kettek/bsk/assets"
)

type StateIntro struct {
	level         *Level
	messages      []string
	messageLife   int
	birdX         int
	captured      bool
	offscreenBird bool
	delay         int
}

func NewStateIntro() *StateIntro {
	data, err := assets.GetLevelSource("intro")
	if err != nil {
		panic(err)
	}

	level := &Level{}

	if err := level.UnmarshalBinary(data); err != nil {
		panic(err)
	}

	return &StateIntro{
		level: level,
		messages: []string{
			"*kis*",
			"",
			"*smoch*",
		},
	}
}

func (s *StateIntro) Update() State {
	if ebiten.IsKeyPressed(ebiten.KeyEnter) || ebiten.IsKeyPressed(ebiten.KeySpace) {
		return NewStatePlay()
	}
	if len(s.messages) > 0 {
		if s.messageLife > 90 {
			s.messages = s.messages[1:]
			s.messageLife = 0
			if len(s.messages) == 0 {
				s.birdX = -16
			}
		}
		s.messageLife++
	} else {
		if s.birdX < 160 {
			s.birdX++
			if s.birdX >= 4*8 && !s.captured {
				s.captured = true
				s.delay = 20
				s.level.Cells[10][7] = 0
			}
		} else {
			if !s.offscreenBird {
				s.offscreenBird = true
				s.delay = 20
			}
		}
	}
	s.delay--
	return nil
}

func (s *StateIntro) Draw(screen *ebiten.Image) {
	s.level.Draw(screen)
	if len(s.messages) > 0 {
		ebitenutil.DebugPrintAt(screen, s.messages[0], 40, 60)
	} else {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(s.birdX), 55)
		screen.DrawImage(assets.GetImage(assets.GetImageID("bigbird")), op)
		if s.captured {
			op.GeoM.Translate(3*8, float64(assets.GetImage(assets.GetImageID("bigbird")).Bounds().Dy()-8))
			screen.DrawImage(assets.GetImage(assets.GetImageID("kit")), op)
			if s.delay <= 0 {
				if s.offscreenBird {
					ebitenutil.DebugPrintAt(screen, "i gotta save kitty!!!", 20, 60)
				} else {
					x := s.birdX - 10
					ebitenutil.DebugPrintAt(screen, "hep me!!!", x, 40)
				}
			}
		}
	}
}
