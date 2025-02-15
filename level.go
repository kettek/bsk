package main

import (
	"image/color"
	"math"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/joonazan/vec2"
	"github.com/kettek/bsk/assets"
)

type CellFlag int16

const (
	CellFlagNone CellFlag = 1 << iota
	CellFlagSolid
	CellFlagClimbable
	CellFlagHurts
	CellFlagDestroyable
	CellFlagEnter
	CellFlagExit
	CellFlagForeground
	CellFlagPushUp
)

func CellFlagFromStrings(ss []string) CellFlag {
	var f CellFlag
	for _, s := range ss {
		switch s {
		case "solid":
			f |= CellFlagSolid
		case "climbable":
			f |= CellFlagClimbable
		case "hurts":
			f |= CellFlagHurts
		case "destroyable":
			f |= CellFlagDestroyable
		case "enter":
			f |= CellFlagEnter
		case "exit":
			f |= CellFlagExit
		case "foreground":
			f |= CellFlagForeground
		case "pushup":
			f |= CellFlagPushUp
		}
	}
	return f
}

// Bits 0-15 represent the image ID
// Bits 16-31 represent the Flags.
type Cell int32

func (c Cell) ImageID() assets.ImageID {
	// Return bits 0 through 15.
	return assets.ImageID(c & 0xFFFF)
}

func (c Cell) Flags() CellFlag {
	// Return bits 16 through 31.
	return CellFlag(c >> 16)
}

type Level struct {
	Cells     [18][20]Cell
	Objects   []Object
	NextLevel string
	nextLevel *Level
	enter     [2]int
}

func (l *Level) UnmarshalBinary(data []byte) error {
	lines := strings.Split(string(data), "\n")

	cellDefs := make(map[string]Cell)

	mode := 0
	y := 0
	for _, line := range lines {
		if line == "" {
			mode++
			continue
		}
		if mode == 0 {
			l.NextLevel = line
		} else if mode == 1 {
			parts := strings.Split(line, " ")
			id := assets.GetImageID(parts[1])
			flag := CellFlagNone
			if len(parts) > 2 {
				flagParts := strings.Split(parts[2], ",")
				flag = CellFlagFromStrings(flagParts)
			}
			cellDefs[parts[0]] = Cell(int32(id) | int32(flag)<<16)
		} else if mode == 2 {
			for x, c := range line {
				if cellDefs[string(c)].Flags()&CellFlagEnter != 0 {
					l.enter[0] = x
					l.enter[1] = y
				}
				l.Cells[y][x] = cellDefs[string(c)]
			}
			y++
		}
	}

	return nil
}

func (l *Level) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	for y, row := range l.Cells {
		op.GeoM.Reset()
		op.GeoM.Translate(0, float64(y*8))
		for _, cell := range row {
			if cell.Flags()&CellFlagForeground == 0 {
				if cell.ImageID() > 0 {
					img := assets.GetImage(cell.ImageID())
					screen.DrawImage(img, op)
				}
			}
			op.GeoM.Translate(8, 0)
		}
	}
}

func (l *Level) DrawForeground(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	for y, row := range l.Cells {
		op.GeoM.Reset()
		op.GeoM.Translate(0, float64(y*8))
		for _, cell := range row {
			if cell.Flags()&CellFlagForeground != 0 {
				if cell.ImageID() > 0 {
					img := assets.GetImage(cell.ImageID())
					screen.DrawImage(img, op)
				}
			}
			op.GeoM.Translate(8, 0)
		}
	}
}

func (l *Level) GetCoord(x, y float64) (int, int) {
	return int(math.Floor(x / 8)), int(math.Floor(y / 8))
}

func (l *Level) GetCellAt(x, y float64) (Cell, bool) {
	cellx, celly := l.GetCoord(x, y)
	if cellx < 0 || cellx >= 20 {
		return 0, false
	}
	if celly < 0 || celly >= 18 {
		return 0, false
	}

	return l.Cells[celly][cellx], true
}

func (l *Level) AddObject(o Object) {
	l.Objects = append(l.Objects, o)
}

func (l *Level) RemoveObject(o Object) {
	for i, o2 := range l.Objects {
		if o == o2 {
			l.Objects = append(l.Objects[:i], l.Objects[i+1:]...)
			return
		}
	}
}

type CellCollision struct {
	cell   Cell
	cellX  int
	cellY  int
	top    float64
	bottom float64
	left   float64
	right  float64
}

func (l *Level) GetCollidingCell(x, y float64, w, h float64, flag CellFlag) (CellCollision, bool) {
	p := vec2.Vector{X: x, Y: y}
	if c, ok := l.GetCellAt(p.X, p.Y); ok && c.Flags()&flag != 0 {
		cellX, cellY := l.GetCoord(p.X, p.Y)
		coll := CellCollision{
			cell:   l.Cells[cellY][cellX],
			cellX:  cellX,
			cellY:  cellY,
			top:    float64(cellY) * 8,
			bottom: float64(cellY+1) * 8,
			left:   float64(cellX) * 8,
			right:  float64(cellX+1) * 8,
		}
		return coll, true
	} else if c, ok := l.GetCellAt(p.X+w, p.Y); ok && c.Flags()&flag != 0 {
		cellX, cellY := l.GetCoord(p.X+w, p.Y)
		coll := CellCollision{
			cell:   l.Cells[cellY][cellX],
			cellX:  cellX,
			cellY:  cellY,
			top:    float64(cellY) * 8,
			bottom: float64(cellY+1) * 8,
			left:   float64(cellX) * 8,
			right:  float64(cellX+1) * 8,
		}
		return coll, true
	} else if c, ok := l.GetCellAt(p.X, p.Y+h); ok && c.Flags()&flag != 0 {
		cellX, cellY := l.GetCoord(p.X, p.Y+h)
		coll := CellCollision{
			cell:   l.Cells[cellY][cellX],
			cellX:  cellX,
			cellY:  cellY,
			top:    float64(cellY) * 8,
			bottom: float64(cellY+1) * 8,
			left:   float64(cellX) * 8,
			right:  float64(cellX+1) * 8,
		}
		return coll, true
	} else if c, ok := l.GetCellAt(p.X+w, p.Y+h); ok && c.Flags()&flag != 0 {
		cellX, cellY := l.GetCoord(p.X+w, p.Y+h)
		coll := CellCollision{
			cell:   l.Cells[cellY][cellX],
			cellX:  cellX,
			cellY:  cellY,
			top:    float64(cellY) * 8,
			bottom: float64(cellY+1) * 8,
			left:   float64(cellX) * 8,
			right:  float64(cellX+1) * 8,
		}
		return coll, true
	}
	return CellCollision{}, false
}

func (l *Level) FindCells(flags CellFlag) []CellCollision {
	var cells []CellCollision
	for y, row := range l.Cells {
		for x, cell := range row {
			if cell.Flags()&flags != 0 {
				cells = append(cells, CellCollision{
					cell:   cell,
					cellX:  x,
					cellY:  y,
					top:    float64(y) * 8,
					bottom: float64(y+1) * 8,
					left:   float64(x) * 8,
					right:  float64(x+1) * 8,
				})
			}
		}
	}
	return cells
}

func (l *Level) Update() (reqs []Request) {
	gravity := vec2.Vector{X: 0, Y: 0.1}
	drag := float64(0.9)
	// Process physiks
	for _, o := range l.Objects {
		size := o.Size()
		// Add gravity.
		velocity := o.Velocity()
		velocity.Add(gravity)
		// Apply drag.
		velocity.X *= drag
		position := o.Position()
		nextPosition := *position
		nextPosition.Add(*velocity)

		// Check against vertical movement
		if c, ok := l.GetCollidingCell(position.X, nextPosition.Y, size.X, size.Y, CellFlagSolid); ok {
			if nextPosition.Y > position.Y { // Moving down.
				if nextPosition.Y < c.top {
					nextPosition.Y = c.top - size.Y
					velocity.Y = 0
					o.Ground()
				}
			} else if nextPosition.Y < position.Y { // Moving up.
				if nextPosition.Y > c.top {
					nextPosition.Y = c.bottom
					velocity.Y = 0
				}
			}
			reqs = append(reqs, o.TouchCell(c))
		}
		// Check against horizontal movement
		if c, ok := l.GetCollidingCell(nextPosition.X, position.Y-1, size.X, size.Y, CellFlagSolid); ok {
			velocity.X = 0
			nextPosition.X = position.X
			if r := o.TouchCell(c); r != nil {
				if req, ok := r.(RequestSetCell); ok {
					l.Cells[req.y][req.x] = Cell(int32(req.image) | int32(req.flag)<<16)
				}
			}
		}

		// Alright, let's see if we hit anything special.
		if c, ok := l.GetCollidingCell(nextPosition.X, nextPosition.Y, size.X, size.Y, CellFlagClimbable|CellFlagExit|CellFlagHurts|CellFlagDestroyable|CellFlagPushUp); ok {
			if r := o.TouchCell(c); r != nil {
				if req, ok := r.(RequestLevel); ok {
					req.object = o
					req.Level = l.NextLevel
					reqs = append(reqs, req)
				} else if _, ok := r.(RequestReset); ok {
					entrances := l.FindCells(CellFlagEnter)
					if len(entrances) > 0 {
						nextPosition.X = float64(entrances[0].cellX * 8)
						nextPosition.Y = float64(entrances[0].cellY * 8)
					}
				} else if req, ok := r.(RequestSetCell); ok {
					l.Cells[req.y][req.x] = Cell(int32(req.image) | int32(req.flag)<<16)
				} else {
					reqs = append(reqs, r)
				}
			}
		}

		// Oops, did we fall out of the map?
		if nextPosition.Y > 144+size.Y || nextPosition.Y < -size.Y || nextPosition.X > 160+size.X || nextPosition.X < -size.X {
			if r := o.FallOut(); r != nil {
				if _, ok := r.(RequestReset); ok {
					entrances := l.FindCells(CellFlagEnter)
					if len(entrances) > 0 {
						nextPosition.X = float64(entrances[0].cellX * 8)
						nextPosition.Y = float64(entrances[0].cellY * 8)
					}
				} else {
					reqs = append(reqs, r)
				}
			}
		}

		position.X = nextPosition.X
		position.Y = nextPosition.Y
	}

	// Update orbjects
	for _, o := range l.Objects {
		reqs = append(reqs, o.Update(nil))
	}

	return
}

func (l *Level) DrawDebug(screen *ebiten.Image) {
	return
	for _, o := range l.Objects {
		position := o.Position()
		ebitenutil.DrawRect(screen, position.X, position.Y, o.Size().X, o.Size().Y, color.RGBA{255, 0, 0, 255})
		if c, ok := l.GetCollidingCell(position.X, position.Y, o.Size().X, o.Size().Y, CellFlagSolid); ok {
			ebitenutil.DrawRect(screen, c.left, c.top, 8, 8, color.RGBA{0, 255, 0, 255})
		}
	}
}
