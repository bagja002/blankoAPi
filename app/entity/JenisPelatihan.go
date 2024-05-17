package entity



type JenisBidangPelatihan struct{
	IdJenisBidangPelatihan uint   `gorm:"primary_key;auto_increment"`
	KodeBidang string
	NamaBidang string 
	CreateAt string
	UpdateAt string
}

type  JenisBidangKompotensi struct{
	IdJenisBidangPelatihan uint   `gorm:"primary_key;auto_increment"`
	KodeKompotensi string
	NamaKompotensi string 
	CreateAt string
	UpdateAt string
}

