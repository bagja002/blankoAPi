package entity


//TUK sendiri diberikan oleh si LSP coy


type Tuk struct{
	IdTuk uint  `gorm:"primary_key;auto_increment"`
	IdLsp uint 
	IdLemdik uint 
	NamaTuk string 
	NoTelpon int 
	Email string
	Pasword string

}