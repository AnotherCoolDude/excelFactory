package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/AnotherCoolDude/excelFactory/projecttransfer/proad/models"
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

// ProadTodo returns an todo that can be used to create a new todo
func ProadTodo(title string, startDate time.Time, managerUrno, projectUrno, responsibleUrno int) (models.PostTodo, error) {
	endDateDays, err := strconv.Atoi(os.Getenv("DEFAULT_DAYS"))
	todo := models.PostTodo{}
	if err != nil {
		return todo, fmt.Errorf("could not parse global var DEFAULT_DAYS: %s", err)
	}
	todo = models.PostTodo{
		Shortinfo:       title,
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
