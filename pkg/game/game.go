package game

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Game struct {
	cells [][]cell
	size  int
}

// New initializes game
// size: int. grid size
func New(size int) (*Game, error) {
	if size <= 0 {
		return nil, fmt.Errorf("malformed game size: %d", size)
	}
	gridCells := make([][]cell, size)
	for y := 0; y < size; y++ {
		row := make([]cell, size)
		for x := 0; x < size; x++ {
			row[x] = cell{
				cellType: DEAD,
				position: position{x: x, y: y},
			}
			gridCells[y] = row
		}
	}
	return &Game{cells: gridCells, size: size}, nil
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
		c := game.cells[y][x]
		x++
		return c, nil
	}
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
			neighbours = append(neighbours, game.cells[posY][posX])
		}
	}
	return neighbours
}

func (game *Game) getCellNextState(c cell) cellType {
	alive := 0
	for _, c := range game.getNeighbors(c) {
		if c.cellType == ALIVE {
			alive++
		}
	}
	if c.cellType == ALIVE {
		if alive < 2 {
			return DEAD
		}
		if alive == 2 || alive == 3 {
			return ALIVE
		}
		if alive > 3 {
			return DEAD
		}
	}
	if c.cellType == DEAD && alive == 3 {
		return ALIVE
	}
	return c.cellType
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

func (game *Game) Start(ctx context.Context, wg *sync.WaitGroup, initialCells []cell, sleepTimeOutMillis int) (<-chan map[position]cell, error) {

	if len(initialCells) > game.size*game.size {
		return nil, fmt.Errorf(
			"number of initial positions %d is more than number of empty initialCells: %d",
			len(initialCells),
			game.size*game.size,
		)
	}
	for _, initialCell := range initialCells {
		if initialCell.position.y > game.size-1 || initialCell.position.x > game.size-1 {
			return nil, fmt.Errorf("cell position %v is out of the game box", initialCell.position)
		}
		game.cells[initialCell.position.y][initialCell.position.x] = initialCell
	}

	out := make(chan map[position]cell)

	go func() {

		timer := time.NewTicker(time.Duration(sleepTimeOutMillis) * time.Millisecond)

		defer func() {
			close(out)
			timer.Stop()
			wg.Done()
		}()

		cellsToCheck := game.getCellsToCheck(initialCells)
		alive := make(map[position]cell, len(initialCells))

		for _, c := range initialCells {
			alive[c.position] = c
		}

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
					game.cells[c.position.y][c.position.x] = c
					if c.cellType == ALIVE {
						alive[c.position] = c
					} else {
						delete(alive, c.position)
					}
				}

				cellsToCheck = game.getCellsToCheck(changed)

				out <- alive
			}
		}
	}()

	return out, nil

}
