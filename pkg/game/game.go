package game

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Game struct {
	size  int
	alive map[position]cell
	rule  composeRule
}

func _new(size int, rule composeRule, initialCells []cell) (*Game, error) {
	if size <= 0 {
		return nil, fmt.Errorf("malformed game size: %d", size)
	}
	if len(initialCells) > size*size {
		return nil, fmt.Errorf(
			"number of initial positions %d is more than number of empty initialCells: %d",
			len(initialCells),
			size*size,
		)
	}

	alive := make(map[position]cell, len(initialCells))

	for _, initialCell := range initialCells {
		if initialCell.position.y > size-1 || initialCell.position.x > size-1 {
			return nil, fmt.Errorf("cell position %v is out of the game box", initialCell.position)
		}

		alive[initialCell.position] = initialCell
	}

	return &Game{alive: alive, size: size, rule: rule}, nil
}

// New initializes game
func New(size int, initialCells []cell) (*Game, error) {
	return _new(size, defaultRule, initialCells)
}

// Iterate over game cells starting with the first row
func (game *Game) gameIterator() func() (cell, error) {
	x, y := 0, 0
	return func() (cell, error) {
		if x == game.size {
			x = 0
			y++
			if y == game.size {
				return cell{}, stopIteration{}
			}

		}
		c := game.getCellByPosition(position{x, y})
		x++
		return c, nil
	}
}

func (game *Game) aliveCells() []cell {
	res := make([]cell, len(game.alive))
	for _, c := range game.alive {
		res = append(res, c)
	}
	return res
}

func (game *Game) getCellByPosition(pos position) cell {
	c, ok := game.alive[pos]
	if !ok {
		c = cell{cellType: DEAD, position: pos}
	}
	return c
}

func (game *Game) getCellsToCheck(changedCells []cell) []cell {
	seen := make(map[position]struct{}, len(changedCells))
	res := make([]cell, 0)
	for _, c := range changedCells {
		for _, neighbourCell := range game.getNeighbors(c) {
			if _, ok := seen[neighbourCell.position]; !ok {
				seen[neighbourCell.position] = struct{}{}
				res = append(res, neighbourCell)
			}
		}
		if _, ok := seen[c.position]; !ok {
			seen[c.position] = struct{}{}
			res = append(res, c)
		}
	}
	return res
}

func (game *Game) getNeighbors(c cell) []cell {
	var neighbours []cell
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if j == 0 && i == 0 {
				continue
			}
			posX := (c.position.x + j + game.size) % game.size
			posY := (c.position.y + i + game.size) % game.size
			pos := position{posX, posY}
			neighbours = append(neighbours, game.getCellByPosition(pos))
		}
	}
	return neighbours
}

func (game *Game) getCellNextState(c cell) cellType {
	return game.rule.getCellType(c, game)
}

// SetRule sets game rule
func (game *Game) SetRule(rule composeRule) {
	game.rule = rule
}

func (game *Game) getCellsToChange(cellsToCheck []cell) []cell {
	var changed []cell
	for _, c := range cellsToCheck {
		newCellType := game.getCellNextState(c)
		if newCellType != c.cellType {
			c.cellType = newCellType
			changed = append(changed, c)
		}
	}
	return changed
}

func GetGliderPattern(gameSize int) []cell {
	startPosition := position{gameSize / 2, gameSize / 2}
	return []cell{
		{
			position: position{startPosition.x, startPosition.y - 1},
			cellType: ALIVE,
		},
		{
			position: position{startPosition.x + 1, startPosition.y},
			cellType: ALIVE,
		},
		{
			position: position{startPosition.x + 1, startPosition.y + 1},
			cellType: ALIVE,
		},
		{
			position: position{startPosition.x, startPosition.y + 1},
			cellType: ALIVE,
		},
		{
			position: position{startPosition.x - 1, startPosition.y + 1},
			cellType: ALIVE,
		},
	}
}

func (game *Game) Start(ctx context.Context, wg *sync.WaitGroup, sleepTimeOutMillis int) (<-chan map[position]cell, error) {

	out := make(chan map[position]cell)

	go func() {

		timer := time.NewTicker(time.Duration(sleepTimeOutMillis) * time.Millisecond)

		defer func() {
			close(out)
			timer.Stop()
			wg.Done()
		}()

		cellsToCheck := game.getCellsToCheck(game.aliveCells())

		for {
			select {
			case <-ctx.Done():
				return
			case <-timer.C:
				changed := game.getCellsToChange(cellsToCheck)
				if len(changed) == 0 {
					return
				}
				for _, c := range changed {
					if c.cellType == ALIVE {
						game.alive[c.position] = c
					} else {
						delete(game.alive, c.position)
					}
				}

				cellsToCheck = game.getCellsToCheck(changed)

				out <- game.alive
			}
		}
	}()

	return out, nil

}
