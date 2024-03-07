package life_test

import (
	"fmt"
	"testing"

	"github.com/hcorreia/go-ws-life/life"
)

func BenchmarkLife(b *testing.B) {
	for _, boardSize := range []int{1, 10, 100, 1_000, 10_000} {
		game := life.NewLifeRandom(boardSize)

		b.Run(fmt.Sprintf("Life-%d", boardSize), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				game.Life()
			}
		})
	}
}
