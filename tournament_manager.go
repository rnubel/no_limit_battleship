package no_limit_battleship

import (
  "battleship"
  "code.google.com/p/go-uuid/uuid"
)

type RegisteredPlayer struct {
  Key         string
  Name        string
  CurrentGame *GameRunner
}

type TournamentManager struct {
  players         []*RegisteredPlayer
  gameRunners     map[string]*GameRunner
}

func StartTournament() *TournamentManager {
  tm := TournamentManager{}
  tm.gameRunners  = make(map[string]*GameRunner)

  return &tm
}

func (tm *TournamentManager) CreateGame(p1Key, p2Key string) (runner *GameRunner, error string) {
  rp1, rp2 := tm.LocatePlayer(p1Key), tm.LocatePlayer(p2Key)
  if rp1 == nil {
    return nil, "player_1_does_not_exist"
  } else if rp2 == nil {
    return nil, "player_2_does_not_exist"
  } else if rp1.CurrentGame != nil {
    return nil, "player_1_already_in_game"
  } else if rp2.CurrentGame != nil {
    return nil, "player_2_already_in_game"
  }

  p1, p2 := battleship.Player{Identifier: p1Key}, battleship.Player{Identifier: p2Key}
	game := battleship.CreateGame(10, 10, p1, p2)
  runner = createGameRunner(&game)

  tm.gameRunners[runner.Id] = runner

  rp1.CurrentGame = runner
  rp2.CurrentGame = runner

  return
}

func (tm *TournamentManager) LocateGame(gameID string) *GameRunner {
  return tm.gameRunners[gameID]
}

func (tm *TournamentManager) LocatePlayer(playerKey string) *RegisteredPlayer {
  for i := range(tm.players) {
    if tm.players[i].Key == playerKey {
      return tm.players[i]
    }
  }

  return nil
}

func (tm *TournamentManager) RegisterPlayer(name string) (*RegisteredPlayer, string) {
  for i := range(tm.players) {
    if tm.players[i].Name == name {
      return nil, "name_already_taken"
    }
  }

  key := uuid.New()
  p := &RegisteredPlayer{Key: key, Name: name}

  tm.players = append(tm.players, p)

  return p, ""
}


