package main

import (
	"github.com/landru29/adsb1090/internal/serialize"
	"github.com/spf13/cobra"
)

func serializerCommand(availableSerializers *[]serialize.Serializer) *cobra.Command {
	return &cobra.Command{
		Use: "serializers",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("Available serializers:")

			for _, serializer := range *availableSerializers {
				cmd.Printf(" - %s\n", serializer)
			}
		},
	}
}
