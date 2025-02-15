package main

import "github.com/hajimehoshi/ebiten/v2"

func main() {
	g := &Game{}
	g.State = NewStateIntro()

	ebiten.SetWindowSize(160*4, 144*4)
	ebiten.SetWindowTitle("BsK!!!")

	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
