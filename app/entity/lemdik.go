package entity


//Lemdik disini adalah si Balai Pelatihan, Politeknik, Sekolah Smk 


type Lemdiklat struct{
	IdLemdik uint 	`gorm:"primary_key;auto_increment"`
	NamaLemdik string
	NoTelpon int 
	Email 	string
	Password string
	Alamat 	string
	CreateAt string
	UpdateAt string
	
}
