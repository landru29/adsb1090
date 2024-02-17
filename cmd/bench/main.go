// Package main is the main application.
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

func main() {
	rootCommand := &cobra.Command{
		Use:   "bench",
		Short: "bench",
		Long:  "bench to test communication with network",
	}

	rootCommand.AddCommand(
		udpCommand(),
		tcpCommand(),
	)

	osSignal := make(chan os.Signal, 1)

	// add any other syscalls that you want to be notified with
	signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()

	go func() {
		<-osSignal

		cancel()
	}()

	if err := rootCommand.ExecuteContext(ctx); err != nil {
		fmt.Println(err) //nolint: forbidigo

		cancel()

		os.Exit(1) //nolint: gocritic
	}
}
