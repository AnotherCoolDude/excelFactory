package main

import (
	"fmt"
	"github.com/AnotherCoolDude/excelFactory/projecttransfer"
	"github.com/urfave/cli"
)

var (
	list bool

	transferCommand = cli.Command{
		Name:    "Transfer",
		Aliases: []string{"t"},
		Usage:   "Transfer transfers todos from basecamp to proad",
		Action:  transferAction,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:        "list, l",
				Usage:       "list all projects that can be transfered",
				Destination: &list,
			},
		},
	}
)

func transferAction(c *cli.Context) error {
	if !list {
		fmt.Println("only list is supported as of now")
		return nil
	}

	projecttransfer.InitClients()

	pp, err := projecttransfer.FetchBasecampProjects()
	if err != nil {
		return err
	}
	for i, p := range pp {
		err := projecttransfer.Basecampclient.FetchTodos(&p)
		if err != nil {
			fmt.Printf("error fetching todos: %s\n", err)
		}
		fmt.Printf("project %d: %s\n", i, p.Name)
		for _, t := range p.Todos {
			fmt.Printf("\t%s\n", t.Title)
		}
	}
	return nil
}
