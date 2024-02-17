package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/fatih/color"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	s := make(chan os.Signal, 1)

	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		<-s
		cancel()

		os.Exit(0)
	}()

	if err := rootCommand().ExecuteContext(ctx); err != nil {
		colored := color.New(color.FgRed).SprintFunc()

		fmt.Fprintf(os.Stderr, colored("[ERROR] ", err.Error(), "\n")) //nolint: staticcheck

		os.Exit(1)
	}
}
