package models

import (
	"time"
	"unicode"
)

// Todo is a struct generated from basecamps json response
type Todo struct {
	ID                    int           `json:"id"`
	Status                string        `json:"status"`
	VisibleToClients      bool          `json:"visible_to_clients"`
	CreatedAt             time.Time     `json:"created_at"`
	UpdatedAt             time.Time     `json:"updated_at"`
	Title                 string        `json:"title"`
	InheritsStatus        bool          `json:"inherits_status"`
	Type                  string        `json:"type"`
	URL                   string        `json:"url"`
	AppURL                string        `json:"app_url"`
	BookmarkURL           string        `json:"bookmark_url"`
	SubscriptionURL       string        `json:"subscription_url"`
	CommentsCount         int           `json:"comments_count"`
	CommentsURL           string        `json:"comments_url"`
	Position              int           `json:"position"`
	Parent                Parent        `json:"parent"`
	Bucket                Bucket        `json:"bucket"`
	Creator               Creator       `json:"creator"`
	Description           string        `json:"description"`
	Completed             bool          `json:"completed"`
	Content               string        `json:"content"`
	StartsOn              string        `json:"starts_on"`
	DueOn                 string        `json:"due_on"`
	Assignees             []Assignee    `json:"assignees"`
	CompletionSubscribers []interface{} `json:"completion_subscribers"`
	CompletionURL         string        `json:"completion_url"`
}

// Assignee represents a person, that is assigned to a todo
type Assignee struct {
	ID             int         `json:"id"`
	AttachableSgid string      `json:"attachable_sgid"`
	Name           string      `json:"name"`
	EmailAddress   string      `json:"email_address"`
	PersonableType string      `json:"personable_type"`
	Title          string      `json:"title"`
	Bio            interface{} `json:"bio"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
	Admin          bool        `json:"admin"`
	Owner          bool        `json:"owner"`
	TimeZone       string      `json:"time_zone"`
	AvatarURL      string      `json:"avatar_url"`
}

// Projectno returns the projectno of the todo
func (t *Todo) Projectno() string {
	if len(t.Bucket.Name) < 14 {
		return ""
	}
	nr := t.Bucket.Name[:14]
	for i := 0; i < 4; i++ {
		r := rune(nr[i])
		if !unicode.IsUpper(r) {
			return ""
		}
	}
	return nr
}

// Timestamp is a identifier for comparing with other todos
func (t Todo) Timestamp() string {
	return t.CreatedAt.Format(time.RFC3339)
}

// Identifier returns a unique identifier
func (t Todo) Identifier() int {
	return t.ID
}

// ClientType returns the type of Todo
func (t Todo) ClientType() string {
	return "basecamp"
}
