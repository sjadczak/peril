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
	fmt.Println("Peril game server starting...")
	rabbitConnString := "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(rabbitConnString)
	if err != nil {
		log.Fatalf("unable to connect to RabbitMQ: %v\n", err)
	}
	defer conn.Close()
	fmt.Println("Peril game server connected to RabbitMQ!")

	// Create channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("unable to get RabbitMQ channel: %v\n", err)
	}
	pubsub.PublishJSON(
		ch,
		routing.ExchangePerilDirect,
		routing.PauseKey,
		routing.PlayingState{IsPaused: true},
	)

	// declare & bind logs queue
	_, queue, err := pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilTopic,
		routing.GameLogSlug,
		routing.GameLogSlug+".*",
		pubsub.Durable,
	)
	if err != nil {
		log.Fatalf("unable to subscribe to game_logs: %v\n", err)
	}
	fmt.Printf("Queue %v declared and bound!\n", queue.Name)

	gamelogic.PrintServerHelp()
	isRunning := true
	for isRunning {
		cmd := gamelogic.GetInput()

		if len(cmd) == 0 {
			continue
		}

		switch cmd[0] {
		case "pause":
			fmt.Println("Publishing pause game state")
			err := pubsub.PublishJSON(
				ch,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{IsPaused: true},
			)
			if err != nil {
				log.Printf("could not publish: %v\n", err)
			}
		case "resume":
			fmt.Println("Publishing resume game state")
			pubsub.PublishJSON(
				ch,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{IsPaused: false},
			)
			if err != nil {
				log.Printf("could not publish: %v\n", err)
			}
		case "quit":
			fmt.Println("Goodbye!")
			isRunning = false
		default:
			fmt.Println("Unknown command")
		}
	}
}
