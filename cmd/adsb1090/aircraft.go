package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/guumaster/logsymbols"
	"github.com/landru29/adsb1090/internal/aircraftdb"
	"github.com/landru29/adsb1090/internal/config"
	"github.com/landru29/adsb1090/internal/model"
	"github.com/spf13/cobra"
)

func aircraftCommand(settings *config.Config) *cobra.Command {
	output := &cobra.Command{
		Use:              "aircraft",
		Short:            "aircraft",
		Long:             "manage aircraft database",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {},
	}

	output.AddCommand(
		aircraftDownloadCommand(settings),
		aircraftSearchCommand(settings),
		aircraftCleanCommand(settings),
	)

	return output
}

func aircraftDownloadCommand(settings *config.Config) *cobra.Command {
	var urlStr string

	output := &cobra.Command{
		Use:   "download",
		Short: "download",
		Long:  "download the latest aircraft database",
		RunE: func(cmd *cobra.Command, args []string) error {
			output := cmd.OutOrStdout()

			if err := download(settings.AircraftDatabaseFile(), urlStr, output); err != nil {
				return err
			}

			stat, err := os.Stat(settings.AircraftDatabaseFile())
			if err != nil {
				fmt.Fprintf(output, " %s\n", logsymbols.Error)

				return fmt.Errorf("checking file %s: %w", settings.AircraftDatabaseFile(), err)
			}

			cmd.Printf("Size: %s %s\n", humanize.Bytes(uint64(stat.Size())), logsymbols.Success)

			return nil
		},
	}

	output.Flags().StringVarP(
		&urlStr,
		"url",
		"u",
		"https://opensky-network.org/datasets/metadata/",
		"URL to aircraft database",
	)

	return output
}

func aircraftSearchCommand(settings *config.Config) *cobra.Command {
	var (
		registration string
		addr         model.ICAOAddr
	)

	output := &cobra.Command{
		Use:   "search",
		Short: "search",
		Long:  "search in the aircraft database",
		RunE: func(cmd *cobra.Command, args []string) error {
			var database aircraftdb.Database
			if err := database.Load(settings.AircraftDatabaseFile(), io.Discard); err != nil {
				return fmt.Errorf(
					"loading aircrafts: %w (please run '%s aircrafts download')", //nolint: misspell
					err,
					cmd.CommandPath(),
				)
			}

			if addr != 0 {
				cmd.Println(database[addr].String())

				return nil
			}

			if registration != "" {
				for _, elt := range database {
					if strings.EqualFold(elt.Registration, registration) {
						cmd.Println(elt.String())
					}
				}
			}

			return nil
		},
	}

	output.Flags().StringVarP(
		&registration,
		"registration",
		"r",
		"",
		"aircraft registration",
	)

	output.Flags().VarP(&addr, "addr", "", "OACI address")

	return output
}

func aircraftCleanCommand(settings *config.Config) *cobra.Command {
	var (
		registration string
		addr         model.ICAOAddr
	)

	output := &cobra.Command{
		Use:   "clean",
		Short: "clean",
		Long:  "clean the aircraft database",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := os.Stat(settings.AircraftDatabaseFile())

			switch {
			case errors.Is(err, os.ErrNotExist):
				cmd.Printf("%s [DELETED] %s\n", settings.AircraftDatabaseFile(), logsymbols.Success)
			case err != nil:
				return err
			default:
				if err := os.Remove(settings.AircraftDatabaseFile()); err != nil {
					return err
				}

				cmd.Printf("%s [DELETED] %s\n", settings.AircraftDatabaseFile(), logsymbols.Success)
			}

			return nil
		},
	}

	output.Flags().StringVarP(
		&registration,
		"registration",
		"r",
		"",
		"aircraft registration",
	)

	output.Flags().VarP(&addr, "addr", "", "OACI address")

	return output
}

func download(filename string, urlStr string, output io.Writer) error {
	database, err := aircraftdb.DownloadLatest(urlStr, output)
	if err != nil {
		return fmt.Errorf("downloading database: %w", err)
	}

	fmt.Fprintf(output, " %s\n", logsymbols.Success)

	if err := database.Save(filename, output); err != nil {
		fmt.Fprintf(output, " %s\n", logsymbols.Error)

		return fmt.Errorf("saving file to %s: %w", filename, err)
	}

	fmt.Fprintf(output, " %s\n", logsymbols.Success)

	return nil
}
