package meet

import (
	"benchutil/pkg/cli"
	"context"
	"fmt"
)

func Command() cli.Command {
	var name string
	return cli.Command{
		Name:        "meet",
		Description: "Приветствует пользователя",
		Flags: []cli.CmdFlag{
			cli.StringFlag{
				Name:        "name",
				Destination: &name,
				Default:     "guest",
				Usage:       "Имя для приветствия",
			},
		},
		Action: func(ctx context.Context) error {
			return action(ctx, name)
		},
	}
}

func action(ctx context.Context, name string) error {
	fmt.Printf("hello, %s! \n", name)
	return nil
}
