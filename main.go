package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hcorreia/go-ws-life/life"
)

const tickTime = time.Millisecond * 33

// const tickTime = time.Millisecond * 16
// const tickTime = time.Millisecond * 120

const boardSize = 100

func main() {
	fmt.Println("Lets GO...")

	homeTempl, err := template.New("homepage.templ").ParseFiles("templates/homepage.templ")
	if err != nil {
		log.Fatalln(err)
	}

	connections := map[*websocket.Conn]bool{}
	// game := life.NewLife(boardSize)
	game := life.NewLifeRandom(boardSize)
	// game := life.NewLifeBeacon()
	// game := life.NewLifeBlinker()
	// game := life.NewLifePentaDecathlon()
	// game := life.NewLifeToad()

	// fmt.Println(">>>>", game.Board)

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
		sleep := tickTime

		ticker := time.NewTicker(sleep)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				ticker.Stop()

				// fmt.Println(">>> TICK")

				if len(connections) > 0 {
					if backof != 0 {
						backof = 0
						sleep = tickTime
					}

					// fmt.Println(">>> TICK go")

					result := game.Draw()

					// fmt.Println(">>>", string(result))

					for conn := range connections {
						if err := conn.WriteMessage(websocket.TextMessage, result); err != nil {
							log.Println("ERROR WriteMessage", err.Error())

							// TODO: maybe remove this conn
						}
					}

					// fmt.Println("1 >>>>", game.Board)
					game.Life()
					// fmt.Println("2 >>>>", game.Board)

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
