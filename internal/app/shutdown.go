package app

import (
	"context"
	"os/signal"
	"syscall"
)

func waitForShutdown() context.Context {
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	go func() {
		<-ctx.Done()
		stop()
	}()

	return ctx
}
