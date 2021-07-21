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
	game, err := New(size, nil)
	require.NoError(t, err)
	require.Equal(t, size, game.size)
}

func TestGameNeighbors(t *testing.T) {
	size := 10
	game, err := New(size, initCells)
	require.NoError(t, err)
	require.Equal(t, size, game.size)

	// cell 0, 0
	neighbours := game.getNeighbors(initCells[0])
	require.Equal(t, 8, len(neighbours))
	require.Equal(t, neighbours[len(neighbours)-1].position, position{1, 1})
	require.Equal(t, neighbours[0].position, position{size - 1, size - 1})

}

func TestGameStartTwoAliveCells(t *testing.T) {
	size := 10
	game, err := New(size, initCells)
	require.NoError(t, err)
	require.Equal(t, size, game.size)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	ctx, cancel := context.WithCancel(context.Background())
	out, err := game.Start(ctx, wg, 200)
	require.NoError(t, err)

	go func() {
		for range out {
		}
	}()
	wg.Wait()
	cancel()

	var alive []cell
	it := game.gameIterator()
	for c, err := it(); err == nil; c, err = it() {
		if c.cellType == ALIVE {
			alive = append(alive, c)
		}
	}
	require.Len(t, alive, 0)
}

func TestGameGliderPattern(t *testing.T) {

	size := 25
	gliderCells := GetGliderPattern(size)
	game, err := New(size, gliderCells)
	require.NoError(t, err)
	require.Equal(t, size, game.size)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	ctx, cancel := context.WithCancel(context.Background())
	out, err := game.Start(ctx, wg, 200)
	require.NoError(t, err)
	go func() {
		for range out {
		}
	}()
	go func() {
		tick := time.NewTicker(time.Duration(1) * time.Second)

		defer tick.Stop()

		for range tick.C {
			cancel()
			return
		}
	}()
	wg.Wait()

	alive := 0
	it := game.gameIterator()
	for c, err := it(); err == nil; c, err = it() {
		if c.cellType == ALIVE {
			alive++
		}
	}
	require.Equal(t, len(gliderCells), alive)

}
