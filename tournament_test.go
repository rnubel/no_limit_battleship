package no_limit_battleship

import (
  "testing"
  "battleship"
)

func endGame(gr *GameRunner, winnerKey string) {
  g := gr.Game
  g.Phase = battleship.FINISHED
  if g.Player1.Identifier == winnerKey {
    g.Winner = g.Player1
  } else {
    g.Winner = g.Player2
  }
}

func TestOneLevelBracketCreation (t *testing.T) {
  tm := StartTournament()
  tm.WinsNeeded = 3
  tm.RegisterPlayer("p1")
  tm.RegisterPlayer("p2")

  if len(tm.players) != 2 { t.Error("Players weren't registered") }

  tm.CreateBracket()
  tm.CreateAvailableGames()
  tm.CreateAvailableGames() // call it twice to make sure this is a no-op

  if len(tm.gameRunners) != 1 { t.Error("One game was not created between the two players") }

  endGame(tm.bracket.rounds[0].matches[0].games[0], tm.players[0].Key) // p1 wins the first game

  tm.CreateAvailableGames()
  if len(tm.gameRunners) != 2 { t.Error("Second game was not created after the first game in a match finished") }

  endGame(tm.bracket.rounds[0].matches[0].games[1], tm.players[0].Key) // p1 wins the second game

  tm.CreateAvailableGames()
  if len(tm.gameRunners) != 3 { t.Error("Third game was not created after the second game in a match finished") }

  endGame(tm.bracket.rounds[0].matches[0].games[2], tm.players[0].Key) // p1 wins the third game and the match
  tm.CreateAvailableGames()

  if len(tm.gameRunners) != 3 { t.Error("Fourth game was created after the match should be over") }
}

func TestByeCreation (t *testing.T) {
  tm := StartTournament()
  tm.RegisterPlayer("p1")
  tm.RegisterPlayer("p2")
  tm.RegisterPlayer("p3")

  tm.CreateBracket()
  tm.CreateAvailableGames()

  if len(tm.bracket.rounds[0].matches) != 2 { t.Error("Two matches were not created for three players") }
  if len(tm.gameRunners) != 1 { t.Error("A bye game was incorrectly started") }
}

func TestSecondLevelBracketCreation (t *testing.T) {
  tm := StartTournament()
  tm.RegisterPlayer("p1")
  tm.RegisterPlayer("p2")
  tm.RegisterPlayer("p3")
  tm.RegisterPlayer("p4")

  tm.CreateBracket()
  tm.CreateAvailableGames()

  if len(tm.bracket.rounds[0].matches) != 2 { t.Error("Two matches were not created for four players") }

  endGame(tm.bracket.rounds[0].matches[0].games[0], tm.bracket.rounds[0].matches[0].playerA.Key)
  endGame(tm.bracket.rounds[0].matches[1].games[0], tm.bracket.rounds[0].matches[1].playerA.Key)
  tm.CreateAvailableGames()
  endGame(tm.bracket.rounds[0].matches[0].games[1], tm.bracket.rounds[0].matches[0].playerA.Key)
  endGame(tm.bracket.rounds[0].matches[1].games[1], tm.bracket.rounds[0].matches[1].playerA.Key)
  tm.CreateAvailableGames()
  endGame(tm.bracket.rounds[0].matches[0].games[2], tm.bracket.rounds[0].matches[0].playerA.Key)
  endGame(tm.bracket.rounds[0].matches[1].games[2], tm.bracket.rounds[0].matches[1].playerA.Key)

  tm.CreateAvailableGames() // this should create a new bracket round

  if len(tm.bracket.rounds) != 2 {
    t.Error("New bracket round was not created after two games in the first level finished")
  }

  if len(tm.bracket.rounds[1].matches) != 1 {
    t.Error("Expected one match created in second round, got", len(tm.bracket.rounds[1].matches))
  }
}
