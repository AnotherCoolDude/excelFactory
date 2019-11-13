package proad

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/AnotherCoolDude/excelFactory/projecttransfer/helper"
	"github.com/AnotherCoolDude/excelFactory/projecttransfer/proad/models"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"sync"
)

// Client represents a client which communicates with proad
type Client struct {
	httpClient  *http.Client
	apiKey      string
	ManagerUrno int
	Employees   map[string]int
}

// DefaultClient returns a default client for basecamp
func DefaultClient() *Client {
	c := &Client{
		httpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
		apiKey: os.Getenv("PROAD_APIKEY"),
	}
	u, err := c.managerUrno()
	if err != nil {
		fmt.Printf("could not retrieve manager urno: %s\n", err)
	}
	c.ManagerUrno = u

	ee, err := c.fetchEmployees()
	if err != nil {
		fmt.Printf("could not retrieve employees: %s\n", err)
	}
	c.Employees = ee
	return c
}

// Do creates and sends a request
func (c *Client) Do(method, URL string, body io.Reader, query map[string]string) (*http.Response, error) {
	requestURL, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}
	if !requestURL.IsAbs() {
		requestURL, _ = url.Parse("https://192.168.0.15/api/v5/")
		requestURL.Path = path.Join(requestURL.Path, URL)
	}
	req, err := http.NewRequest(method, requestURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("apikey", c.apiKey)
	q := req.URL.Query()
	for key, value := range query {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ManagerUrno returns the urno of this clients user
func (c *Client) managerUrno() (int, error) {
	resp, err := c.Do("GET", "/me", http.NoBody, helper.Query{})
	if err != nil {
		return 0, err
	}
	var Manager struct {
		Urno     int    `json:"urno"`
		Endpoint string `json:"endpoint"`
	}
	bb, err := helper.ResponseBytes(resp)
	if err != nil {
		return 0, err
	}
	err = json.Unmarshal(bb, &Manager)
	if err != nil {
		return 0, err
	}
	return Manager.Urno, nil
}

// FetchEmployees returns a map that contains the names as keys and their urnos as values
func (c *Client) fetchEmployees() (map[string]int, error) {
	ee := map[string]int{}
	resp, err := c.Do("GET", "/staffs", http.NoBody, helper.Query{})
	if err != nil {
		return ee, err
	}
	var Employees struct {
		EmployeeList []struct {
			Urno      int    `json:"urno"`
			Firstname string `json:"firstname"`
			Lastname  string `json:"lastname"`
		} `json:"person_list"`
	}

	bb, err := helper.ResponseBytes(resp)
	if err != nil {
		return ee, err
	}
	err = json.Unmarshal(bb, &Employees)
	if err != nil {
		return ee, err
	}
	for _, e := range Employees.EmployeeList {
		ee[e.Firstname+" "+e.Lastname] = e.Urno
	}
	return ee, nil

}

// PostTodo creates a new todo in proad
func (c *Client) PostTodo(todo models.PostTodo) error {
	bb, err := json.Marshal(todo)
	if err != nil {
		return fmt.Errorf("unable to marshal todo: %s", err)
	}
	resp, err := c.Do("POST", "/tasks", bytes.NewBuffer(bb), helper.Query{})
	if err != nil {
		return fmt.Errorf("unable to send a post request: %s", err)
	}

	bb, err = helper.ResponseBytes(resp)
	if err != nil {
		return fmt.Errorf("unable to unmarshal response: %s", err)
	}
	var responseStruct struct {
		Error string `json:"error"`
		Urno  int    `json:"urno"`
	}
	err = json.Unmarshal(bb, &responseStruct)
	if err != nil {
		return fmt.Errorf("unable to unmarshal response: %s", err)
	}
	if responseStruct.Error != "" {
		return fmt.Errorf("got error response posting todo: %s", responseStruct.Error)
	}
	assignee := ""
	for n, u := range c.Employees {
		if u == todo.ResponsibleUrno {
			assignee = n
		}
	}
	fmt.Printf("todo %s created with assignee %s (urno: %d)\n", todo.Shortinfo, assignee, responseStruct.Urno)
	return nil
}

// FetchProject gets a proad project by projectnumber
func (c *Client) FetchProject(projectno string, project *models.Project) error {
	projectresp, err := c.Do("GET", "projects", http.NoBody, helper.Query{"projectno": projectno})
	if err != nil {
		return err
	}
	var pp []models.Project
	err = unmarshal(projectresp, &pp)
	*project = pp[0]
	if err != nil {
		return err
	}
	return nil
}

// FetchProjectAsnyc fetches a proad project by projectnumber asynchronus
func (c *Client) fetchProjectAsync(projectno string, project *models.Project, sem chan int, wg *sync.WaitGroup, errChan chan error) {
	defer wg.Done()
	sem <- 1
	if err := c.FetchProject(projectno, project); err != nil {
		select {
		case errChan <- err:
			// we are the first worker to fail
		default:
			// there allready happend an error
		}
	}
	<-sem
}

// FetchTodos fetches todos from project
func (c *Client) FetchTodos(project *models.Project) error {
	todosresp, err := c.Do("GET", "tasks", http.NoBody, helper.Query{"project": strconv.Itoa(project.Urno)})
	if err != nil {
		return err
	}
	var todos []models.Todo
	err = unmarshal(todosresp, &todos)
	if err != nil {
		return err
	}
	for i := range todos {
		todos[i].Project = project
	}
	project.Todos = todos
	return nil
}

// FetchTodosAsync fetches todos from project asynchronus
func (c *Client) FetchTodosAsync(project *models.Project, sem chan int, wg *sync.WaitGroup, errChan chan error) {
	defer wg.Done()
	sem <- 1
	if err := c.FetchTodos(project); err != nil {
		select {
		case errChan <- err:
			// we are the first worker to fail
		default:
			// there allready happend an error
		}
	}
	<-sem
}

// unmarshal parses the json body of response into a struct
func unmarshal(response *http.Response, model interface{}) error {
	var dd map[string]interface{}

	b, e := helper.ResponseBytes(response)
	if e != nil {
		return e
	}
	e = json.Unmarshal(b, &dd)
	if e != nil {
		return e
	}

	var d []byte
	for _, v := range dd {
		d, e = json.Marshal(v)
		if e != nil {
			return e
		}
	}
	json.Unmarshal(d, &model)
	return nil
}
