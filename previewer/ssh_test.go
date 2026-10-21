package previewer

import (
	"testing"

	"github.com/joshmedeski/sesh/v2/icon"
	"github.com/joshmedeski/sesh/v2/lister"
	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/tmux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSSHPreview(t *testing.T) {
	for _, live := range []bool{false, true} {
		t.Run(map[bool]string{false: "host details", true: "live pane"}[live], func(t *testing.T) {
			config := model.Config{SSHConfigs: []model.SSHConfig{{Name: "dev", Host: "alice@dev-host"}}}
			mLister := lister.NewMockLister(t)
			mTmux := tmux.NewMockTmux(t)
			mLister.EXPECT().FindTmuxSession("dev").Return(model.SeshSession{Name: "dev"}, live)
			if live {
				mTmux.EXPECT().CapturePane("dev").Return("remote shell", nil)
			}
			ic := icon.NewIcon(config)
			p := NewPreviewer(mLister, mTmux, ic, nil, nil, nil, config, nil)
			out, err := p.Preview(ic.AddIconNoColor(config.SSHConfigs[0].Session()))
			require.NoError(t, err)
			if live {
				assert.Equal(t, "remote shell", out)
			} else {
				assert.Contains(t, out, "Host: alice@dev-host")
				assert.Contains(t, out, "ssh -- 'alice@dev-host'")
			}
		})
	}
}
