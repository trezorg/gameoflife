package game

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	initCells = []cell{
		{
			position: position{0, 0},
			cellType: ALIVE,
		},
		{
			position: position{1, 0},
			cellType: ALIVE,
		},
	}
)

func TestGameInit(t *testing.T) {

	size := 10
	game, err := New(size, initCells...)
	require.NoError(t, err)
	require.Equal(t, size, game.size)
	alive, dead := game.countByCellType()

	require.Equal(t, len(initCells), alive)
	require.Equal(t, size*size-alive, dead)

}

func TestGameNeighbors(t *testing.T) {
	size := 10
	game, err := New(size, initCells...)
	require.NoError(t, err)
	require.Equal(t, size, game.size)

	// cell 0, 0
	neighbours := game.getNeighbors(initCells[0])
	require.Equal(t, 8, len(neighbours))
	require.Equal(t, neighbours[len(neighbours)-1].position, position{1, 1})
	require.Equal(t, neighbours[0].position, position{size - 1, size - 1})

}

func TestGameCellNextState(t *testing.T) {
	size := 10
	game, err := New(size, initCells...)
	require.NoError(t, err)
	require.Equal(t, size, game.size)

	c := initCells[0]
	state := game.getCellNextState(c)
	require.Equal(t, DEAD, state)

	c = initCells[1]
	state = game.getCellNextState(c)
	require.Equal(t, DEAD, state)

}

func TestGameWithTwoAliveCellsFinished(t *testing.T) {
	size := 10
	game, err := New(size, initCells...)
	require.NoError(t, err)
	require.Equal(t, size, game.size)

	changed := game.getItemsToChange()
	require.Len(t, changed, 2)

}

func TestGameStartTwoAliveCells(t *testing.T) {
	size := 10
	game, err := New(size, initCells...)
	require.NoError(t, err)
	require.Equal(t, size, game.size)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for range game.Start(ctx, wg, 1) {
		}
	}()
	wg.Wait()
	cancel()
	changed := game.getItemsToChange()
	require.Len(t, changed, 0)

}

func TestGameGliderPattern(t *testing.T) {

	size := 25
	gliderCells := GetGliderPattern(size)
	game, err := New(size, gliderCells...)
	require.NoError(t, err)
	require.Equal(t, size, game.size)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for range game.Start(ctx, wg, 1) {
		}
	}()
	go func() {
		tick := time.NewTicker(time.Duration(5) * time.Second)
		for range tick.C {
			cancel()
			return
		}
	}()
	wg.Wait()

	alive, _ := game.countByCellType()
	require.Equal(t, len(gliderCells), alive)

}
