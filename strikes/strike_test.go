package strikes_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/RenugaParamalingam/carrom/strikes"
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
		actualResult := strikes.AddPlayersToGame(tc.playerNames)

		if tc.expectedResult != actualResult {
			t.Errorf("AddPlayersToGame(%s)= %t , want= %t", tc.playerNames, actualResult, tc.expectedResult)
		}
	}
}

var invalidErr = fmt.Errorf("invalid strike input")

func TestMultiStrike(t *testing.T) {
	p := new(strikes.Player)
	p.PlayerName = "p1"

	strikes.CoinsOnBoard = &strikes.Coins{
		Red:   1,
		Black: 9,
	}

	testCases := []struct {
		coinsCount         int
		isRedPockted       bool
		expectedErr        error
		expectedPoints     int
		expectedBlackCount int
		expectedRedCount   int
	}{
		{0, false, invalidErr, 0, 9, 1},
		{10, false, invalidErr, 0, 9, 1},
		{2, true, nil, 2, 8, 0},
		{4, true, invalidErr, 2, 8, 0},
	}

	for _, tc := range testCases {
		actualErr := p.MultiStrike(tc.coinsCount, tc.isRedPockted)
		actualCoins := strikes.CoinsOnBoard

		if tc.expectedPoints != p.Points ||
			tc.expectedBlackCount != actualCoins.Black || tc.expectedRedCount != actualCoins.Red {
			t.Errorf("MultiStrike(%d,%t)= err: %v, score: %d, black: %d, red: %d , want= err: %v, score: %d, black: %d, red:%d ",
				tc.coinsCount, tc.isRedPockted,
				actualErr, p.Points, actualCoins.Black, actualCoins.Red,
				tc.expectedErr, tc.expectedPoints, tc.expectedBlackCount, tc.expectedRedCount)
		}
	}
}

func TestRedStrike(t *testing.T) {
	p := new(strikes.Player)
	p.PlayerName = "p2"

	strikes.CoinsOnBoard = &strikes.Coins{
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
		actualCoins := strikes.CoinsOnBoard

		if tc.expectedPoints != p.Points || tc.expectedRedCount != actualCoins.Red {
			t.Errorf("RedStrike()= err: %v, score: %d, red: %d , want= err: %v, score: %d, red:%d ",
				actualErr, p.Points, actualCoins.Red,
				tc.expectedErr, tc.expectedPoints, tc.expectedRedCount)
		}
	}
}

func TestDefunct(t *testing.T) {
	p := new(strikes.Player)
	p.PlayerName = "p1"

	strikes.CoinsOnBoard = &strikes.Coins{
		Red:   1,
		Black: 9,
	}

	testCases := []struct {
		coinsCount         int
		isRedPockted       bool
		expectedErr        error
		expectedPoints     int
		expectedBlackCount int
		expectedRedCount   int
	}{
		{0, false, invalidErr, 0, 9, 1},
		{10, false, invalidErr, 0, 9, 1},
		{2, true, nil, -2, 8, 0},
		{4, true, invalidErr, -2, 8, 0},
	}

	for _, tc := range testCases {
		actualErr := p.Defunct(tc.coinsCount, tc.isRedPockted)
		actualCoins := strikes.CoinsOnBoard

		if tc.expectedPoints != p.Points ||
			tc.expectedBlackCount != actualCoins.Black || tc.expectedRedCount != actualCoins.Red {
			t.Errorf("Defunct(%d,%t)= err: %v, score: %d, black: %d, red: %d , want= err: %v, score: %d, black: %d, red:%d ",
				tc.coinsCount, tc.isRedPockted,
				actualErr, p.Points, actualCoins.Black, actualCoins.Red,
				tc.expectedErr, tc.expectedPoints, tc.expectedBlackCount, tc.expectedRedCount)
		}
	}
}

func TestNoPocket(t *testing.T) {
	p := new(strikes.Player)
	p.PlayerName = "p3"

	strikes.CoinsOnBoard = &strikes.Coins{
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
		actualResult := strikes.IsGameOver()

		if tc.expectedResult != actualResult {
			t.Errorf("TestIsGameOver()= %t , want= %t", actualResult, tc.expectedResult)
		}
	}
}

func setGameInDraw() func() {
	return func() {
		strikes.AddPlayersToGame([]string{"p1", "p2"})
		strikeInput := strikes.NewBoard()

		strikeInput <- strikes.Input{StrikeCode: 0}
		strikeInput <- strikes.Input{StrikeCode: 1, CoinsCount: 5, IsRedPockted: false}
		strikeInput <- strikes.Input{StrikeCode: 2}
		strikeInput <- strikes.Input{StrikeCode: 3}
		strikeInput <- strikes.Input{StrikeCode: 4, CoinsCount: 2, IsRedPockted: false}
		strikeInput <- strikes.Input{StrikeCode: 5}
		strikeInput <- strikes.Input{StrikeCode: 0}
		strikeInput <- strikes.Input{StrikeCode: 0}
		strikeInput <- strikes.Input{StrikeCode: 0}
		strikeInput <- strikes.Input{StrikeCode: 0}
		time.Sleep(time.Second * 1)
	}
}

func setGameInWin() func() {
	return func() {
		strikes.AddPlayersToGame([]string{"p3", "p4"})
		strikeInput := strikes.NewBoard()

		strikeInput <- strikes.Input{StrikeCode: 1, CoinsCount: 5, IsRedPockted: false}
		strikeInput <- strikes.Input{StrikeCode: 0}
		strikeInput <- strikes.Input{StrikeCode: 2}
		strikeInput <- strikes.Input{StrikeCode: 3}
		strikeInput <- strikes.Input{StrikeCode: 0}
		strikeInput <- strikes.Input{StrikeCode: 5}
	}
}

func setGameUnFinished() func() {
	return func() {
		strikes.AddPlayersToGame([]string{"p3", "p4"})
		strikeInput := strikes.NewBoard()

		strikeInput <- strikes.Input{StrikeCode: 1, CoinsCount: 5, IsRedPockted: false}
		strikeInput <- strikes.Input{StrikeCode: 0}
		strikeInput <- strikes.Input{StrikeCode: 2}
	}
}
