package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/remmakoshino/kirby-inspired-go/internal/game"
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Kirby-Inspired RPG",
		Bounds: pixel.R(0, 0, game.WindowWidth, game.WindowHeight),
		VSync:  true,
	}
	
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	
	// ゲーム作成と実行
	g := game.NewGame(win)
	g.Run()
}

func main() {
	pixelgl.Run(run)
}
