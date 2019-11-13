package models

// Todo wraps the proad response for todos in a struct
type Todo struct {
	Urno             int         `json:"urno"`
	Company          Company     `json:"company"`
	Project          *Project    `json:"project"`
	ProjectUrno      int         `json:"urno_project"`
	ManagerUrno      int         `json:"urno_manager"`
	ServiceCode      ServiceCode `json:"service_code"`
	Responsible      Responsible `json:"responsible"`
	Manager          Manager     `json:"manager"`
	Shortinfo        string      `json:"shortinfo"`
	FromDatetime     string      `json:"from_datetime"`
	UntilDatetime    string      `json:"until_datetime"`
	ReminderDatetime string      `json:"reminder_datetime"`
	Status           string      `json:"status"`
	Priority         string      `json:"priority"`
	Description      string      `json:"description"`
}

// Company wraps the proad response for Company in a struct
type Company struct {
	Urno       int    `json:"urno"`
	Shortname  string `json:"shortname"`
	Name       string `json:"name"`
	Firstname  string `json:"firstname"`
	Lastname   string `json:"lastname"`
	Type       string `json:"type"`
	Active     bool   `json:"active"`
	ExternalID string `json:"external_id"`
}

// Manager wraps the proad response for manager in a struct
type Manager struct {
	Urno              int                `json:"urno"`
	Shortname         string             `json:"shortname"`
	Name              string             `json:"name"`
	PrivateMainAdress PrivateMainAddress `json:"private_main_address"`
	Firstname         string             `json:"firstname"`
	Lastname          string             `json:"lastname"`
	Type              string             `json:"type"`
	Active            bool               `json:"active"`
	ExternalID        string             `json:"external_id"`
}

// ServiceCode wraps the proad response for servicecode in a struct
type ServiceCode struct {
	Urno                      int    `json:"urno"`
	Shortname                 string `json:"shortname"`
	Name                      string `json:"name"`
	Istimecost                bool   `json:"istimecost"`
	Useintimeregistration     bool   `json:"useintimeregistration"`
	Isexternalcost            bool   `json:"isexternalcost"`
	Useinpurchaseinvoice      bool   `json:"useinpurchaseinvoice"`
	Ismaterialcost            bool   `json:"ismaterialcost"`
	Useinmaterialregistration bool   `json:"useinmaterialregistration"`
	Iscategory1               bool   `json:"iscategory1"`
	Iscategory2               bool   `json:"iscategory2"`
	Iscategory3               bool   `json:"iscategory3"`
}

// Responsible wraps the proad response for the responsible manager in a struct
type Responsible struct {
	Urno               int                `json:"urno"`
	Shortname          string             `json:"shortname"`
	Firstname          string             `json:"firstname"`
	Lastname           string             `json:"lastname"`
	Type               string             `json:"type"`
	PrivateMainAddress PrivateMainAddress `json:"private_main_address"`
	Salutation         string             `json:"salutation"`
	Title              string             `json:"title"`
	Gender             string             `json:"gender"`
	Department         string             `json:"department"`
	Function           string             `json:"function"`
	Business1          string             `json:"business1"`
	Business2          string             `json:"business2"`
	Birthday           string             `json:"birthday"`
	Active             bool               `json:"active"`
	ExternalID         string             `json:"external_id"`
}

// PrivateMainAddress wraps the proad response for the private main adressof a person in a struct
type PrivateMainAddress struct {
	Urno    int     `json:"urno"`
	Street  string  `json:"street"`
	Zipcode string  `json:"zipcode"`
	City    string  `json:"city"`
	Phone   string  `json:"phone"`
	Fax     string  `json:"fax"`
	Mobile  string  `json:"mobile"`
	Country Country `json:"country"`
	Email   string  `json:"email"`
	Type    string  `json:"type"`
}

// Country wraps the proad response for a country in a struct
type Country struct {
	Urno        int    `json:"urno"`
	CountryName string `json:"country_name"`
	Shortname   string `json:"shortname"`
}

// PostTodo is a Todo that contains fields neccessary for creating a todo
type PostTodo struct {
	Shortinfo       string `json:"shortinfo"`
	ProjectUrno     int    `json:"urno_project"`
	ManagerUrno     int    `json:"urno_manager"`
	ResponsibleUrno int    `json:"urno_responsible"`
	FromDatetime    string `json:"from_datetime"`
	UntilDatetime   string `json:"until_datetime"`
}
