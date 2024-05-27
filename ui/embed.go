package ui

import (
	_ "embed"
)

//go:embed icons/status-active.ico
var ActiveIcon []byte

//go:embed icons/status-inactive.ico
var InactiveIcon []byte
