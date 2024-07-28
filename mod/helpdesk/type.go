package helpdesk

type Helpdesk struct {
	Team         string `json:"team,omitempty" bson:"team,omitempty"`
	Scope        string `json:"scope,omitempty" bson:"scope,omitempty"`
	Name         string `json:"name,omitempty" bson:"name,omitempty"`
	Phonenumbers string `json:"phonenumbers,omitempty" bson:"phonenumbers,omitempty"`
	Occupied     bool   `json:"occupied,omitempty" bson:"occupied,omitempty"`
}

type User struct {
	Team         string   `json:"team,omitempty" bson:"team,omitempty"`
	Scope        string   `json:"scope,omitempty" bson:"scope,omitempty"`
	Name         string   `json:"name,omitempty" bson:"name,omitempty"`
	Phonenumbers string   `json:"phonenumbers,omitempty" bson:"phonenumbers,omitempty"`
	Terlayani    bool     `json:"terlayani,omitempty" bson:"terlayani,omitempty"`
	Masalah      string   `json:"masalah,omitempty" bson:"masalah,omitempty"`
	Solusi       string   `json:"solusi,omitempty" bson:"solusi,omitempty"`
	Operator     Helpdesk `json:"operator,omitempty" bson:"operator,omitempty"`
}
