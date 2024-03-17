package life_test

import (
	"fmt"
	"testing"

	"github.com/hcorreia/go-ws-life/life"
)

func BenchmarkNewLifeNativeRustRandom(b *testing.B) {
	for _, boardSize := range []int{1, 10, 100, 1_000, 1_440, 10_000} {
		b.Run(fmt.Sprintf("NewLifeNativeRustRandom-Board-%d", boardSize), func(b *testing.B) {
			life.NewLifeNativeRust(boardSize)
		})
	}
}

func BenchmarkNativeRustDrawImageDataUrl(b *testing.B) {
	for _, boardSize := range []int{1, 10, 100, 1_000, 1_440, 10_000} {
		game := life.NewLifeNativeRust(boardSize)

		b.Run(fmt.Sprintf("NativeRustDrawImageDataUrl-Board-%d", boardSize), func(b *testing.B) {
			game.DrawImageDataUrl()
		})
	}
}

func BenchmarkLifeNativeRust(b *testing.B) {
	for _, boardSize := range []int{1, 10, 100, 1_000, 1_440, 10_000} {
		game := life.NewLifeNativeRust(boardSize)

		b.Run(fmt.Sprintf("LifeNativeRust-Board-%d", boardSize), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				game.Life()
			}
		})
	}
}
