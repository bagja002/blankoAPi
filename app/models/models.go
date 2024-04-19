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

type SuperAdmin struct{ 
	IdSuperAdmin uint  `json:"id_admin"`
	Nama string	`json:"nama"`
	Email string`json:"email"`
	Username string`json:"username"`
	Password string`json:"password"`
	Status string`json:"status"`
}

//Entity Untuk Lemdiklat

type AdminPusat struct{
	IdAdminPusat uint 
	Nama string
	Email string
	Password string
	NoTelpon string
	Nip string	
	Status string
}