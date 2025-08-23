//go:build tools
// +build tools

// workaround to keep xcaddy in auto-update by dependabot

package tools

import (
	_ "github.com/caddyserver/xcaddy"
)
