package main

import (
	"fmt"
	"math/rand"

	l "github.com/sirupsen/logrus"

	"github.com/RenugaParamalingam/carrom/carrom"
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
	if !carrom.AddPlayersToGame(playerNames) {
		return fmt.Errorf("invalid player names. player names provided: %v", playerNames)
	}

	l.WithField("playerNames", playerNames).Println("Players on board")

	strikeInput := carrom.NewBoard()
	shouldEndGame := false

	redCoinRandomness := []bool{true, false}

	for !shouldEndGame {
		strikeInput <- carrom.Input{
			StrikeCode: rand.Intn(6),
			CoinsPocketedCount: carrom.CoinsPocketedCount{
				Black:         rand.Intn(10),
				White:         rand.Intn(10),
				IsRedPocketed: redCoinRandomness[rand.Intn(2)],
			},
		}

		shouldEndGame = carrom.IsGameOver()
	}

	return nil
}
