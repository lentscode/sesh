package previewer

import (
	"fmt"

	"github.com/joshmedeski/sesh/v2/model"
)

type SSHPreviewStrategy struct {
	config model.Config
}

func NewSSHStrategy(config model.Config) *SSHPreviewStrategy {
	return &SSHPreviewStrategy{config: config}
}

func (s *SSHPreviewStrategy) Execute(name string) (string, error) {
	for _, config := range s.config.SSHConfigs {
		if config.Name == name {
			return fmt.Sprintf("SSH: %s\nHost: %s\n\n%s\n", config.Name, config.Host, config.Session().StartupCommand), nil
		}
	}
	return "", nil
}
