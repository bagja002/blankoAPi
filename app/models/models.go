package models

type Admin struct {
	Name      string `json:"name"`
	ID        string `json:"id"`
	Role      string `json:"role"`
	Satminkal string `json:"satminkal"`
}

// Moderator represents the moderator entity
type Moderator struct {
	Name      string `json:"name"`
	ID        string `json:"id"`
	Role      string `json:"role"`
	Satminkal string `json:"satminkal"`
}

// User represents the user entity
type User struct {
	Name      string `json:"name"`
	ID        string `json:"id"`
	Role      string `json:"role"`
	Satminkal string `json:"satminkal"`
}
