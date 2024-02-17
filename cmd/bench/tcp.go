package main

import (
	"fmt"
	"net"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/landru29/adsb1090/internal/model"
	"github.com/landru29/adsb1090/internal/serialize/nmea"
	"github.com/spf13/cobra"
)

const (
	tcpDefaultPort = 2000
	tcpBufferSize  = 2048
	tcpMaxRetries  = 300
	tcpDefaultAddr = "127.0.0.1:3000"
)

func tcpCommand() *cobra.Command {
	output := &cobra.Command{
		Use:   "tcp",
		Short: "tcp",
		Long:  "tcp operations",
	}

	output.AddCommand(
		tcpBindCommand(),
		tcpDialCommand(),
	)

	return output
}

func tcpBindCommand() *cobra.Command {
	var port uint32

	output := &cobra.Command{
		Use:   "bind",
		Short: "bind",
		Long:  "bind tcp addr",
		RunE: func(cmd *cobra.Command, args []string) error {
			tcpServer, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
			if err != nil {
				return err
			}
			defer func() {
				_ = tcpServer.Close()
			}()

			cmd.Printf("Listening on port %d\n", port)

			groundSpeed := 120.0
			track := 110.0

			aircraft := model.Aircraft{
				Addr: 0x391217, //nolint: gomnd
				Position: &model.Position{
					Longitude: -1.8600325,
					Latitude:  48.1157851, //nolint: gomnd
				},
				Altitude:    500, //nolint: gomnd
				GroundSpeed: &groundSpeed,
				Track:       &track,
			}

			serializer := nmea.New(nmea.VesselTypeHelicopter, 0)

			go func() {
				for {
					// Wait for a connection.
					conn, err := tcpServer.Accept()
					if err != nil {
						cmd.PrintErr(err)

						return
					}

					cmd.Println("Incoming connexion")

					defer func(connexion net.Conn) {
						_ = connexion.Close()
					}(conn)

					go func(connexion net.Conn) {
						for {
							data, err := serializer.Serialize(aircraft)
							if err != nil {
								fmt.Fprintf(cmd.OutOrStderr(), "ERROR: %s\n", err)

								continue
							}
							cmd.Printf("Sending: %s\n", string(data))
							_, _ = connexion.Write(data)
							time.Sleep(time.Second)
							aircraft.Position.Longitude += 0.02
						}
					}(conn)
				}
			}()

			<-cmd.Context().Done()
			cmd.Println("Quitting")

			return nil
		},
	}

	output.PersistentFlags().Uint32VarP(&port, "port", "p", tcpDefaultPort, "port to bind")

	return output
}

func tcpDialCommand() *cobra.Command {
	var address string

	output := &cobra.Command{
		Use:   "dial",
		Short: "dial",
		Long:  "dial tcp addr",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				dialer net.Dialer
				conn   net.Conn
			)

			cmd.Printf("trying to connect to %s ...\n", address)

			bckoff := backoff.WithMaxRetries(backoff.NewConstantBackOff(1*time.Second), tcpMaxRetries)
			err := backoff.Retry(func() error {
				var err error
				conn, err = dialer.DialContext(cmd.Context(), "tcp", address)
				if err != nil {
					return err
				}

				return nil
			}, bckoff)
			if err != nil {
				return err
			}

			defer func() {
				_ = conn.Close()
			}()

			cmd.Printf("connected to %s\n", address)

			errChan := make(chan error)

			go func() {
				var cnt int
				packet := make([]byte, tcpBufferSize)
				for {
					cnt, err = conn.Read(packet)
					if err != nil {
						cmd.PrintErr(err)

						errChan <- err

						return
					}

					cmd.Printf("%d =>%s\n", cnt, string(packet[:cnt]))
				}
			}()

			for {
				select {
				case <-cmd.Context().Done():
					cmd.Println("Quitting")

					return nil
				case err := <-errChan:
					cmd.Println("error occurred")

					return err
				}
			}
		},
	}

	output.PersistentFlags().StringVarP(&address, "addr", "a", tcpDefaultAddr, "address to dial")

	return output
}
