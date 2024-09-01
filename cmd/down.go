package cmd

import (
	command "github.com/codescalersinternships/tfvpn/internal/cmd"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "stop vpn server",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := command.Down(cmd.Context()); err != nil {
			log.Fatal().Msg("failed to stop vpn server")
		}

		return nil
	},
}
