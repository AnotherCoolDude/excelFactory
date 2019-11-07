package helper

import (
	"io/ioutil"
	"net/http"
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
