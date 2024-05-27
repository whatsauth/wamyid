package module

type Module struct {
	Name    string   `json:"name,omitempty" bson:"name,omitempty"`
	Keyword []string `json:"keyword,omitempty" bson:"keyword,omitempty"`
}
type Typo struct {
	From string `json:"from,omitempty" bson:"from,omitempty"`
	To   string `json:"to,omitempty" bson:"to,omitempty"`
}
