package json

// Etats wraps a etat director in a struct
type Etats struct {
	Directors []struct {
		Name      string   `json:"Name"`
		Campaigns []string `json:"Campaigns"`
	} `json:"Directors"`
}
