package no_limit_battleship

import "battleship"

type PlayerStatus struct {
	PlayerKey     string
	Name          string
	CurrentGameID string
}

type SalvoResult struct {
  Hits  int
}

type GameStatus struct {
	GameID        string
	Player1       string
	Player2       string
	Phase         string
  CurrentPlayer string
}

func gameStatus(gr *GameRunner) GameStatus {
	g := gr.Game
	gs := GameStatus{
    Player1:  g.Player1.Name,
    Player2:  g.Player2.Name,
		GameID:  gr.Id,
		Phase:   phaseName(g.Phase)}

  if g.Phase == battleship.BATTLE {
    gs.CurrentPlayer = g.CurrentPlayer.Name
  }

  return gs
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
