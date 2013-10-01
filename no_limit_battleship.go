package no_limit_battleship

import (
  "battleship"
  "net/http"
  "fmt"
  "strings"
  "strconv"
)

type Action struct {
  actionType    string
  x             int
  y             int
  size          int
  horizontal    bool
}

func printBoard(b *battleship.Board) (output string) {
  for y := 0; y < b.Height; y++ {
    cells := []string{}
    for x := 0; x < b.Width; x++ {
      var val string;

      if b.HitAt(x, y) {
        val = "[X]"
      } else if b.MissAt(x, y) {
        val = " O "
      } else if b.ShipAt(x, y) {
        val = "[ ]"
      } else {
        val = " . "
      }

      cells = append(cells, val)
    }

    output += strings.Join(cells, "") + "\n"
  }
  return
}

func initializeBoard() *battleship.Board {
  board := battleship.Board{Width: 10, Height: 10}
  return &board
}

func manageBoard(board *battleship.Board) (chan Action, chan bool) {
  actionChannel := make(chan Action)
  confChannel := make(chan bool)

  go func() { // Do this all in parallel.
    for a := range(actionChannel) {
      if a.actionType == "shoot" {
        board.RecordShot(a.x, a.y)
      } else if a.actionType == "place" {
        board.PlaceShip(a.x, a.y, a.size, a.horizontal)
      }
      confChannel <- true
    }
  }()

  return actionChannel, confChannel;
}

func NoLimitBattleship() {
  fmt.Println("Starting up the server...")

  board := initializeBoard()
  input, confirmation := manageBoard(board)

  http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
    fmt.Println("Handling request [%s]", req.URL.String())
    req.ParseForm()
    if req.Form.Get("x") != "" {
      x, _ := strconv.Atoi(req.Form.Get("x"))
      y, _ := strconv.Atoi(req.Form.Get("y"))
      a := Action{x: x, y: y}

      if req.Form.Get("action") == "shoot" {
        a.actionType = "shoot"
      } else if req.Form.Get("action") == "place" {
        a.actionType = "place"
        size, err1 := strconv.Atoi(req.Form.Get("size"))
        horizontal, err2 := strconv.ParseBool(req.Form.Get("horizontal"))

        if err1 != nil || err2 != nil {
          fmt.Println("INVALID REQUEST")
          fmt.Fprintf(w, "Invalid request!")
          return;
        } else {
          a.size = size
          a.horizontal = horizontal
        }
      }

      input <- a
      result := <-confirmation
      fmt.Println(result)
    }

    fmt.Fprintf(w, "Board:\n%s", printBoard(board))
  })
  http.ListenAndServe(":8080", nil)
}
