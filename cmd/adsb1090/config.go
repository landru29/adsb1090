package main

import (
	conf "github.com/landru29/adsb1090/internal/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func configCommand(settings *conf.Config) *cobra.Command {
	output := &cobra.Command{
		Use:              "config",
		Short:            "config",
		Long:             "manage configuration",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {},
	}

	output.AddCommand(
		configPathCommand(settings),
		configDisplayCommand(settings),
	)

	return output
}

func configPathCommand(settings *conf.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "path",
		Short: "path",
		Long:  "display configuration path",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(settings.Path())
		},
	}
}

func configDisplayCommand(settings *conf.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "current",
		Short: "current",
		Long:  "display configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := yaml.Marshal(settings)
			if err != nil {
				return err
			}

			cmd.Println(string(data))

			return nil
		},
	}
}
