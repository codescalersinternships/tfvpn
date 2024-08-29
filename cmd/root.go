package cmd

import (
	"context"
	"os"

	"github.com/codescalersinternships/tfvpn/pkg/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Send()
	}
}

var rootCmd = &cobra.Command{
	Use:   "grid-compose",
	Short: "Grid-Compose is a tool for running multi-vm applications on TFGrid defined using a Yaml formatted file.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		network := os.Getenv("NETWORK")
		mnemonics := os.Getenv("MNEMONIC")
		config := config.Config{
			Mnemonics: mnemonics,
			Network:   network,
		}
		if err := config.Validate(); err != nil {
			log.Fatal().Err(err).Msg("failed to read config")
		}

		client, err := setup(config)
		if err != nil {
			log.Fatal().Err(err).Send()
		}

		//nolint:staticcheck
		ctx := context.WithValue(cmd.Context(), "client", client)
		cmd.SetContext(ctx)
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}
