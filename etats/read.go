package json

import (
	"encoding/json"
	"io/ioutil"
)

// ReadFile Directors reads the file at path and returns the etat directors
func ReadFile(path string) (Etats, error) {
	var ee Etats
	bb, err := ioutil.ReadFile(path)
	if err != nil {
		return ee, err
	}
	json.Unmarshal(bb, &ee)
	return ee, nil
}
