package entity


//Entiti LSP data yang di buat oleh si lembaga diklat
type Lsp struct{
	IdLsp uint   `gorm:"primary_key;auto_increment"`
	IdLemdik uint 
	NamaLsp string
	NoTelpon int 
	Email string
	Password string
	AlamatLsp string
	CreateAt string
	UpdateAt string
}