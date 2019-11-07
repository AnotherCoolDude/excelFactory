package basecamp

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/AnotherCoolDude/excelFactory/projecttransfer/basecamp/models"
	"github.com/AnotherCoolDude/excelFactory/projecttransfer/helper"
	"strconv"

	"github.com/rs/xid"
	"golang.org/x/oauth2"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"sync"
)

// Client represents a client which communicates with basecamp
type Client struct {
	email       string
	appName     string
	id          int
	oauthConfig *oauth2.Config
	state       string
	code        string
	token       *oauth2.Token
	httpclient  *http.Client
}

// DefaultClient returns a default client for basecamp
func DefaultClient() *Client {
	return &Client{
		appName: os.Getenv("BASECAMP_APPNAME"),
		email:   os.Getenv("BASECAMP_EMAIL"),
		id:      0,
		state:   xid.New().String(),
		code:    "",
		token:   &oauth2.Token{},
		oauthConfig: &oauth2.Config{
			RedirectURL:  os.Getenv("BASECAMP_CALLBACK"),
			ClientID:     os.Getenv("BASECAMP_CLIENT"),
			ClientSecret: os.Getenv("BASECAMP_SECRET"),
			Scopes:       []string{},
			Endpoint: oauth2.Endpoint{
				AuthStyle: oauth2.AuthStyleAutoDetect,
				AuthURL:   "https://launchpad.37signals.com/authorization/new",
				TokenURL:  "https://launchpad.37signals.com/authorization/token",
			},
		},
		httpclient: http.DefaultClient,
	}
}

// CachedAuth attempts to load a token from os.Env "BASECAMP_CODE" and id from os.Env "BASECAMP_ID"
func (c *Client) CachedAuth() error {
	code := os.Getenv("BASECAMP_CODE")
	id := os.Getenv("BASECAMP_ID")
	if code == "" || id == "" {
		return fmt.Errorf("env variables are empty")
	}
	t, err := c.oauthConfig.Exchange(oauth2.NoContext, code, oauth2.SetAuthURLParam("type", "web_server"))
	if err != nil {
		return err
	}
	if !t.Valid() {
		return fmt.Errorf("cashed token has expired")
	}
	c.token = t
	idnr, err := strconv.Atoi(id)
	if err != nil {
		return fmt.Errorf("error parsing id into int")
	}
	c.id = idnr
	return nil
}

// CacheAuth caches relevant infos in environment variables for convenience
func (c *Client) CacheAuth() error {
	if c.code == "" || c.id == 0 {
		return fmt.Errorf("could not cache basecamp auth: code or id empty")
	}

	var err error
	err = os.Setenv("BASECAMP_ID", strconv.Itoa(c.id))
	if err != nil {
		return err
	}
	os.Setenv("BASECAMP_CODE", c.code)
	if err != nil {
		return err
	}
	return nil
}

// Unmarshal sends a request and returns the unmarshalled response
func (c *Client) Unmarshal(URL string, query map[string]string, model interface{}) error {
	resp, err := c.Do("GET", URL, http.NoBody, query)
	if err != nil {
		return err
	}
	err = unmarshal(resp, &model)
	if err != nil {
		return err
	}
	return nil
}

// Do creates and sends a request
func (c *Client) Do(method, URL string, body io.Reader, query map[string]string) (*http.Response, error) {
	requestURL, err := url.Parse(URL)
	if err != nil {
		fmt.Println("could not parse url")
		return nil, err
	}
	if !requestURL.IsAbs() {
		requestURL = c.baseURL()
		requestURL.Path = path.Join(requestURL.Path, URL)
	}
	req, err := http.NewRequest(method, requestURL.String(), body)
	if err != nil {
		fmt.Println("error making request " + err.Error())
		return nil, err
	}
	c.addHeader(req)
	q := req.URL.Query()
	for key, value := range query {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()
	resp, err := c.httpclient.Do(req)
	if err != nil {
		fmt.Println("error making request " + err.Error())
		return nil, err
	}
	return resp, nil
}

// baseURL returnes the base of every request for basecamp api
func (c *Client) baseURL() *url.URL {
	urlString := fmt.Sprintf("https://3.basecampapi.com/%d/", c.id)
	url, _ := url.Parse(urlString)
	return url
}

// addHeader modifies the header of request to be accepted by basecamp
func (c *Client) addHeader(request *http.Request) {
	request.Header.Add("Authorization", "Bearer "+c.token.AccessToken)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("User-Agent", fmt.Sprintf("%s (%s)", c.appName, c.email))
}

// AuthCodeURL returns the url where the client needs to authenticate
func (c *Client) AuthCodeURL() string {
	return c.oauthConfig.AuthCodeURL(c.state, oauth2.SetAuthURLParam("type", "web_server"))
}

// HandleCallback processes callback from basecamp (and fetches the ID)
func (c *Client) HandleCallback(request *http.Request) error {
	code := request.FormValue("code")
	state := request.FormValue("state")
	if state != c.state {
		return errors.New("[basecamp.go/handleCallback] state doesn't match")
	}
	c.code = code
	t, err := c.oauthConfig.Exchange(oauth2.NoContext, code, oauth2.SetAuthURLParam("type", "web_server"))
	if err != nil {
		return err
	}
	c.token = t
	err = c.receiveID()
	if err != nil {
		return err
	}
	return nil
}

// IsValid returns wether the the basecamp client has a valid token
func (c *Client) IsValid() bool {
	if !c.token.Valid() {
		return false
	}
	if c.id == 0 {
		return false
	}
	return true
}

// receiveID makes a request to basecamp and extracts the ID from the response. The ID is required for further requests
func (c *Client) receiveID() error {
	if !c.token.Valid() {
		return errors.New("need valid token to receive ID")
	}
	resp, err := c.Do("GET", "https://launchpad.37signals.com/authorization.json", http.NoBody, helper.Query{})
	if err != nil {
		return errors.New("[basecamp.go/receiveID] couldn't make request to auth endpoint")
	}
	respbytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("[basecamp.go/receiveID] couldn't read response body")
	}
	defer resp.Body.Close()
	var result map[string]interface{}
	json.Unmarshal(respbytes, &result)
	accounts := result["accounts"].([]interface{})
	accDetails := accounts[0].(map[string]interface{})
	c.id = int(accDetails["id"].(float64))
	return nil
}

// FetchTodos gets todos from basecamp
func (c *Client) FetchTodos(project *models.Project) error {
	setresp, err := c.Do("GET", project.Dock[2].URL, http.NoBody, helper.Query{})
	if err != nil {
		return err
	}
	var set models.Todoset
	err = unmarshal(setresp, &set)
	if err != nil {
		return err
	}
	if set.TodolistsCount == 0 {
		return nil
	}
	listresp, err := c.Do("GET", set.TodolistsURL, http.NoBody, helper.Query{})
	if err != nil {
		return err
	}
	var lists []models.Todolist
	err = unmarshal(listresp, &lists)
	if err != nil {
		return err
	}
	todos := []models.Todo{}
	for _, l := range lists {
		todoresp, err := c.Do("GET", l.TodosURL, http.NoBody, helper.Query{})
		if err != nil {
			return err
		}
		var tt []models.Todo
		err = unmarshal(todoresp, &tt)
		if err != nil {
			return err
		}
		todos = append(todos, tt...)
	}
	project.Todos = todos
	return nil
}

// fetchTodosAsync fetches todos from basecamp asynchronus
func (c *Client) fetchTodosAsync(project *models.Project, sem chan int, wg *sync.WaitGroup, errChan chan error) {
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
	b, e := helper.ResponseBytes(response)
	if e != nil {
		return e
	}

	e = json.Unmarshal(b, &model)

	if e != nil {
		return e
	}
	return nil
}
