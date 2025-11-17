package models

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"` // Don't send password in JSON
	Role     string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
