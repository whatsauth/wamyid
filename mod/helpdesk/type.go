package helpdesk

type Helpdesk struct {
	Team         string `json:"team,omitempty" bson:"team,omitempty"`
	Scope        string `json:"scope,omitempty" bson:"scope,omitempty"`
	Name         string `json:"name,omitempty" bson:"name,omitempty"`
	Phonenumbers string `json:"phonenumbers,omitempty" bson:"phonenumbers,omitempty"`
	Occupied     bool   `json:"occupied,omitempty" bson:"occupied,omitempty"`
}

type HelpdeskUser struct {
	Team         string `json:"team,omitempty" bson:"team,omitempty"`
	Scope        string `json:"scope,omitempty" bson:"scope,omitempty"`
	Name         string `json:"name,omitempty" bson:"name,omitempty"`
	Phonenumbers string `json:"phonenumbers,omitempty" bson:"phonenumbers,omitempty"`
	Occupied     bool   `json:"occupied,omitempty" bson:"occupied,omitempty"`
}
