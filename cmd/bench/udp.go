package main

import (
	"fmt"
	"net"

	"github.com/spf13/cobra"
)

const (
	defaultUDPport    = 2000
	defaultBufferSize = 1024
)

func udpCommand() *cobra.Command {
	output := &cobra.Command{
		Use:   "udp",
		Short: "udp",
		Long:  "udp operations",
	}

	output.AddCommand(
		udpBindCommand(),
	)

	return output
}

func udpBindCommand() *cobra.Command {
	var port uint32

	output := &cobra.Command{
		Use:   "bind",
		Short: "bind",
		Long:  "bind port",
		RunE: func(cmd *cobra.Command, args []string) error {
			udpServer, err := net.ListenPacket("udp", fmt.Sprintf(":%d", port))
			if err != nil {
				return err
			}
			defer func() {
				_ = udpServer.Close()
			}()

			cmd.Printf("Listening on port %d\n", port)

			go func() {
				for {
					buf := make([]byte, defaultBufferSize)
					length, _, err := udpServer.ReadFrom(buf)
					if err != nil {
						continue
					}

					cmd.Print(string(buf[:length]))
				}
			}()

			<-cmd.Context().Done()
			cmd.Println("Quitting")

			return nil
		},
	}

	output.PersistentFlags().Uint32VarP(&port, "port", "p", defaultUDPport, "port to bind")

	return output
}
