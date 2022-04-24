package cli

import (
	"context"
	"os"
)

type App struct {
	commands []executable
}

// NewApp создаёт экземпляр приложения cli
// при отсутствии команд приложение не создастся а вернётся ошибка EmptyCommandsErr
func NewApp(cmds ...Command) (App, error) {
	if len(cmds) == 0 {
		return App{}, EmptyCommandsErr
	}
	executeList := make([]executable, 0, len(cmds))
	for _, c := range cmds {
		executeList = append(executeList, c)
	}

	help := helpCommand{name: "help", commands: cmds}

	executeList = append(executeList, help)

	return App{commands: executeList}, nil
}

// Run запуск приложения, если в os.Args один и меньше элементов,
// то не задана команда и запуск приложения вернёт NoCommandErr
// если нужно команды не надётся, Run вернёт UnknownCommandErr
func (app App) Run(ctx context.Context) error {
	if len(os.Args) == 1 {
		return NoCommandErr
	}

	commandName := os.Args[1]
	for _, c := range app.commands {
		if c.getName() == commandName {
			return c.execute(ctx, os.Args[2:])
		}
	}

	return UnknownCommandErr
}
