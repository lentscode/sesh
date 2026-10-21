package connector

import (
	"errors"
	"testing"

	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/shell"
	"github.com/joshmedeski/sesh/v2/startup"
	"github.com/joshmedeski/sesh/v2/tmux"
	"github.com/joshmedeski/sesh/v2/zoxide"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSSHDoesNotCreateTmuxSession(t *testing.T) {
	for _, opts := range []model.ConnectOpts{{}, {Switch: true}, {Command: "nvim"}} {
		t.Run(opts.Command+map[bool]string{false: "attach", true: "switch"}[opts.Switch], func(t *testing.T) {
			config := model.Config{
				SSHConfigs:           []model.SSHConfig{{Name: "dev-server", Host: "dev'host;echo"}},
				DefaultSessionConfig: model.DefaultSessionConfig{StartupCommand: "nvim", Windows: []string{"editor"}},
			}
			mShell := shell.NewMockShell(t)
			mShell.EXPECT().CmdWithOutput("ssh", "--", "dev'host;echo").Return("", nil).Once()
			// Any tmux call, session lookup, home expansion, local startup or
			// zoxide write fails these mocks, even with tmux-specific options.
			c := NewConnector(config, nil, home.NewMockHome(t), lister.NewMockLister(t), nil,
				startup.NewMockStartup(t), tmux.NewMockTmux(t), zoxide.NewMockZoxide(t), nil, nil, mShell)
			message, err := c.Connect("dev-server", opts)
			require.NoError(t, err)
			assert.Empty(t, message)
		})
	}
}

func TestSSHConnectionFailure(t *testing.T) {
	config := model.Config{SSHConfigs: []model.SSHConfig{{Name: "dev", Host: "dev-host"}}}
	mShell := shell.NewMockShell(t)
	runErr := errors.New("SSH exited with status 255")
	mShell.EXPECT().CmdWithOutput("ssh", "--", "dev-host").Return("", runErr)
	c := NewConnector(config, nil, nil, nil, nil, nil, tmux.NewMockTmux(t), nil, nil, nil, mShell)
	_, err := c.Connect("dev", model.ConnectOpts{})
	assert.ErrorContains(t, err, `failed to connect to SSH host "dev-host"`)
	assert.ErrorIs(t, err, runErr)
}

func TestConnectWindowRejectsSSHHost(t *testing.T) {
	c := NewConnector(model.Config{SSHConfigs: []model.SSHConfig{{Name: "dev", Host: "dev-host"}}},
		nil, nil, lister.NewMockLister(t), nil, nil, tmux.NewMockTmux(t), nil, nil, nil, shell.NewMockShell(t))
	_, err := c.ConnectWindow(model.WindowConnectOpts{Session: "dev"})
	assert.ErrorContains(t, err, "is an SSH host, not a tmux session")
}
