package connector

import (
	"errors"
	"testing"

	"github.com/joshmedeski/sesh/v2/home"
	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/startup"
	"github.com/joshmedeski/sesh/v2/tmux"
	"github.com/joshmedeski/sesh/v2/zoxide"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnectSSH(t *testing.T) {
	for _, existing := range []bool{false, true} {
		t.Run(map[bool]string{false: "new", true: "reuse"}[existing], func(t *testing.T) {
			config := model.Config{
				SSHConfigs:           []model.SSHConfig{{Name: "dev-server", Host: "dev'host;echo"}},
				DefaultSessionConfig: model.DefaultSessionConfig{StartupCommand: "nvim", Windows: []string{"editor"}},
			}
			mLister := lister.NewMockLister(t)
			mHome := home.NewMockHome(t)
			mTmux := tmux.NewMockTmux(t)
			mTmux.EXPECT().IsAttached().Return(false)
			mZoxide := zoxide.NewMockZoxide(t)
			var session model.SeshSession
			if existing {
				session = model.SeshSession{Src: "tmux", Name: "dev-server", Path: "/home/test"}
			} else {
				mHome.EXPECT().ExpandPath("~").Return("/home/test", nil)
				mTmux.EXPECT().NewSession("dev-server", "/home/test").Return("", nil)
				mTmux.EXPECT().SendKeys("dev-server:", `ssh -- 'dev'\''host;echo'`).Return("", nil).Once()
			}
			mLister.EXPECT().FindTmuxSession("dev-server").Return(session, existing).Once()
			mTmux.EXPECT().SwitchOrAttach("dev-server", model.ConnectOpts{}).Return("connected", nil)
			s := startup.NewStartup(config, mLister, mTmux, mHome, nil)
			c := NewConnector(config, nil, mHome, mLister, nil, s, mTmux, mZoxide, nil, nil)
			message, err := c.Connect("dev-server", model.ConnectOpts{})
			require.NoError(t, err)
			assert.Equal(t, "connected", message)
			// Unexpected calls fail the mocks: no zoxide write or local startup.
		})
	}
}

func TestSSHStartupFailure(t *testing.T) {
	config := model.Config{SSHConfigs: []model.SSHConfig{{Name: "dev", Host: "dev-host"}}}
	mLister := lister.NewMockLister(t)
	mHome := home.NewMockHome(t)
	mTmux := tmux.NewMockTmux(t)
	mTmux.EXPECT().IsAttached().Return(false)
	mLister.EXPECT().FindTmuxSession("dev").Return(model.SeshSession{}, false)
	mHome.EXPECT().ExpandPath("~").Return("/home/test", nil)
	mTmux.EXPECT().NewSession("dev", "/home/test").Return("", nil)
	mTmux.EXPECT().SendKeys("dev:", "ssh -- 'dev-host'").Return("", errors.New("send failed"))
	s := startup.NewStartup(config, mLister, mTmux, mHome, nil)
	c := NewConnector(config, nil, mHome, mLister, nil, s, mTmux, nil, nil, nil)
	_, err := c.Connect("dev", model.ConnectOpts{})
	assert.ErrorContains(t, err, "failed to start SSH session")
}
