package module

type Module struct {
	Name    string   `json:"name,omitempty" bson:"name,omitempty"`
	Keyword []string `json:"keyword,omitempty" bson:"keyword,omitempty"`
}
