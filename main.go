package main

import (
	"fmt"
	"math/rand"

	l "github.com/sirupsen/logrus"

	"github.com/RenugaParamalingam/carrom/strikes"
)

func main() {
	names := []string{"p1", "p2", "p3", "p4"}
	resetGame := true

	for resetGame {
		if err := startGame(names[rand.Intn(4):rand.Intn(4)]); err != nil {
			l.WithError(err).Errorln("invalid request")
			resetGame = true

			continue
		}

		resetGame = false
	}
}

func startGame(playerNames []string) error {
	if !strikes.AddPlayersToGame(playerNames) {
		return fmt.Errorf("invalid player names. player names provided: %v", playerNames)
	}

	l.WithField("playerNames", playerNames).Println("Players on board")

	strikeInput := strikes.NewBoard()
	shouldEndGame := false

	redCoinRandomness := []bool{true, false}

	for !shouldEndGame {
		strikeInput <- strikes.Input{
			StrikeCode:   rand.Intn(6),
			CoinsCount:   rand.Intn(10),
			IsRedPockted: redCoinRandomness[rand.Intn(2)],
		}

		shouldEndGame = strikes.IsGameOver()
	}

	return nil
}
