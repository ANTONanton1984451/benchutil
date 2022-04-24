package cli

import "errors"

var (
	EmptyCommandsErr  = errors.New("no commands to run")
	NoCommandErr      = errors.New("commands is not provided")
	UnknownCommandErr = errors.New("command is unknown")
)
