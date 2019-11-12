package main

import (
	"fmt"
	"github.com/AnotherCoolDude/excelFactory/projecttransfer/basecamp"
	bcmodels "github.com/AnotherCoolDude/excelFactory/projecttransfer/basecamp/models"
	"github.com/manifoldco/promptui"
	"strings"

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
	fmt.Println("fetching basecamp projects...")
	pp, err := basecampclient.FetchProjects()
	if err != nil {
		return err
	}
	fmt.Println("fetching basecamp todos...")
	for i := range pp {
		err := basecampclient.FetchTodos(&pp[i])
		if err != nil {
			fmt.Printf("error fetching todos: %s\n", err)
		}
	}

	idx, err := basecampSelection(pp)

	if err != nil {
		fmt.Printf("error occurred: %s", err)
		os.Exit(1)
	}
	fmt.Printf("you choose %s\n", pp[idx].Name)

	fmt.Println("checking for corresponding proad projects...")

	var project models.Project
	proadclient.FetchProject(pp[idx].Projectno(), &project)
	proadclient.FetchTodos(&project)
	fmt.Printf("found project %s\nselect assignee\n", project.Projectno)
	urno, err := employeeSelection(proadclient.Employees)
	if err != nil {
		fmt.Printf("error selecting employee: %s", err)
		os.Exit(1)
	}
	fmt.Printf("you choose urno %d", urno)
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

func basecampSelection(projects []bcmodels.Project) (int, error) {
	tt := promptui.SelectTemplates{
		Label:    "Choose a project",
		Active:   `▸ {{ .Name | blue }}`,
		Inactive: `{{ .Name | blue }}`,
		Selected: `▸ {{ .Name | blue | cyan}}`,
		Details: `{{ range .Todos }}
		{{ .Title | blue }}{{ end }}`,
	}

	searcher := func(input string, index int) bool {
		p := projects[index]
		name := strings.Replace(strings.ToLower(p.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	selection := promptui.Select{
		Label:     "Select Project",
		Items:     projects,
		Templates: &tt,
		Searcher:  searcher,
	}

	i, _, err := selection.Run()
	return i, err
}

func employeeSelection(employees map[string]int) (int, error) {
	names := []string{}
	for e := range employees {
		names = append(names, e)
	}

	searcher := func(input string, index int) bool {
		for name := range employees {
			return strings.Contains(strings.ToLower(name), strings.ToLower(input))
		}
		return false
	}
	selection := promptui.Select{
		Label:             "Choose Assignee",
		Items:             names,
		Searcher:          searcher,
		StartInSearchMode: true,
	}
	i, n, err := selection.Run()
	for name, urno := range employees {
		if n == name {
			i = urno
			fmt.Println(urno)
		}
	}
	return i, err
}
