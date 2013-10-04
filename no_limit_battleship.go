package no_limit_battleship

import (
	"battleship"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
	"time"
)


func printBoard(b *battleship.Board) (output string) {
	for y := 0; y < b.Height; y++ {
		cells := []string{}
		for x := 0; x < b.Width; x++ {
			var val string

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

// stores the state of the server in a shared construct.
type ServerState struct {
  tm  *TournamentManager
}

// decorator that times and logs the request
type BaseRequestHandler struct {
	router *mux.Router
}

func (h BaseRequestHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	t0 := time.Now()
	fmt.Printf("Handling request [%s]\n", req.URL.String())
	req.ParseForm()

	w.Header().Set("Content-Type", "application/json")

	h.router.ServeHTTP(w, req)

	fmt.Printf("Request took %fms to serve.\n", float64(time.Since(t0).Nanoseconds())/1e6)
}

func renderJSON(w http.ResponseWriter, json []byte, err error) {
	if err != nil {
		http.Error(w, fmt.Sprintf("{\"error\": \"%s\"}", err), 500)
    return
	}

	fmt.Fprintf(w, "%s", json)
}

func renderErrorJSON(w http.ResponseWriter, err string, status int) {
  http.Error(w, fmt.Sprintf("{\"error\": \"%s\"}", err), status)
  fmt.Printf("Returned error: %s\n", err)
}


// root page, doesn't do much
func rootHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "{\"status\": \"ok\"}")
}

// creates a new game. might remove this, probably done internally.
func createGameHandler(w http.ResponseWriter, req *http.Request, ss *ServerState) {
	gameRunner, gameErr := ss.tm.CreateGame(req.FormValue("player1"), req.FormValue("player2"))

  if gameErr != "" {
    renderErrorJSON(w, gameErr, 422)
    return
  }

	json, err := json.Marshal(gameStatus(gameRunner))
  renderJSON(w, json, err)
}


// returns the game status
func gameStatusHandler(w http.ResponseWriter, req *http.Request, ss *ServerState) {
  // locate the requested game
  vars := mux.Vars(req)
  gameID := vars["gameID"]

  gameRunner := ss.tm.LocateGame(gameID)

  if gameRunner == nil {
    renderErrorJSON(w, "game_not_found", 404)
    return
  }

  json, err := json.Marshal(gameStatus(gameRunner))
  renderJSON(w, json, err)
}

// registers a new player, fetches a unique ID for them
func registerPlayerHandler(w http.ResponseWriter, req *http.Request, ss *ServerState) {
  name := req.FormValue("name")
  player, regErr := ss.tm.RegisterPlayer(name)

  if regErr != "" {
    renderErrorJSON(w, regErr, 422)
    return
  }

  json, err := json.Marshal(playerStatus(player))
  renderJSON(w, json, err)
}

// fetches a player's status from their key
func playerStatusHandler(w http.ResponseWriter, req *http.Request, ss *ServerState) {
  vars := mux.Vars(req)
  playerKey := vars["playerKey"]

  player := ss.tm.LocatePlayer(playerKey)

  if player == nil {
    http.Error(w, "{\"error\": \"invalid_player_key\"}", 422)
    return
  }

  json, err := json.Marshal(playerStatus(player))
  renderJSON(w, json, err)
}

// main server routine
func NoLimitBattleship() {
	fmt.Println("Starting up the server...")

  tm          := StartTournament()
  serverState := &ServerState{tm: tm}

  // inject the server state
  makeHandler := func(fn func(http.ResponseWriter, *http.Request, *ServerState)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
      fn(w, r, serverState)
    }
  }

	r := mux.NewRouter()
	r.HandleFunc("/", rootHandler)
  r.HandleFunc("/players",              makeHandler(registerPlayerHandler)).Methods("POST")
  r.HandleFunc("/players/{playerKey}",  makeHandler(playerStatusHandler)).Methods("GET")
	r.HandleFunc("/games",                makeHandler(createGameHandler)).Methods("POST")
	r.HandleFunc("/games/{gameID}",       makeHandler(gameStatusHandler)).Methods("GET")

	// we want all requests to be logged and timed
	bh := BaseRequestHandler{router: r}
	http.Handle("/", bh)

	http.ListenAndServe(":8080", nil)
}
