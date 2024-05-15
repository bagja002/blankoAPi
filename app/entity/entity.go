package entity

//Users Membuat Akun sendiri 
type Users struct{
	IdUsers uint 		`gorm:"primary_key;auto_increment"`
	Nama	string
	NoTelpon int  
	Email string
	Password string
	Kota string
	Provinsi string
	Alamat string
	Nik int 
	TempatLahir string
	TanggalLahir string
	JenisKelamin string
	Pekerjaan string
	GolonganDarah string
	StatusMenikah string
	Kewarganegaraan string
	IbuKandung string
	NegaraTujuanBekerja string
	PendidikanTerakhir string
	Agama string
	Foto string
	Status 	string
	CreateAt string
	UpdateAt string

	UsersPelatihan []UsersPelatihan `gorm:"foreignKey:IdUsers"`
}

//Auto generate
type SuperAdmin struct {
    IdSuperAdmin uint `gorm:"primary_key;auto_increment"`
    Nama         string
    Email        string
    Username     string
    Password     string
    Status       string
}


//Entity Untuk Lemdiklat
//Auto Generate or Input From admin pusat
type AdminPusat struct{
	IdAdminPusat uint `gorm:"primary_key;auto_increment"`
	Nama string
	Email string
	Password string
	NoTelpon string
	Nip string	
	Status string
}

