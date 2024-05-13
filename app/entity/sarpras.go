package entity


type Sarpras struct{
	IdSarpras uint  `gorm:"primary_key;auto_increment"`
	IdLemdik uint 
	NamaSarpras string
	Harga int 
	Deskripsi string
	Jenis string
	CreateAt string
	UpdateAt string
}
//Perlemdik
type SarprasPelatihan struct{
	IdSarprasPelatihan uint  `gorm:"primary_key;auto_increment"`
	IdPelatihan uint 
	IdLemdik uint 
	IdSarpras uint 
	
}