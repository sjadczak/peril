package main

import (
	"fmt"

	"github.com/sjadczak/peril/internal/gamelogic"
	"github.com/sjadczak/peril/internal/routing"
)

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) {
	return func(ps routing.PlayingState) {
		defer fmt.Printf("> ")
		gs.HandlePause(ps)
	}
}
