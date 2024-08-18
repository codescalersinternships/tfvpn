// Package cmd for parsing command line arguments
package cmd

import (
	command "github.com/codescalersinternships/tfvpn/internal/cmd"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login with mnemonics to a grid network",
	Run: func(cmd *cobra.Command, args []string) {
		err := command.Login()
		if err != nil {
			log.Fatal().Err(err).Send()
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
