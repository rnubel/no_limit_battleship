package no_limit_battleship

import (
	"battleship"
	"code.google.com/p/go-uuid/uuid"
)

type Turn struct {
	actionType string
	x          int
	y          int
	size       int
	horizontal bool
}

type Result struct {
	ok  bool
	err string
}

type GameRunner struct {
	Id           string
	Game         *battleship.Game
	TurnInput    chan Turn
	ResultOutput chan Result
}

func createGameRunner(g *battleship.Game) *GameRunner {
	gr := &GameRunner{Id: uuid.New(), Game: g}

	return gr
}

func (gr *GameRunner) start() {
  gr.Game.Start()
}

func (gr *GameRunner) WinnerIdentifier() string {
  return gr.Game.Winner.Identifier
}
