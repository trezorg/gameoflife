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
// size: grid size
func New(size int, initialCells ...cell) (*Game, error) {
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
	for _, initialCell := range initialCells {
		if initialCell.position.y > size-1 || initialCell.position.x > size-1 {
			return nil, fmt.Errorf("cell position %v is out of the game box", initialCell.position)
		}
		gridCells[initialCell.position.y][initialCell.position.x] = initialCell
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

func (game *Game) countByCellType() (int, int) {
	alive := 0
	it := game.gameIterator()
	for c, err := it(); err == nil; c, err = it() {
		if c.cellType == ALIVE {
			alive++
		}
	}
	return alive, game.size*game.size - alive
}

func (game *Game) alive() []cell {
	var alive []cell
	it := game.gameIterator()
	for c, err := it(); err == nil; c, err = it() {
		if c.cellType == ALIVE {
			alive = append(alive, c)
		}
	}
	return alive
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

func (game *Game) getItemsToChange() []cell {
	var changed []cell
	it := game.gameIterator()
	for c, err := it(); err == nil; c, err = it() {
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

func (game *Game) Start(ctx context.Context, wg *sync.WaitGroup, sleepTimeOut int) <-chan []cell {

	timer := time.NewTicker(time.Duration(sleepTimeOut) * time.Second)
	out := make(chan []cell)

	go func() {

		defer func() {
			close(out)
			timer.Stop()
			wg.Done()
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case <-timer.C:
				changed := game.getItemsToChange()
				if len(changed) == 0 {
					return
				}
				for _, c := range changed {
					game.cells[c.position.y][c.position.x] = c
				}

				out <- game.alive()
			}
		}
	}()

	return out

}
