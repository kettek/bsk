package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/bsk/assets"
)

type StatePlay struct {
	levels map[string]*Level
	level  *Level
}

func NewStatePlay() *StatePlay {
	return &StatePlay{
		levels: make(map[string]*Level),
	}
}

func (s *StatePlay) Update() State {
	if s.level == nil {
		s.level = s.LoadLevel("start")
		entrances := s.level.FindCells(CellFlagEnter)
		birb := NewBirb()
		if len(entrances) > 0 {
			birb.Position().X = float64(entrances[0].cellX * 8)
			birb.Position().Y = float64(entrances[0].cellY * 8)
		}
		s.level.AddObject(birb)
	}

	for _, req := range s.level.Update() {
		switch req := req.(type) {
		case RequestLevel:
			if req.Level == "outro" {
				return NewStateOutro(s.level)
			}
			nextLevel := s.LoadLevel(req.Level)
			entrances := nextLevel.FindCells(CellFlagEnter)
			if len(entrances) == 0 {
				panic("no entrances")
			}
			s.level.RemoveObject(req.object)
			s.level = s.LoadLevel(req.Level)
			s.level.AddObject(req.object)
			req.object.Position().X = float64(entrances[0].cellX * 8)
			req.object.Position().Y = float64(entrances[0].cellY * 8)
		case RequestDelete:
			fmt.Println("deleting")
			s.level.RemoveObject(req.object)
		case RequestAdd:
			fmt.Println("adding")
			s.level.AddObject(req.object)
		}
	}

	return nil
}

func (s *StatePlay) Draw(screen *ebiten.Image) {
	if s.level == nil {
		return
	}
	s.level.Draw(screen)
	for _, o := range s.level.Objects {
		o.Draw(screen)
	}
	s.level.DrawForeground(screen)

	s.level.DrawDebug(screen)
}

func (s *StatePlay) LoadLevel(n string) *Level {
	if s.levels[n] == nil {
		data, err := assets.GetLevelSource(n)
		if err != nil {
			panic(err)
		}
		s.levels[n] = &Level{}
		if err := s.levels[n].UnmarshalBinary(data); err != nil {
			panic(err)
		}
	}
	return s.levels[n]
}
