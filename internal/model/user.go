package model

// Registration request body
type RegisterRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// Login request body
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Authentication response
type AuthResponse struct {
	Token   string      `json:"token"`
	User    UserSummary `json:"user"`
	Profile interface{} `json:"profile"`
}

// Safe user data (no password)
type UserSummary struct {
	UserID    int    `json:"user_id"`
	Username  string `json:"username"`
	Role      string `json:"role"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
