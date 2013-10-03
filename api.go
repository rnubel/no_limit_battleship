package no_limit_battleship

import "battleship"

type GameStatus struct {
	Player1 string
	Player2 string
	Phase   string
}

func gameStatus(g *battleship.Game) GameStatus {
	return GameStatus{
		Player1: g.Player1.Identifier,
		Player2: g.Player2.Identifier,
		Phase:   phaseName(g.Phase)}
}

func phaseName(phase battleship.GamePhase) string {
	switch phase {
	case battleship.NOTSTARTED:
		return "not_started"
	case battleship.PLACEMENT:
		return "placement"
	case battleship.BATTLE:
		return "battle"
	case battleship.FINISHED:
		return "finished"
	default:
		return ""
	}
}
