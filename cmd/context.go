package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func WithContext(ctx context.Context) context.Context {
	return WithContextFunc(
		ctx, func() {
			fmt.Println("Interrupt received, terminating process.")
		},
	)
}

func WithContextFunc(ctx context.Context, f func()) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		c := make(chan os.Signal, 1)

		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		defer signal.Stop(c)

		select {
		case <-ctx.Done():
		case <-c:
			f()
			cancel()
		}
	}()

	return ctx
}
