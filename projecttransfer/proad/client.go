package proad

import (
	"crypto/tls"
	"encoding/json"
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
	httpClient *http.Client
	apiKey     string
}

// DefaultClient returns a default client for basecamp
func DefaultClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
		apiKey: os.Getenv("PROAD_APIKEY"),
	}
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
