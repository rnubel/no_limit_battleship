package no_limit_battleship

import (
	"battleship"
	"code.google.com/p/go-uuid/uuid"
  "math/rand"
)

type BracketMatch struct {
  playerA     *RegisteredPlayer
  playerB     *RegisteredPlayer
  games       []*GameRunner
  winner      *RegisteredPlayer
  nextMatch   *BracketMatch
}

type BracketRound struct {
  num       int
  matches   []*BracketMatch
}

type Bracket struct {
  rounds    []*BracketRound
}

type RegisteredPlayer struct {
	Key         string
	Name        string
	CurrentGame *GameRunner
}

type TournamentManager struct {
	players     []*RegisteredPlayer
	gameRunners map[string]*GameRunner

  bracket     *Bracket

  WinsNeeded  int
}

func StartTournament() *TournamentManager {
	tm := TournamentManager{WinsNeeded: 3}
	tm.gameRunners = make(map[string]*GameRunner)

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

	return tm.createGame(rp1, rp2), ""
}

func (tm *TournamentManager) createGame(rp1, rp2 *RegisteredPlayer) (runner *GameRunner) {
	p1, p2 := battleship.Player{Identifier: rp1.Key, Name: rp1.Name},
              battleship.Player{Identifier: rp2.Key, Name: rp2.Name}
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
	for i := range tm.players {
		if tm.players[i].Key == playerKey {
			return tm.players[i]
		}
	}

	return nil
}

func (tm *TournamentManager) RegisterPlayer(name string) (*RegisteredPlayer, string) {
	for i := range tm.players {
		if tm.players[i].Name == name {
			return nil, "name_already_taken"
		}
	}

	key := uuid.New()
	p := &RegisteredPlayer{Key: key, Name: name}

	tm.players = append(tm.players, p)

	return p, ""
}

func (tm *TournamentManager) CreateBracket() {
  // Create bracket
  tm.bracket = &Bracket{}

  // Create first round
  round := &BracketRound{num: 1}
  tm.bracket.rounds = append(tm.bracket.rounds, round)

  tm.createMatchesForPlayers(round, tm.players, true)
}

func (tm *TournamentManager) createMatchesForPlayers(round *BracketRound, _players []*RegisteredPlayer, createBye bool) {
  // Knuth shuffle the players
  players := []*RegisteredPlayer{}
  perm := rand.Perm(len(_players))
  for _, v := range perm {
    players = append(players, _players[v])
  }

  // Create matches
  for len(players) >= 2 {
    p1, p2, ps := players[0], players[1], players[2:]
    match := &BracketMatch{playerA: p1, playerB: p2}
    round.matches = append(round.matches, match)

    players = ps
  }

  // Bye?
  if len(players) > 0 && createBye {
    // create bye
    p := players[0]
    match := &BracketMatch{playerA: p, winner: p}
    round.matches = append(round.matches, match)
  }
}

func (tm *TournamentManager) CreateAvailableGames() {
  // Start a new game for unfinished matches
  for i, round := range(tm.bracket.rounds) {
    matchesOverAndPending := []*BracketMatch{}
    numMatchesInProgress := 0

    for _, match := range(round.matches) {
      if match.winner == nil {
        // are all their games over?
        gamesOver, aWins, bWins := true, 0, 0
        for _, gr := range(match.games) {
          if gr.Game.IsOver() {
            if gr.WinnerIdentifier() == match.playerA.Key {
              aWins += 1
            } else {
              bWins += 1
            }
          } else {
            gamesOver = false
            break
          }
        }

        if gamesOver {
          if aWins >= tm.WinsNeeded && aWins > bWins {
            match.winner = match.playerA
            matchesOverAndPending = append(matchesOverAndPending, match)
          } else if bWins >= tm.WinsNeeded && bWins > aWins {
            match.winner = match.playerB
            matchesOverAndPending = append(matchesOverAndPending, match)
          } else {
            // create a new game
            var game *GameRunner
            if len(match.games) % 2 == 0 {
              game = tm.createGame(match.playerA, match.playerB)
            } else {
              game = tm.createGame(match.playerB, match.playerA)
            }
            match.games = append(match.games, game)
            numMatchesInProgress += 1
          }
        } else {
          numMatchesInProgress += 1
        }
      } else if match.nextMatch == nil {
        matchesOverAndPending = append(matchesOverAndPending, match)
      }
    }

    if len(matchesOverAndPending) >= 2 {
      // we need to create a new match in the next round. Does it exist?
      if len(tm.bracket.rounds) <= i + 1 {
        tm.bracket.rounds = append(tm.bracket.rounds, &BracketRound{num: round.num + 1})
      }

      nextRound := tm.bracket.rounds[len(tm.bracket.rounds) - 1]

      players := []*RegisteredPlayer{}
      for _, match := range(matchesOverAndPending) {
        players = append(players, match.winner)
      }

      tm.createMatchesForPlayers(nextRound, players, numMatchesInProgress == 0)
    }
  }
}
