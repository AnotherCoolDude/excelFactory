package main

import (
	"fmt"
	"github.com/AnotherCoolDude/excelFactory/projecttransfer/basecamp"

	"github.com/AnotherCoolDude/excelFactory/projecttransfer/proad"
	"github.com/AnotherCoolDude/excelFactory/projecttransfer/proad/models"
	"github.com/joho/godotenv"
	"github.com/urfave/cli"
	"os"
)

var (
	list           bool
	basecampclient *basecamp.Client
	proadclient    *proad.Client

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

	initClients()

	pp, err := basecampclient.FetchProjects()
	if err != nil {
		return err
	}
	for i := range pp {
		err := basecampclient.FetchTodos(&pp[i])
		if err != nil {
			fmt.Printf("error fetching todos: %s\n", err)
		}
		fmt.Printf("project %d: %s\n", i, pp[i].Name)
		for _, t := range pp[i].Todos {
			if len(t.Assignees) > 0 {
				fmt.Printf("\t%s\t%s\n", t.Title, t.Assignees[0].Name)
			} else {
				fmt.Printf("\t%s\n", t.Title)
			}
		}
	}
	fmt.Println("checking for corresponding proad projects...")
	proadpp := []models.Project{}
	for i, p := range pp {
		if p.Projectno() == "" {
			continue
		}
		var project models.Project
		proadclient.FetchProject(p.Projectno(), &project)
		proadclient.FetchTodos(&project)
		proadpp = append(proadpp, project)
		fmt.Printf("project %d: %s\n", i, project.Projectno)
		for _, t := range project.Todos {
			fmt.Printf("\t%s\n", t.Shortinfo)
		}
	}
	// fmt.Println("attempt to create a new proad todo")
	// todo := helper.CreateProadTodo(pp[1].Todos[0], proadpp[0], proadclient.ManagerUrno)
	// err = proadclient.CreateTodo(todo)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	return nil
}

func initClients() {
	err := godotenv.Load()

	if err != nil {
		fmt.Printf("error loading .env file: %s\n", err)
		os.Exit(1)
	}
	proadclient = proad.DefaultClient()
	basecampclient = basecamp.DefaultClient()

	// Authenticate
	err = basecampclient.Authenticate()
	if err != nil {
		fmt.Printf("failed to authenticate basecamp client: %s", err)
		os.Exit(1)
	}
}
