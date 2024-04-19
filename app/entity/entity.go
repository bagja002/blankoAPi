package entity


type Users struct{
	IdUsers uint 		`gorm:"primary_key;auto_increment"`
	Nama	string

	Status 	string
}

type SuperAdmin struct {
    IdSuperAdmin uint `gorm:"primary_key;auto_increment"`
    Nama         string
    Email        string
    Username     string
    Password     string
    Status       string
}


//Entity Untuk Lemdiklat

type AdminPusat struct{
	IdAdminPusat uint `gorm:"primary_key;auto_increment"`
	Nama string
	Email string
	Password string
	NoTelpon string
	Nip string	
	Status string
}


type Lemdik struct {

	IdLemdik uint 
	NamaLemdik string
	Status string 
}
