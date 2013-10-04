package no_limit_battleship

import (
  "battleship"
  "fmt"
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
  ok          bool
  err         string
}

type GameRunner struct {
  Id            string
  Game          *battleship.Game
  TurnInput     chan Turn
  ResultOutput  chan Result
}

func createGameRunner(g *battleship.Game) *GameRunner {
  gr := &GameRunner{Id: uuid.New(), Game: g}

  go gr.run()

  return gr
}

func (gr *GameRunner) run() {
  select {
  case turn := <-gr.TurnInput:
    fmt.Println("Received turn: ", turn)

    r := Result{ok: true}
    gr.ResultOutput <- r
  }
}
