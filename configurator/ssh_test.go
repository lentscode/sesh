package configurator

import (
	"fmt"
	"testing"

	"github.com/joshmedeski/sesh/v2/model"
	"github.com/joshmedeski/sesh/v2/pathwrap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSSHConfig(t *testing.T) {
	config, err := configFromTOML(t, `
strict_mode = true
[[ssh]]
name = "dev-server"
host = "alice@dev.example.com"
[icons.ssh]
icon = "S"
color = "magenta"
`)
	require.NoError(t, err)
	assert.Equal(t, []model.SSHConfig{{Name: "dev-server", Host: "alice@dev.example.com"}}, config.SSHConfigs)
}

func TestInvalidSSHConfig(t *testing.T) {
	for _, entry := range []model.SSHConfig{
		{Name: "", Host: "dev"},
		{Name: "dev.example.com", Host: "dev"},
		{Name: "dev:1", Host: "dev"},
		{Name: " dev", Host: "dev"},
		{Name: "dev\n", Host: "dev"},
		{Name: "dev", Host: ""},
		{Name: "dev", Host: "-V"},
		{Name: "dev", Host: "dev -p 2222"},
		{Name: "dev", Host: "dev\nexit"},
		{Name: "dev", Host: "*.example.com"},
	} {
		t.Run(fmt.Sprintf("%q/%q", entry.Name, entry.Host), func(t *testing.T) {
			_, err := configFromTOML(t, fmt.Sprintf("[[ssh]]\nname = %q\nhost = %q\n", entry.Name, entry.Host))
			assert.ErrorContains(t, err, "invalid ssh")
		})
	}
}

func TestSSHImportsAndValidation(t *testing.T) {
	for _, explicitPath := range []bool{false, true} {
		for _, duplicate := range []bool{false, true} {
			t.Run(fmt.Sprintf("explicit=%t/duplicate=%t", explicitPath, duplicate), func(t *testing.T) {
				main := "import = [\"~/hosts.toml\"]\n"
				if duplicate {
					main += "[[ssh]]\nname = \"dev\"\nhost = \"other\"\n"
				}
				mockOs := &testOs{homeDir: "/home/test", files: map[string][]byte{
					"/home/test/.config/sesh/sesh.toml": []byte(main),
					"/home/test/hosts.toml":             []byte("[[ssh]]\nname = \"dev\"\nhost = \"dev-host\"\n"),
				}}
				path := ""
				if explicitPath {
					path = "/home/test/.config/sesh/sesh.toml"
				}
				config, err := NewConfiguratorWithPath(mockOs, pathwrap.NewPath(), nil, path).GetConfig()
				if duplicate {
					assert.ErrorContains(t, err, "duplicate ssh name")
				} else {
					require.NoError(t, err)
					assert.Equal(t, []model.SSHConfig{{Name: "dev", Host: "dev-host"}}, config.SSHConfigs)
				}
			})
		}
	}
}

func TestSSHNameConflicts(t *testing.T) {
	for _, field := range []string{"name", "alias"} {
		_, err := configFromTOML(t, fmt.Sprintf("[[session]]\n%s = \"dev\"\n[[ssh]]\nname = \"dev\"\nhost = \"dev-host\"\n", field))
		assert.ErrorContains(t, err, "conflicts with a configured session")
	}
}
