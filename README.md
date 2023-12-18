# carrom

A game in carrom-board called ​Clean Strike.

The game is described as follows:
● There are 9 black coins, a red coin and a striker on the carrom-board
● Strike​ - When a player pockets a coin he/she wins a point
● Multi-strike - When a player pockets more than one coin he/she wins 2 points. All, but 2
coins, that were pocketed, get back on to the carrom-board
● Red strike - When a player pockets red coin he/she wins 3 points. If other coins are
pocketed along with red coin in the same turn, other coins get back on to the
carrom-board
● Striker strike​ - When a player pockets the striker he/she loses a point
● Defunct coin - When a coin is thrown out of the carrom-board, due to a strike, the player
loses 2 points, and the coin goes out of play
● When a player does not pocket a coin for 3 successive turns he/she loses a point
● When a player ​fouls 3 times (a ​foul is a turn where a player loses, at least, 1 point),
he/she loses an additional point
● A ​game is won by the first player to have won at least 5 points, in total, and, at least, 3
points more than the opponent
● When the coins are exhausted on the board, if the highest scorer is not leading by, at
least, 3 points or does not have a minimum of 5 points, the game is considered a draw

### Local build and run

**Build**

```
go build -mod vendor
```

**Run**

```
./carrom
```

### Run test case

```
cd carrom
go test
```

A new game in carrom-board called ​Clean Strike is played by 2 players with multiple ​turn​s. A turn has a player attempting to strike a coin with the striker. Players alternate in taking turns.
