package life

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"image"
	"image/color"
	"image/png"
	"log"
	"math/rand"
)

const ALIVE = 1
const DEAD = 0

var aliveColor = color.RGBA{51, 51, 51, 255}
var deadColor = color.RGBA{204, 204, 204, 255}

type Life struct {
	Width int   `json:"width"`
	Size  int   `json:"size"`
	Board []int `json:"board"`
}

func (life *Life) Life() {
	newBoard := make([]int, life.Size)

	for idx, state := range life.Board {
		x := idx % life.Width
		y := idx / life.Width

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
				if ((x+xi >= 0) && (x+xi < life.Width)) && ((y+yi >= 0) && (y+yi < life.Width)) {
					// fmt.Println("CHECK", (x + xi), (y + yi))

					// Check is Alive
					if life.Board[(y+yi)*life.Width+(x+xi)] == ALIVE {
						neighboursAlive += 1
					}
				}
			}
		}

		newBoard[idx] = life.Board[idx]

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

	life.Board = newBoard

	// return newBoard
}

func (life *Life) Draw() []byte {
	str, err := json.Marshal(life)
	if err != nil {
		log.Fatalln(err)
	}

	return str
}

func (life *Life) DrawImageDataUrl() ([]byte, error) {
	m := image.NewRGBA(image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{life.Width, life.Width},
	})

	for idx, state := range life.Board {
		x := idx % life.Width
		y := idx / life.Width

		if state == ALIVE {
			m.SetRGBA(x, y, aliveColor)
		} else {
			m.SetRGBA(x, y, deadColor)
		}
	}

	// b := []byte{}
	w := bytes.NewBuffer([]byte{})

	err := png.Encode(w, m)
	if err != nil {
		log.Fatalln(err)

		return nil, err
	}

	b := w.Bytes()

	b64 := make([]byte, base64.StdEncoding.EncodedLen(len(b)))
	base64.RawStdEncoding.Encode(b64, b)

	result := bytes.Join([][]byte{[]byte("data:image/png;base64,"), b64}, []byte{})

	return result, nil
}

func NewLife(width int) *Life {
	return &Life{
		Width: width,
		Size:  width * width,
		Board: make([]int, width*width),
	}
}

func NewLifeRandom(width int) *Life {
	board := make([]int, width*width)

	for idx := range board {
		board[idx] = rand.Intn(2)
	}

	return &Life{
		Width: width,
		Size:  width * width,
		Board: board,
	}
}

func NewLifeBeacon() *Life {
	return &Life{
		Width: 10,
		Size:  10 * 10,
		Board: []int{
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 1, 1, 0, 0, 0, 0, 0,
			0, 0, 0, 1, 1, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 1, 1, 0, 0, 0,
			0, 0, 0, 0, 0, 1, 1, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		},
	}
}

func NewLifeBlinker() *Life {
	return &Life{
		Width: 5,
		Size:  5 * 5,
		Board: []int{
			0, 0, 0, 0, 0,
			0, 0, 1, 0, 0,
			0, 0, 1, 0, 0,
			0, 0, 1, 0, 0,
			0, 0, 0, 0, 0,
		},
	}
}

func NewLifePentaDecathlon() *Life {
	return &Life{
		Width: 18,
		Size:  18 * 18,
		Board: []int{
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		},
	}
}

func NewLifeToad() *Life {
	return &Life{
		Width: 10,
		Size:  10 * 10,
		Board: []int{
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
		},
	}
}
