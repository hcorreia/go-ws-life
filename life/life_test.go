package life_test

import (
	"fmt"
	"testing"

	"github.com/hcorreia/go-ws-life/life"
)

func BenchmarkNewLifeRandom(b *testing.B) {
	for _, boardSize := range []int{1, 10, 100, 1_000, 1_440, 10_000} {
		b.Run(fmt.Sprintf("NewLifeRandom-Board-%d", boardSize), func(b *testing.B) {
			life.NewLifeRandom(boardSize)
		})
	}
}

func BenchmarkDrawImageDataUrl(b *testing.B) {

	for _, boardSize := range []int{1, 10, 100, 1_000, 1_440, 10_000} {
		game := life.NewLifeRandom(boardSize)

		b.Run(fmt.Sprintf("DrawImageDataUrl-Board-%d", boardSize), func(b *testing.B) {
			game.DrawImageDataUrl()
		})
	}
}

func BenchmarkLife(b *testing.B) {
	for _, boardSize := range []int{1, 10, 100, 1_000, 1_440, 10_000} {
		game := life.NewLifeRandom(boardSize)

		b.Run(fmt.Sprintf("Life-Board-%d", boardSize), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				game.Life()
			}
		})
	}
}
