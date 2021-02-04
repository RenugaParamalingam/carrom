package carrom

import (
	"fmt"

	l "github.com/sirupsen/logrus"
)

// Strike adds a point to player removes the pocketed coin out of game.
func (p *Player) Strike(coinsPocketed CoinsPocketedCount) error {
	if coinsPocketed.Black < 1 && coinsPocketed.White < 1 {
		l.WithField("coinsPocketedCount", coinsPocketed).Errorln("invalid request. Ignoring request")

		return fmt.Errorf(invalid)
	}

	p.Points++

	removeCoin(black, coinsPocketed.Black)
	removeCoin(white, coinsPocketed.White)

	return nil
}

// MultiStrike takes count of coins pocketed and a flag to determine red coin is pocketed.
// And error is returned incase of invalid coins count or red pocketed flag.
func (p *Player) MultiStrike(coinsPocketed CoinsPocketedCount) error {
	if coinsPocketed.Black > CoinsOnBoard.Black || coinsPocketed.White > CoinsOnBoard.White ||
		(coinsPocketed.Black == 0 && coinsPocketed.White == 0) ||
		(coinsPocketed.IsRedPocketed && CoinsOnBoard.Red == 0) {
		l.WithFields(l.Fields{
			"blackCoinOnBoard":    *CoinsOnBoard,
			"coinsCountRequested": coinsPocketed,
		}).Errorln("invalid multistrike request. Ignoring request")

		return fmt.Errorf(invalid)

	}

	p.Points += 2

	if coinsPocketed.IsRedPocketed {
		removeCoin(red, 1)
	}

	removeCoin(black, coinsPocketed.Black)
	removeCoin(white, coinsPocketed.White)

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
func (p *Player) Defunct(coinsPocketed CoinsPocketedCount) error {
	if (coinsPocketed.Black == 0 && coinsPocketed.White == 0 && !coinsPocketed.IsRedPocketed) ||
		coinsPocketed.Black > CoinsOnBoard.Black || coinsPocketed.White > CoinsOnBoard.White ||
		(coinsPocketed.IsRedPocketed && CoinsOnBoard.Red == 0) {
		l.WithFields(l.Fields{
			"blackCoinOnBoard":    *CoinsOnBoard,
			"coinsCountRequested": coinsPocketed,
		}).Errorln("invalid defunct request. Ignoring request")

		return fmt.Errorf(invalid)
	}

	p.Points -= 2
	p.foul()

	if coinsPocketed.IsRedPocketed {
		removeCoin(red, 1)
	}

	removeCoin(black, coinsPocketed.Black)
	removeCoin(white, coinsPocketed.White)

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

// ​foul is a turn where a player loses, at least, 1 point.
// player loses a point on three fouls.
func (p *Player) foul() {
	p.FoulCount++

	if p.FoulCount >= 3 {
		p.Points--
		p.FoulCount = 0
	}
}
