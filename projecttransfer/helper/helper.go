package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	proadDateTimeFormat = "2006-01-02T15:04:05"
)

// Query is a convenience for map[string]string
type Query map[string]string

// ResponseBytes returns the json body of response as bytes
func ResponseBytes(response *http.Response) ([]byte, error) {
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []byte{}, err
	}
	defer response.Body.Close()
	return bytes, nil
}

// // CreateProadTodo creates a proad todo from a basecamptodo and the project, to which the todo is going to belong to
// func CreateProadTodo(basecampTodo bcmodels.Todo, paproject pamodels.Project, managerUrno int) interface{} {
// 	fmt.Println("creating...")
// 	todo := struct {
// 		Shortinfo     string `json:"shortinfo"`
// 		ProjectUrno   int    `json:"urno_project"`
// 		ManagerUrno   int    `json:"urno_manager"`
// 		FromDatetime  string `json:"from_datetime"`
// 		UntilDatetime string `json:"until_datetime"`
// 	}{
// 		Shortinfo:     basecampTodo.Title,
// 		ProjectUrno:   paproject.Urno,
// 		ManagerUrno:   managerUrno,
// 		FromDatetime:  basecampTodo.CreatedAt.Format(proadDateTimeFormat),
// 		UntilDatetime: basecampTodo.CreatedAt.Add(168 * time.Hour).Format(proadDateTimeFormat),
// 	}
// 	fmt.Printf("%+v\n", todo)
// 	return todo
// }

// ProadTodo returns an todo that can be used to create a new todo
func ProadTodo(info string, startDate time.Time, managerUrno, projectUrno, responsibleUrno int) (interface{}, error) {
	endDateDays, err := strconv.Atoi(os.Getenv("DEFAULT_DAYS"))
	if err != nil {
		return nil, fmt.Errorf("could not parse global var DEFAULT_DAYS: %s", err)
	}
	todo := struct {
		Shortinfo       string `json:"shortinfo"`
		ProjectUrno     int    `json:"urno_project"`
		ManagerUrno     int    `json:"urno_manager"`
		ResponsibleUrno int    `json:"urno_responsible"`
		FromDatetime    string `json:"from_datetime"`
		UntilDatetime   string `json:"until_datetime"`
	}{
		Shortinfo:       info,
		ProjectUrno:     projectUrno,
		ManagerUrno:     managerUrno,
		ResponsibleUrno: responsibleUrno,
		FromDatetime:    startDate.Format(proadDateTimeFormat),
		UntilDatetime:   startDate.Add(time.Duration(endDateDays) * 24 * time.Hour).Format(proadDateTimeFormat),
	}
	return todo, nil
}

// PrettyPrintBytes prints out bytes (e.g. from a response) in a readable way
func PrettyPrintBytes(bb []byte) error {
	var jsonPretty bytes.Buffer
	err := json.Indent(&jsonPretty, bb, "", "\t")
	if err != nil {
		return err
	}
	fmt.Println(string(jsonPretty.Bytes()))
	return nil
}
