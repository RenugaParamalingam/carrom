package strikes

import (
	"fmt"

	l "github.com/sirupsen/logrus"
)

// Player is details of the player in game.
type Player struct {
	PlayerName    string
	Points        int
	FoulCount     int
	NoPocketCount int
}

// Coins holds different coins on board.
type Coins struct {
	Red   int
	Black int
}

var (
	red     = "red"
	black   = "black"
	invalid = "invalid strike input"
)

// A Input is source for strikes
type Input struct {
	StrikeCode   int
	CoinsCount   int
	IsRedPockted bool
}

// CoinsOnBoard gives coins and it's current count.
var CoinsOnBoard *Coins
var playersOnBoard []*Player
var playerIDForTurn int

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
		Black: 9,
	}

	playerIDForTurn = 0

	StrikeCodeInput = make(chan Input, 1)
	go mapInputToStrike()

	return StrikeCodeInput
}

func newPlayer(name string) *Player {
	return &Player{
		PlayerName: name,
	}
}

// Strike adds a point to player removes the pocketed black coin out of game.
func (p *Player) Strike() {
	p.Points++
	removeCoin(black, 1)
}

// MultiStrike takes count of coins pocketed and a flag to determine red coin is pocketed.
// And error is returned incase of invalid coins count or red pocketed flag.
func (p *Player) MultiStrike(coinsCount int, isRedPockted bool) error {
	if coinsCount > CoinsOnBoard.Black || coinsCount == 0 || (isRedPockted && CoinsOnBoard.Red == 0) {
		l.WithFields(l.Fields{
			"blackCoinOnBoard":    CoinsOnBoard.Black,
			"coinsCountRequested": coinsCount,
			"isRedPockted":        isRedPockted,
		}).Errorln("invalid multistrike request. Ignoring request")

		return fmt.Errorf(invalid)

	}

	p.Points += 2

	if isRedPockted {
		removeCoin(red, 1)
		removeCoin(black, coinsCount-1)

		return nil
	}

	removeCoin(black, 2)

	return nil
}

// RedStrike returns error if red coin is already out of game.
func (p *Player) RedStrike() error {
	if CoinsOnBoard.Red <= 0 {
		l.Errorln("red coin is not in board. Ignoring request.")

		return fmt.Errorf(invalid)
	}

	p.Points += 3
	removeCoin(red, 1)

	return nil
}

// StrikerStrike adds a foul count as player loses a point.
func (p *Player) StrikerStrike() {
	p.Points--
	p.foul()
}

// Defunct takes count of coins pocketed and a flag to determine red coin is pocketed
// and removes coins out of game provided in.
// An error is returned incase of invalid coins count or red pocketed flag.
func (p *Player) Defunct(coinsCount int, isRedPockted bool) error {
	if CoinsOnBoard.Black < coinsCount || (coinsCount == 0 ||
		(isRedPockted && CoinsOnBoard.Red == 0)) {
		l.WithFields(l.Fields{
			"blackCoinOnBoard":    CoinsOnBoard.Black,
			"coinsCountRequested": coinsCount,
			"redCoinOnBoard":      CoinsOnBoard.Red,
			"isRedPockted":        isRedPockted,
		}).Errorln("invalid defunct request. Ignoring request")

		return fmt.Errorf(invalid)
	}

	p.Points -= 2
	p.foul()

	if isRedPockted {
		removeCoin(red, 1)
		removeCoin(black, coinsCount-1)

		return nil
	}

	removeCoin(black, coinsCount)

	return nil
}

// NoPocket removes a point when player does not pocket a coin for 3 successive turns.
func (p *Player) NoPocket() {
	p.NoPocketCount++

	if p.NoPocketCount >= 3 {
		p.Points--
		p.foul()
		p.NoPocketCount = 0
	}
}

func removeCoin(coinColor string, removalCount int) {
	switch coinColor {
	case black:
		if CoinsOnBoard.Black > removalCount {
			CoinsOnBoard.Black -= removalCount

			return
		}
		CoinsOnBoard.Black = 0
	case red:
		CoinsOnBoard.Red = 0

	default:
		l.Errorln("invalid color: ", coinColor)
	}
}

// ​foul is a turn where a player loses, at least, 1 point.
// player loses a point on three fouls.
func (p *Player) foul() {
	p.FoulCount++

	if p.FoulCount >= 3 {
		p.Points--
	}
}

// AddPlayersToGame returns true if provided player names are valid.
// Unique player names and more than one player is considered as valid.
func AddPlayersToGame(playerNames []string) bool {
	if !isValidPlayers(playerNames) {
		return false
	}

	// remove players of last game if any.
	playersOnBoard = []*Player{}

	for _, name := range playerNames {
		playersOnBoard = append(playersOnBoard, newPlayer(name))
	}

	return true
}

func isValidPlayers(playerNames []string) bool {
	if len(playerNames) < 2 {
		return false
	}

	players := make(map[string]struct{}, 0)

	for _, name := range playerNames {
		if _, ok := players[name]; ok {
			return false
		}

		players[name] = struct{}{}
	}

	return true
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
	return CoinsOnBoard.Red == 0 && CoinsOnBoard.Black == 0
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

		switch c.StrikeCode {
		case 0:
			p.Strike()
		case 1:
			err = p.MultiStrike(c.CoinsCount, c.IsRedPockted)
		case 2:
			err = p.RedStrike()
		case 3:
			p.StrikerStrike()
		case 4:
			err = p.Defunct(c.CoinsCount, c.IsRedPockted)
		case 5:
			p.NoPocket()
		default:
			l.WithField("strikeCode", c.StrikeCode).Errorln("invalid strike code")
		}
	}
}

// passPlayer internally rotates player.
func passPlayer() (p *Player) {
	if playerIDForTurn == len((playersOnBoard)) {
		playerIDForTurn = 0
	}

	p = playersOnBoard[playerIDForTurn]
	playerIDForTurn++

	return p
}
