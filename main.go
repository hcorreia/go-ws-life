package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const width = 10
const boardSize = width * width
const ALIVE = 1
const DEAD = 0

type Board = [boardSize]int

func life(board Board) Board {
	var newBoard Board

	for idx, state := range board {
		x := idx % width
		y := idx / width

		neighboursAlive := 0

		// fmt.Println("IDX", idx)
		// fmt.Println("X Y", x, y)

		for _, xi := range [3]int{-1, 0, 1} {
			for _, yi := range [3]int{-1, 0, 1} {

				// Ignore self
				if xi == 0 && yi == 0 {
					continue
				}

				// Check boundries
				if ((x+xi >= 0) && (x+xi < width)) && ((y+yi >= 0) && (y+yi < width)) {
					// fmt.Println("CHECK", (x + xi), (y + yi))

					// Check is Alive
					if board[(y+yi)*width+(x+xi)] == ALIVE {
						neighboursAlive += 1
					}
				}
			}
		}

		newBoard[idx] = board[idx]

		if state == ALIVE {
			if neighboursAlive < 2 {
				newBoard[idx] = DEAD
			} else if neighboursAlive > 3 {
				newBoard[idx] = DEAD
			}
		} else {
			if neighboursAlive == 3 {
				newBoard[idx] = ALIVE
			}
		}

		// fmt.Println("neighboursAlive", neighboursAlive)
		// fmt.Println("")

		// fmt.Println("IDX", idx)
		// fmt.Println("isAlive", state)
		// fmt.Println("X Y", x, y)
	}

	return newBoard
}

func genRandomBoard() Board {
	var result Board

	for idx := range result {
		result[idx] = rand.Intn(2)
	}

	// fmt.Println(">>>>", result)

	return result
}

func genBlinkerBoard() Board {
	return Board{
		0, 0, 0, 0, 0,
		0, 0, 1, 0, 0,
		0, 0, 1, 0, 0,
		0, 0, 1, 0, 0,
		0, 0, 0, 0, 0,
	}
}

func genToadBoard() Board {
	return Board{
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 1, 1, 1, 0, 0, 0, 0,
		0, 0, 0, 0, 1, 1, 1, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}
}

func draw(board Board) []byte {
	str, err := json.Marshal(board)
	if err != nil {
		log.Fatalln(err)
	}

	// fmt.Println(">>>>", string(str))

	return str
}

func main() {
	fmt.Println("Lets GO...")

	homeTempl, err := template.New("homepage.templ").ParseFiles("templates/homepage.templ")
	if err != nil {
		log.Fatalln(err)
	}

	connections := map[*websocket.Conn]bool{}
	// gameState := genRandomBoard()
	// gameState := genBlinkerBoard()
	gameState := genToadBoard()

	// fmt.Println(">>>>", gameState)

	http.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.String())

		err = homeTempl.Execute(w, struct {
			URL string
		}{
			URL: r.URL.String(),
		})
		if err != nil {
			log.Fatalln(err)
		}

	})

	upgrader := websocket.Upgrader{}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {

		log.Println(r.Method, r.URL.String())

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("ERROR", err.Error())

			w.WriteHeader(http.StatusBadRequest)
			// fmt.Fprintf(w, "400 - Unable to Upgrade")

			return
		}

		connections[conn] = true

		fmt.Println("CONS", connections)

		// Echo
		for {
			// Read message from browser
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				break
			}

			// Print the message to the console
			fmt.Printf("%s TYPE:%d sent: %s\n", conn.RemoteAddr(), msgType, string(msg))

			// Write message back to browser
			if err = conn.WriteMessage(msgType, msg); err != nil {
				break
			}
		}

		fmt.Println("BEFORE", connections)
		delete(connections, conn)
		fmt.Println("AFTER ", connections)

		fmt.Printf("%s DISCONNECT\n", conn.RemoteAddr())
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func(ctx context.Context) {
		backof := 0
		tickTime := time.Second * 1
		// tickTime := time.Millisecond * 16

		sleep := tickTime

		ticker := time.NewTicker(sleep)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				ticker.Stop()

				fmt.Println(">>> TICK")

				if len(connections) > 0 {
					if backof != 0 {
						backof = 0
						sleep = tickTime
					}

					fmt.Println(">>> TICK go")

					result := draw(gameState)

					for conn := range connections {
						if err := conn.WriteMessage(websocket.TextMessage, result); err != nil {
							log.Println("ERROR WriteMessage", err.Error())

							// TODO: maybe remove this conn
						}
					}

					// fmt.Println("1 >>>>", gameState)
					gameState = life(gameState)
					// fmt.Println("2 >>>>", gameState)

				} else if backof <= 5 {
					backof += 1
					sleep = time.Second * (2 * time.Duration(backof))

				}

				if backof > 0 {
					log.Println(">>> TICK SLEEP", sleep)
				}

				ticker.Reset(sleep)

			case <-ctx.Done():
				fmt.Println(">>> TICK DONE")

				return
			}
		}
	}(ctx)

	log.Fatalln(http.ListenAndServe(":8080", nil))
}
