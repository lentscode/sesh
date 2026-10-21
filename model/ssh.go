package model

import "strings"

// Session keeps the SSH destination separate from Path, which always refers
// to a local directory. StartupCommand is quoted for display in the preview;
// connections pass SSHHost directly to SSH as a single argument.
func (c SSHConfig) Session() SeshSession {
	return SeshSession{
		Src:            "ssh",
		Name:           c.Name,
		SSHHost:        c.Host,
		StartupCommand: "ssh -- '" + strings.ReplaceAll(c.Host, "'", `'\''`) + "'",
	}
}
