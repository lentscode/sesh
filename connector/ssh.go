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
		// Reuse an existing session, and never record its local working
		// directory as a zoxide visit just because we connected to a host.
		if session, exists := c.lister.FindTmuxSession(name); exists {
			return model.Connection{Found: true, Session: session}, nil
		}
		session := config.Session()
		path, err := c.home.ExpandPath("~")
		if err != nil {
			return model.Connection{}, fmt.Errorf("couldn't expand home for SSH session: %w", err)
		}
		session.Path = path
		return model.Connection{Found: true, New: true, Session: session}, nil
	}
	return model.Connection{}, nil
}
