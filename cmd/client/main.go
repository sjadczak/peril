package main

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sjadczak/peril/internal/gamelogic"
	"github.com/sjadczak/peril/internal/pubsub"
	"github.com/sjadczak/peril/internal/routing"
)

func main() {
	fmt.Println("Starting Peril client...")
	rabbitConnString := "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(rabbitConnString)
	if err != nil {
		log.Fatalf("unable to connect to RabbitMQ: %v\n", err)
	}
	defer conn.Close()
	fmt.Println("Peril game client connected to RabbitMQ!")

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("invalid username: %v\n", err)
	}
	gs := gamelogic.NewGameState(username)

	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+gs.GetUsername(),
		routing.PauseKey,
		pubsub.Transient,
		handlerPause(gs),
	)
	if err != nil {
		log.Fatalf("unable to subscribe to pause: %v\n", err)
	}

	isRunning := true
	for isRunning {
		cmds := gamelogic.GetInput()

		if len(cmds) == 0 {
			continue
		}

		switch cmds[0] {
		case "spawn":
			err := gs.CommandSpawn(cmds)
			if err != nil {
				log.Printf("unable to spawn: %v\n", err)
				continue
			}
		case "move":
			_, err := gs.CommandMove(cmds)
			if err != nil {
				log.Printf("unable to move: %v\n", err)
				continue
			}
		case "status":
			gs.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			fmt.Println("Spamming not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			isRunning = false
		default:
			fmt.Println("unknown command")
		}
	}
}
