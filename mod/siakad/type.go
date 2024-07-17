package siakad

type Config struct {
	PhoneNumber    string `json:"phonenumber,omitempty" bson:"phonenumber,omitempty"`
	SiakadLoginURL string `json:"siakadloginurl,omitempty" bson:"siakadloginurl,omitempty"`
	BapURL         string `json:"bapurl,omitempty" bson:"bapurl,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
