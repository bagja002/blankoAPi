package entity

type Sertifikat struct {
	IdSertifikat uint `gorm:"primary_key;auto_increment"`
	IdLemdik     uint
	IdPelatihan  uint
	NoSertfikat  string

	CreateAt string
	UpdateAt string
}
