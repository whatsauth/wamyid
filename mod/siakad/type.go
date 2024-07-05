package siakad

type Config struct {
	SiakadLoginURL string `bson:"siakad_login_url"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
