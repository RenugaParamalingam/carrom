package carrom

var playersOnBoard []*Player
var playerIDForTurn int

// Player is details of the player in game.
type Player struct {
	PlayerName    string
	Points        int
	FoulCount     int
	NoPocketCount int
}

func newPlayer(name string) *Player {
	return &Player{
		PlayerName: name,
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

// passPlayer internally rotates player.
func passPlayer() (p *Player) {
	if playerIDForTurn == len((playersOnBoard)) {
		playerIDForTurn = 0
	}

	p = playersOnBoard[playerIDForTurn]
	playerIDForTurn++

	return p
}
