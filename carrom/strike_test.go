package carrom_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/RenugaParamalingam/carrom/carrom"
)

func TestAddPlayersToGame(t *testing.T) {
	testCases := []struct {
		playerNames    []string
		expectedResult bool
	}{
		{[]string{"p1", "p2"}, true},
		{[]string{"p1", "p1"}, false},
		{[]string{}, false},
	}

	for _, tc := range testCases {
		actualResult := carrom.AddPlayersToGame(tc.playerNames)

		if tc.expectedResult != actualResult {
			t.Errorf("AddPlayersToGame(%s)= %t , want= %t", tc.playerNames, actualResult, tc.expectedResult)
		}
	}
}

var invalidErr = fmt.Errorf("invalid strike input")

func TestMultiStrike(t *testing.T) {
	p := new(carrom.Player)
	p.PlayerName = "p1"

	carrom.CoinsOnBoard = &carrom.Coins{
		Red:   1,
		White: 9,
		Black: 9,
	}

	testCases := []struct {
		coinsPocketed      carrom.CoinsPocketedCount
		expectedErr        error
		expectedPoints     int
		expectedBlackCount int
		expectedWhiteCount int
		expectedRedCount   int
	}{
		{carrom.CoinsPocketedCount{0, 0, false}, invalidErr, 0, 9, 9, 1},
		{carrom.CoinsPocketedCount{10, 0, false}, invalidErr, 0, 9, 9, 1},
		{carrom.CoinsPocketedCount{2, 0, false}, nil, 2, 7, 9, 1},
		{carrom.CoinsPocketedCount{2, 2, false}, nil, 4, 5, 7, 1},
		{carrom.CoinsPocketedCount{0, 0, true}, invalidErr, 4, 5, 7, 1},
		{carrom.CoinsPocketedCount{1, 0, true}, invalidErr, 6, 4, 7, 0},
	}

	for _, tc := range testCases {
		actualErr := p.MultiStrike(tc.coinsPocketed)
		actualCoins := carrom.CoinsOnBoard

		if tc.expectedPoints != p.Points ||
			tc.expectedBlackCount != actualCoins.Black || tc.expectedRedCount != actualCoins.Red {
			t.Errorf("MultiStrike(%v)= err: %v, score: %d, black: %d, white: %d, red: %d , want= err: %v, score: %d, black: %d, white: %d, red:%d ",
				tc.coinsPocketed,
				actualErr, p.Points, actualCoins.Black, actualCoins.White, actualCoins.Red,
				tc.expectedErr, tc.expectedPoints, tc.expectedBlackCount, tc.expectedWhiteCount, tc.expectedRedCount)
		}
	}
}

func TestRedStrike(t *testing.T) {
	p := new(carrom.Player)
	p.PlayerName = "p2"

	carrom.CoinsOnBoard = &carrom.Coins{
		Red:   1,
		Black: 9,
	}

	testCases := []struct {
		expectedErr      error
		expectedPoints   int
		expectedRedCount int
	}{
		{nil, 3, 0},
		{invalidErr, 3, 0},
	}

	for _, tc := range testCases {
		actualErr := p.RedStrike()
		actualCoins := carrom.CoinsOnBoard

		if tc.expectedPoints != p.Points || tc.expectedRedCount != actualCoins.Red {
			t.Errorf("RedStrike()= err: %v, score: %d, red: %d , want= err: %v, score: %d, red:%d ",
				actualErr, p.Points, actualCoins.Red,
				tc.expectedErr, tc.expectedPoints, tc.expectedRedCount)
		}
	}
}

func TestDefunct(t *testing.T) {
	p := new(carrom.Player)
	p.PlayerName = "p1"

	carrom.CoinsOnBoard = &carrom.Coins{
		Red:   1,
		White: 9,
		Black: 9,
	}

	testCases := []struct {
		coinsPocketed      carrom.CoinsPocketedCount
		expectedErr        error
		expectedPoints     int
		expectedBlackCount int
		expectedWhiteCount int
		expectedRedCount   int
	}{

		{carrom.CoinsPocketedCount{0, 0, false}, invalidErr, 0, 9, 9, 1},
		{carrom.CoinsPocketedCount{10, 0, false}, invalidErr, 0, 9, 9, 1},
		{carrom.CoinsPocketedCount{2, 0, false}, nil, -2, 7, 9, 1},
		{carrom.CoinsPocketedCount{0, 1, false}, nil, -4, 7, 8, 1},
		{carrom.CoinsPocketedCount{0, 0, true}, nil, -7, 7, 8, 0},
		{carrom.CoinsPocketedCount{0, 0, true}, invalidErr, -7, 7, 8, 0},
	}

	for _, tc := range testCases {
		actualErr := p.Defunct(tc.coinsPocketed)
		actualCoins := carrom.CoinsOnBoard

		if tc.expectedPoints != p.Points ||
			tc.expectedBlackCount != actualCoins.Black || tc.expectedRedCount != actualCoins.Red {
			t.Errorf("Defunct(%v)= err: %v, score: %d, black: %d, white: %d, red: %d , want= err: %v, score: %d, black: %d, white: %d, red:%d ",
				tc.coinsPocketed,
				actualErr, p.Points, actualCoins.Black, actualCoins.White, actualCoins.Red,
				tc.expectedErr, tc.expectedPoints, tc.expectedBlackCount, tc.expectedWhiteCount, tc.expectedRedCount)
		}
	}
}

func TestNoPocket(t *testing.T) {
	p := new(carrom.Player)
	p.PlayerName = "p3"

	carrom.CoinsOnBoard = &carrom.Coins{
		Red:   1,
		Black: 9,
	}

	testCases := []struct {
		expectedPoints int
	}{
		{0},
		{0},
		{-1},
		{-1},
		{-1},
		{-2},
	}

	for _, tc := range testCases {
		p.NoPocket()

		if tc.expectedPoints != p.Points {
			t.Errorf("NoPocket()= score: %d, want= score: %d", p.Points, tc.expectedPoints)
		}
	}
}

func TestIsGameOver(t *testing.T) {
	testCases := []struct {
		setGame        func()
		expectedResult bool
	}{
		{setGameInDraw(), true},
		{setGameInWin(), true},
		{setGameUnFinished(), false},
	}

	for _, tc := range testCases {
		tc.setGame()

		actualResult := carrom.IsGameOver()

		fmt.Println(carrom.CoinsOnBoard)
		if tc.expectedResult != actualResult {
			t.Errorf("TestIsGameOver()= %t , want= %t", actualResult, tc.expectedResult)
		}
	}
}

func setGameInDraw() func() {
	return func() {
		carrom.AddPlayersToGame([]string{"p1", "p2"})
		strikeInput := carrom.NewBoard()

		strikeInput <- carrom.Input{StrikeCode: 0}
		strikeInput <- carrom.Input{1, carrom.CoinsPocketedCount{Black: 5, IsRedPocketed: false}}
		strikeInput <- carrom.Input{StrikeCode: 2}
		strikeInput <- carrom.Input{StrikeCode: 3}
		strikeInput <- carrom.Input{4, carrom.CoinsPocketedCount{Black: 2, White: 8, IsRedPocketed: false}}
		strikeInput <- carrom.Input{StrikeCode: 5}
		strikeInput <- carrom.Input{0, carrom.CoinsPocketedCount{Black: 2, White: 4, IsRedPocketed: false}}
		time.Sleep(time.Second * 1)
	}
}

func setGameInWin() func() {
	return func() {
		carrom.AddPlayersToGame([]string{"p3", "p4"})
		strikeInput := carrom.NewBoard()

		strikeInput <- carrom.Input{1, carrom.CoinsPocketedCount{Black: 5, IsRedPocketed: false}}
		strikeInput <- carrom.Input{0, carrom.CoinsPocketedCount{Black: 2, White: 4, IsRedPocketed: false}}
		strikeInput <- carrom.Input{StrikeCode: 2}
		strikeInput <- carrom.Input{StrikeCode: 3}
		strikeInput <- carrom.Input{StrikeCode: 0}
		strikeInput <- carrom.Input{StrikeCode: 5}
		strikeInput <- carrom.Input{0, carrom.CoinsPocketedCount{Black: 2, White: 4, IsRedPocketed: false}}
	}
}

func setGameUnFinished() func() {
	return func() {
		carrom.AddPlayersToGame([]string{"p3", "p4"})
		strikeInput := carrom.NewBoard()

		strikeInput <- carrom.Input{1, carrom.CoinsPocketedCount{Black: 5, IsRedPocketed: false}}
		strikeInput <- carrom.Input{StrikeCode: 0}
		strikeInput <- carrom.Input{StrikeCode: 2}
	}
}

func TestStrike(t *testing.T) {
	p := new(carrom.Player)
	p.PlayerName = "p3"

	carrom.CoinsOnBoard = &carrom.Coins{
		Red:   1,
		Black: 9,
		White: 9,
	}

	actualCoins := carrom.CoinsOnBoard

	testCases := []struct {
		coinsPocketed      carrom.CoinsPocketedCount
		expectedErr        error
		expectedBlackCount int
		expectedWhiteCount int
	}{
		{carrom.CoinsPocketedCount{0, 0, false}, invalidErr, 9, 9},
		{carrom.CoinsPocketedCount{1, 0, false}, nil, 8, 9},
		{carrom.CoinsPocketedCount{0, 1, false}, nil, 8, 8},
	}

	for _, tc := range testCases {
		actualErr := p.Strike(tc.coinsPocketed)

		if tc.expectedBlackCount != actualCoins.Black || tc.expectedWhiteCount != actualCoins.White {
			t.Errorf("Strike()= %v, white: %d, black: %d , want= err: %v, white: %d, black:%d ",
				actualErr, actualCoins.White, actualCoins.Black, tc.expectedErr, tc.expectedWhiteCount, tc.expectedBlackCount)
		}
	}
}
