package helpdesk

type Helpdesk struct {
	Name         string   `json:"name,omitempty" bson:"name,omitempty"`
	Phonenumbers []string `json:"phonenumbers,omitempty" bson:"phonenumbers,omitempty"`
	Occupied     bool     `json:"occupied,omitempty" bson:"occupied,omitempty"`
}
