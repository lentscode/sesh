package lister

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/formatter"
	"github.com/joshmedeski/sesh/v2/icon"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/tmux"
	"github.com/joshmedeski/sesh/v2/tmuxinator"
	"github.com/joshmedeski/sesh/v2/zoxide"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListSSH(t *testing.T) {
	config := model.Config{SSHConfigs: []model.SSHConfig{
		{Name: "dev-server", Host: "dev"},
		{Name: "backup", Host: "alice@backup.example.com"},
	}}
	for _, opts := range []ListOptions{{}, {SSH: true}, {SSH: true, Icons: true, NoColor: true}, {SSH: true, Json: true}} {
		t.Run("list", func(t *testing.T) {
			mTmux := tmux.NewMockTmux(t)
			mZoxide := zoxide.NewMockZoxide(t)
			mTmuxinator := tmuxinator.NewMockTmuxinator(t)
			if !opts.SSH {
				mTmux.EXPECT().ListSessions().Return(nil, nil)
				mZoxide.EXPECT().ListResults().Return(nil, nil)
				mTmuxinator.EXPECT().List().Return(nil, nil)
			}
			ic := icon.NewIcon(config)
			l := NewLister(config, nil, mTmux, mZoxide, mTmuxinator, formatter.NewFormatter(ic), nil)
			sessions, err := l.List(opts)
			require.NoError(t, err)
			assert.Equal(t, []string{"ssh:dev-server", "ssh:backup"}, sessions.OrderedIndex)
			dev := sessions.Directory["ssh:dev-server"]
			assert.Equal(t, "ssh", dev.Src)
			assert.Equal(t, "dev", dev.SSHHost)
			assert.Empty(t, dev.Path)
			assert.Equal(t, "dev-server", ic.RemoveIcon(dev.Name))
			if opts.Icons {
				assert.NotEqual(t, "dev-server", dev.Name)
			}
		})
	}
}

func TestSSHHideDuplicates(t *testing.T) {
	mTmux := tmux.NewMockTmux(t)
	mTmux.EXPECT().ListSessions().Return([]*model.TmuxSession{{Name: "dev-server", Path: "/home/user"}}, nil)
	l := NewLister(model.Config{SSHConfigs: []model.SSHConfig{{Name: "dev-server", Host: "dev"}}}, nil, mTmux, nil, nil, nil, nil)
	sessions, err := l.List(ListOptions{SSH: true, Tmux: true, HideDuplicates: true})
	require.NoError(t, err)
	require.Len(t, sessions.OrderedIndex, 1)
	assert.Equal(t, "tmux", sessions.Directory[sessions.OrderedIndex[0]].Src)
}
