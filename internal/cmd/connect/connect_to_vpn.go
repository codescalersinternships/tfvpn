package connect

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/rs/zerolog/log"
)

func ConnectToVPN(ipAddr string) error {
	log.Info().Msg("connecting to the vpn server...")
	cmd := exec.Command("sshuttle", "-r", "root@"+ipAddr, "0.0.0.0/0")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to connect to the vpn server %w output %s stderr %s", err, cmd.Stdout, cmd.Stderr)
	}
	log.Info().Msg("connection established successfully!")
	return nil
}
