package no_limit_battleship

import "battleship"

type PlayerStatus struct {
	PlayerKey     string
	Name          string
	CurrentGameID string
}

type GameStatus struct {
	GameID  string
	Player1 string
	Player2 string
	Phase   string
}

func gameStatus(gr *GameRunner) GameStatus {
	g := gr.Game

	return GameStatus{
		GameID:  gr.Id,
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

func playerStatus(p *RegisteredPlayer) PlayerStatus {
	ps := PlayerStatus{
		PlayerKey: p.Key,
		Name:      p.Name}

	if p.CurrentGame != nil {
		ps.CurrentGameID = p.CurrentGame.Id
	}

	return ps
}
