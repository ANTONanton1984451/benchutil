package cli

import (
	"context"
	"errors"
	"flag"
)

type executable interface {
	execute(ctx context.Context, args []string) error
	getName() string
}

type Command struct {
	Name        string
	Description string
	Flags       []CmdFlag
	Action      ActionFunc

	fs *flag.FlagSet
}

// ActionFunc функция в котором происходит действие команды
// без инициализации данной функции команда не запустится
type ActionFunc func(ctx context.Context) error

func (c Command) execute(ctx context.Context, args []string) error {
	if c.Action == nil {
		return errors.New("action for command not set")
	}

	if c.fs == nil {
		c.fs = flag.NewFlagSet(c.Name, flag.ContinueOnError)
	}

	if len(c.Flags) > 0 {
		for _, f := range c.Flags {
			f.bind(c.fs)
		}
		if err := c.fs.Parse(args); err != nil {
			return err
		}
	}
	return c.Action(ctx)
}

func (c Command) getName() string {
	return c.Name
}
