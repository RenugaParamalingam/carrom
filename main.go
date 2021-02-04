package main

import (
	"fmt"
	"math/rand"
	"time"

	l "github.com/sirupsen/logrus"

	"github.com/RenugaParamalingam/carrom/carrom"
)

func main() {
	names := []string{"p1", "p2", "p3", "p4"}
	resetGame := true

	rand.Seed(int64(time.Now().Nanosecond()))

	for resetGame {
		start := rand.Intn(4)
		end := rand.Intn(4)

		if end < start {
			continue
		}

		if err := startGame(names[start:end]); err != nil {
			l.WithError(err).Errorln("invalid request")

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

	rand.Seed(int64(time.Now().Nanosecond()))

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
