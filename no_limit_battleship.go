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

type Action struct {
	actionType string
	x          int
	y          int
	size       int
	horizontal bool
}

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

func rootHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "{\"status\": \"ok\"}")
}

func createGameHandler(w http.ResponseWriter, req *http.Request) {
	p1, p2 := battleship.Player{Identifier: req.FormValue("player1")},
		battleship.Player{Identifier: req.FormValue("player2")}

	game := battleship.CreateGame(10, 10, p1, p2)

	json, err := json.Marshal(gameStatus(&game))
	if err != nil {
		http.Error(w, fmt.Sprintf("{\"error\": \"%s\"}", err), 500)
	}

	fmt.Fprintf(w, "%s", json)
}

func NoLimitBattleship() {
	fmt.Println("Starting up the server...")

	//  board := initializeBoard()
	//  input, confirmation := manageBoard(board)

	r := mux.NewRouter()
	r.HandleFunc("/", rootHandler)
	r.HandleFunc("/game", createGameHandler).Methods("POST")

	// we want all requests to be logged and timed
	bh := BaseRequestHandler{router: r}
	http.Handle("/", bh)

	http.ListenAndServe(":8080", nil)
}
