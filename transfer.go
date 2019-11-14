package main

import (
	"fmt"
	"path"
	"strings"

	"github.com/AnotherCoolDude/excelFactory/projecttransfer/basecamp"
	bcmodels "github.com/AnotherCoolDude/excelFactory/projecttransfer/basecamp/models"
	"github.com/AnotherCoolDude/excelFactory/projecttransfer/helper"
	"github.com/manifoldco/promptui"

	"os"

	"github.com/AnotherCoolDude/excelFactory/projecttransfer/proad"
	"github.com/AnotherCoolDude/excelFactory/projecttransfer/proad/models"
	"github.com/joho/godotenv"
	"github.com/urfave/cli"
)

var (
	list           bool
	basecampclient *basecamp.Client
	proadclient    *proad.Client

	transferCommand = &cli.Command{
		Name:    "Transfer",
		Aliases: []string{"t"},
		Usage:   "Transfer transfers todos from basecamp to proad",
		Action:  transferAction,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "list, l",
				Usage:       "list all projects that can be transfered",
				Destination: &list,
			},
		},
	}
)

func transferAction(c *cli.Context) error {
	// init clients
	initClients()
	// fetch basecamp projects and todos
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
	// if list flag has been passed, only show projects
	if list {
		for _, p := range pp {
			fmt.Printf("> %s\n", p.Name)
			for _, t := range p.Todos {
				assigenee := ""
				if len(t.Assignees) > 0 {
					assigenee = t.Assignees[0].Name
				}
				fmt.Printf("\t%s\t%s\n", t.Title, assigenee)
			}
		}
		return nil
	}

	// let user choose which project to transfer
	idx, err := basecampSelection(pp)
	if err != nil {
		fmt.Printf("error occurred: %s", err)
		os.Exit(1)
	}

	// fetch proad project
	fmt.Println("checking for corresponding proad projects...")
	var project models.Project
	err = proadclient.FetchProject(pp[idx].Projectno(), &project)
	if err != nil {
		fmt.Printf("could not fetch project with projectnumber %s.\nDoes the project exist in proad?\n", pp[idx].Projectno())
		os.Exit(1)
	}
	// let user choose the employees responsible for the todos
	fmt.Printf("found project %s\nselect assignees\n", project.Projectno)
	urnos := []int{}
	// get urnos of employees
	for _, t := range pp[idx].Todos {
		// check if assignees exist
		if len(t.Assignees) == 1 {
			if urno, ok := proadclient.Employees[t.Assignees[0].Name]; ok {
				urnos = append(urnos, urno)
				continue
			}
		}
		// if not, let user choose assignee
		urno, _, err := employeeSelection(proadclient.Employees, t.Title)
		if err != nil {
			fmt.Printf("error selecting employee: %s", err)
			os.Exit(1)
		}
		urnos = append(urnos, urno)
	}
	// we now have all infos we need, lets create todos
	fmt.Println("attempting to create a new proad todos")
	for i, t := range pp[idx].Todos {
		pt, err := helper.ProadTodo(t.Title, t.CreatedAt, proadclient.ManagerUrno, project.Urno, urnos[i])
		if err != nil {
			fmt.Printf("could not create todo: %s", err)
			os.Exit(1)
		}
		err = proadclient.PostTodo(pt)
		if err != nil {
			fmt.Printf("could not post todo: %s", err)
			os.Exit(1)
		}
	}
	return nil
}

func initClients() {
	p, err := os.Executable()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	envPath := path.Join(path.Dir(p), ".env")
	fmt.Printf("expecting .env at path %s\n", envPath)
	err = godotenv.Load(envPath)

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

func employeeSelection(employees map[string]int, title string) (int, string, error) {
	names := []string{}
	for e := range employees {
		names = append(names, e)
	}

	searcher := func(input string, index int) bool {
		return strings.Contains(strings.ToLower(names[index]), strings.ToLower(input))
	}

	selection := promptui.Select{
		Label:             "Choose Assignee for " + title,
		Items:             names,
		Searcher:          searcher,
		StartInSearchMode: true,
	}
	i, n, err := selection.Run()
	for name, urno := range employees {
		if n == name {
			i = urno
		}
	}
	return i, n, err
}
