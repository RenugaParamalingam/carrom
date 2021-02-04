package carrom

import l "github.com/sirupsen/logrus"

// Coins holds different coins on board.
type Coins struct {
	Red   int
	Black int
	White int
}

// CoinsPocketedCount is count of coins pocketed
type CoinsPocketedCount struct {
	Black         int
	White         int
	IsRedPocketed bool
}

var (
	red   = "red"
	black = "black"
	white = "white"
)

// CoinsOnBoard gives coins and it's current count.
var CoinsOnBoard *Coins

func removeCoin(coinColor string, removalCount int) {
	switch coinColor {
	case black:
		if CoinsOnBoard.Black > removalCount {
			CoinsOnBoard.Black -= removalCount

			return
		}
		CoinsOnBoard.Black = 0
	case white:
		if CoinsOnBoard.White > removalCount {
			CoinsOnBoard.White -= removalCount

			return
		}
		CoinsOnBoard.White = 0
	case red:
		CoinsOnBoard.Red = 0

	default:
		l.Errorln("invalid color: ", coinColor)
	}
}
