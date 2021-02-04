package carrom

import (
	"fmt"

	l "github.com/sirupsen/logrus"
)

var invalid = "invalid strike input"

// A Input is source for strikes
type Input struct {
	StrikeCode int
	CoinsPocketedCount
}

// StrikeCodeInput is a channel to flow in the strike type and coins for the game.
// When games ends, channel will be closed internally. Don't close channel by yourself.
// Input code and it's strike name are,
// 0 - strike
// 1 - multi strike
// 2 - red strike
// 3 - striker strike
// 4 - defunct
// 5 - no coin pocketed
var StrikeCodeInput chan Input

// NewBoard resets or prepare pre-requesties for game.
// A channel is returned to feed the input for game.
func NewBoard() chan Input {
	CoinsOnBoard = &Coins{
		Red:   1,
		White: 9,
		Black: 9,
	}

	playerIDForTurn = 0

	StrikeCodeInput = make(chan Input, 1)
	go mapInputToStrike()

	return StrikeCodeInput
}

// IsGameOver returns true if any player won or match ended in draw.
func IsGameOver() (gameEnd bool) {
	playersScore := make(map[string]int, 0)

	for _, p := range playersOnBoard {
		playersScore[p.PlayerName] = p.Points
	}

	highestScorer := gethighestScore(playersOnBoard...)

	var winner *Player

	for _, p := range playersOnBoard {
		if highestScorer.Points-p.Points >= 3 && highestScorer.Points >= 5 {
			winner = highestScorer
			break
		}
	}

	if winner != nil {
		l.Printf("\n Player named %q won the game by scoring %v points. \n", winner.PlayerName, winner.Points)
	} else if isBoardEmpty() {
		l.Println("\n Coins exhausted and no players won. Game ends in draw.")
	}

	gameEnd = winner != nil || isBoardEmpty()

	// stop receiving input as game is end.
	if gameEnd {
		printScore(playersOnBoard)
		close(StrikeCodeInput)
	}

	return gameEnd
}

func printScore(players []*Player) {
	fmt.Printf("\n Score board \n -----------------------  \n | Player Name | Score | \n ----------------------- \n")

	for _, p := range players {
		fmt.Printf(" | %-11v | %-5v | \n", p.PlayerName, p.Points)
	}
}

func gethighestScore(players ...*Player) *Player {
	max := players[0]

	for _, p := range players {
		if p.Points > max.Points {
			max = p
		}
	}

	return max
}

func isBoardEmpty() bool {
	return CoinsOnBoard.Red == 0 && CoinsOnBoard.Black == 0 && CoinsOnBoard.White == 0
}

func mapInputToStrike() {
	var err error
	p := playersOnBoard[0]

	for c := range StrikeCodeInput {
		// turn will be passed to next player only if current player
		// successfuly completes his turn by providing valid input.
		if err == nil {
			p = passPlayer()
		}

		l.WithField("input", c).Infoln("input received")

		switch c.StrikeCode {
		case 0:
			err = p.Strike(c.CoinsPocketedCount)
		case 1:
			err = p.MultiStrike(c.CoinsPocketedCount)
		case 2:
			err = p.RedStrike()
		case 3:
			p.StrikerStrike()
		case 4:
			err = p.Defunct(c.Black, c.IsRedPocketed)
		case 5:
			p.NoPocket()
		default:
			l.WithField("strikeCode", c.StrikeCode).Errorln("invalid strike code")
		}
	}
}
