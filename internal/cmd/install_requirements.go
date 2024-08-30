package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/rs/zerolog/log"
)

func installRequirements() error {
	log.Info().Msg("checking if sshuttle, and python3 are installed on the system")
	cmd := exec.Command("sudo", "bash", "-c", "./install_requirements.sh")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("installing requirements failed %w output %s stderr %s", err, cmd.Stdout, cmd.Stderr)
	}

	log.Info().Msg("all requirements are installed/present successfully!")
	return nil
}
