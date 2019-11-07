package projecttransfer

import (
	"fmt"
	"github.com/AnotherCoolDude/excelFactory/projecttransfer/basecamp"
	bcmodels "github.com/AnotherCoolDude/excelFactory/projecttransfer/basecamp/models"
	"github.com/AnotherCoolDude/excelFactory/projecttransfer/helper"
	"github.com/AnotherCoolDude/excelFactory/projecttransfer/proad"
	"github.com/joho/godotenv"
	"github.com/pkg/browser"
	"net/http"
	"os"
	"sync"
)

var (
	// Proadclient performs tasks on proad
	Proadclient *proad.Client
	// Basecampclient performs tasks on basecamp
	Basecampclient *basecamp.Client

	server *http.Server
	wg     sync.WaitGroup
)

// InitClients returns initialized clients (start server first)
func InitClients() {
	err := godotenv.Load()

	if err != nil {
		fmt.Printf("error loading .env file: %s\n", err)
		os.Exit(1)
	}
	Proadclient = proad.DefaultClient()
	Basecampclient = basecamp.DefaultClient()

	// check if basecamp's cached auth is available and valid
	err = Basecampclient.CachedAuth()
	if err != nil {
		fmt.Println(err)
	} else {
		// if auth is valid, we dont need the rest
		return
	}

	// init routes
	http.HandleFunc("/basecamp/callback", handleCallback)

	// start server
	startServer()

	// authenticate basecamp, wait for callback
	wg.Add(1)
	err = browser.OpenURL(Basecampclient.AuthCodeURL())

	if err != nil {
		fmt.Printf("could not open browser for authentication: %s", err)
		os.Exit(1)
	}
	fmt.Println("waiting for callback")
	wg.Wait()
}

// StartServer starts the sever needed for receiving the callback from basecamp
func startServer() {
	fmt.Println("starting server")
	server = &http.Server{Addr: ":3000"}
	go server.ListenAndServe()
}

// StopServer stops the running server
func stopServer() {
	fmt.Println("closing server")
	server.Close()
}

// handleCallback handles the callback from basecamp
func handleCallback(w http.ResponseWriter, r *http.Request) {
	err := Basecampclient.HandleCallback(r)
	if err != nil {
		fmt.Printf("error processing callback from basecamp: %s", err)
		os.Exit(1)
	}
	wg.Done()
	err = Basecampclient.CacheAuth()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("successfully received callback")
	stopServer()
}

// TransferTodos creates proad todos based on basecamp todos for a project with projectnr
// func TransferTodos(projectnr string) error {

// }

// FetchBasecampProjects returns user specific projects
func FetchBasecampProjects() ([]bcmodels.Project, error) {
	// get projects from basecamp
	var pp []bcmodels.Project

	// fetch basecampprojects
	err := Basecampclient.Unmarshal("/projects.json", helper.Query{}, &pp)
	if err != nil {
		fmt.Printf("error unmarshalling basecamp projects: %s\n", err)
		return pp, err
	}
	return pp, nil
}
