package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/trezorg/gameoflife/pkg/game"
	"github.com/trezorg/gameoflife/pkg/writer"
)

func main() {

	size := 25
	gliderCells := game.GetGliderPattern(size)
	g, err := game.New(size)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	wrt := writer.New(size)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	ctx, cancel := context.WithCancel(context.Background())
	out, err := g.Start(ctx, wg, gliderCells, 500)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	wg.Add(1)
	go func() {

		defer wg.Done()

		for items := range out {
			for _, item := range items {
				wrt.AddItem(item)
			}
			wrt.Draw()
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		for range stop {
			cancel()
			return
		}
	}()

	wg.Wait()

}
