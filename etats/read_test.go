package etats

import (
	"testing"
)

func TestRead(t *testing.T) {
	path := "etatdirector.json"
	etats, err := ReadFile(path)
	if err != nil {
		t.Errorf("error reading file at path %s: %s", path, err)
	}

	for _, d := range etats.Directors {
		t.Logf("Director %s has the following campaigns: %v", d.Name, d.Campaigns)
	}
}
