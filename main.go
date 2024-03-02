package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"slices"
	"time"

	"github.com/gorilla/websocket"
)

const boardSize = 10 * 10

func draw() []byte {
	result := make([]int, boardSize)

	for idx := range boardSize {
		result[idx] = rand.Intn(2)
	}

	str, err := json.Marshal(result)
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

	connections := []*websocket.Conn{}

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

		connections = append(connections, conn)

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

		for idx, c := range connections {
			if conn == c {
				fmt.Println("BEFORE", connections)
				connections = slices.Delete(connections, idx, idx+1)
				fmt.Println("AFTER ", connections)
				break
			}
		}

		fmt.Printf("%s DISCONNECT\n", conn.RemoteAddr())
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	go func() {
		for {
			time.Sleep(time.Second * 1)
			// time.Sleep(time.Millisecond * 16)
			fmt.Println(">>> TICK")

			if len(connections) > 0 {

				result := draw()

				fmt.Println(">>> TICK go ")

				for _, conn := range connections {

					if err := conn.WriteMessage(websocket.TextMessage, result); err != nil {
						log.Println("ERROR WriteMessage", err.Error())
					}
				}
			}
		}
	}()

	log.Fatalln(http.ListenAndServe(":8080", nil))
}
