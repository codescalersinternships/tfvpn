package connect

import (
	"fmt"
	"net"
	"path/filepath"
	"time"

	"github.com/melbahja/goph"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh"
)

// ConnectToVM attempts to establish an SSH connection within the given timeout.
func ConnectToVM(timeout time.Duration, user, ipAddr string) error {
	var auth goph.Auth
	var err error

	if goph.HasAgent() {
		log.Info().Msg("using ssh agent for authentication")
		auth, err = goph.UseAgent()
	} else {
		log.Info().Msg("using private key for authentication")
		var sshDir, privKeyPath string
		sshDir, err = GetUserSSHDir()
		if err != nil {
			return err
		}
		privKeyPath = filepath.Join(sshDir, "id_rsa")
		auth, err = goph.Key(privKeyPath, "")
	}

	if err != nil {
		return fmt.Errorf("failed to get ssh auth %v", err)
	}

	startTime := time.Now()
	var client *goph.Client

	for {
		elapsedTime := time.Since(startTime)
		if elapsedTime >= timeout {
			return fmt.Errorf("timeout reached while waiting for SSH connection %v", err)
		}

		remainingTime := timeout - elapsedTime

		client, err = goph.NewConn(&goph.Config{
			User:     user,
			Port:     22,
			Addr:     ipAddr,
			Auth:     auth,
			Callback: verifyHost,
		})

		if err == nil {
			defer client.Close()
			return nil
		}

		log.Info().Msg("SSH connection attempt failed, retrying...")
		time.Sleep(time.Second)

		if remainingTime < time.Second {
			time.Sleep(remainingTime)
			break
		}
	}

	return fmt.Errorf("failed to establish connection %v", err)
}

// verifyHost verifies the host key of the server.
func verifyHost(host string, remote net.Addr, key ssh.PublicKey) error {
	// Check if the host is in known hosts file.
	hostFound, err := goph.CheckKnownHost(host, remote, key, "")

	// Host in known hosts but key mismatch!
	// Maybe because of MAN IN THE MIDDLE ATTACK!
	if hostFound && err != nil {

		return err
	}

	// handshake because public key already exists.
	if hostFound && err == nil {
		return nil
	}

	// Add the new host to known hosts file.
	return goph.AddKnownHost(host, remote, key, "")
}
