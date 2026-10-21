package connector

import (
	"fmt"

	"github.com/joshmedeski/sesh/v2/model"
)

func sshStrategy(c *RealConnector, name string) (model.Connection, error) {
	for _, config := range c.config.SSHConfigs {
		if config.Name != name {
			continue
		}
		return model.Connection{Found: true, Session: config.Session()}, nil
	}
	return model.Connection{}, nil
}

func connectToSSH(c *RealConnector, connection model.Connection, _ model.ConnectOpts) (string, error) {
	// Run in the current terminal so SSH can read input and manage its TTY.
	// Pass the destination as one argument, never through a local shell.
	if _, err := c.shell.CmdWithOutput("ssh", "--", connection.Session.SSHHost); err != nil {
		return "", fmt.Errorf("failed to connect to SSH host %q: %w", connection.Session.SSHHost, err)
	}
	return "", nil
}
