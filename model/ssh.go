package model

import "strings"

// Session keeps the SSH destination separate from Path, which always refers
// to a local directory. Quote the host as one shell argument because tmux
// starts SSH by typing the command into the new session's shell.
func (c SSHConfig) Session() SeshSession {
	return SeshSession{
		Src:            "ssh",
		Name:           c.Name,
		SSHHost:        c.Host,
		StartupCommand: "ssh -- '" + strings.ReplaceAll(c.Host, "'", `'\''`) + "'",
	}
}
