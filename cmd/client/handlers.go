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

func handlerMove(gs *gamelogic.GameState) func(gamelogic.ArmyMove) {
	return func(move gamelogic.ArmyMove) {
		defer fmt.Printf("> ")
		gs.HandleMove(move)
	}
}
