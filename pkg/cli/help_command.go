package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
)

type helpCommand struct {
	name     string
	commands []Command
}

func (hc helpCommand) getName() string {
	return hc.name
}

func (hc helpCommand) execute(ctx context.Context, args []string) error {
	if len(args) > 1 {
		return errors.New("unknown command")
	}

	var observeCommands []Command
	if len(args) == 0 {
		for _, c := range hc.commands {
			observeCommands = append(observeCommands, c)
		}
	} else {
		var found bool
		needleCommand := args[0]
		for _, c := range hc.commands {
			if c.Name == needleCommand {
				observeCommands = append(observeCommands, c)
				found = true
			}
		}

		if !found {
			return fmt.Errorf("command %s not found", needleCommand)
		}

	}
	var commandsDescription string
	for i, oc := range observeCommands {
		commandsDescription += hc.makeDescription(oc)
		if i < len(observeCommands)-1 {
			commandsDescription += "-----------------------\n"
		}
	}

	os.Stdout.WriteString(commandsDescription)

	return nil
}

func (hc helpCommand) makeDescription(cmd Command) string {
	flagTmpl := "@padding-%s   %s\n@paddingЗначение по умолчанию - \"%v\" \n"
	flagTmpl = strings.Replace(flagTmpl, "@padding", "     ", -1)

	var description string
	description += fmt.Sprintf("Команда: %s   %s \n", cmd.Name, cmd.Description)
	if len(cmd.Flags) > 0 {
		description += "Флаги:\n"
		for _, f := range cmd.Flags {
			description += fmt.Sprintf(flagTmpl, f.name(), f.usage(), f.defaultVal())
		}
	}

	return description
}
