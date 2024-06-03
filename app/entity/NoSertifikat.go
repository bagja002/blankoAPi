package entity

type NoSertfikat struct {
	IdNoSertfikat       uint `gorm:"primary_key;auto_increment"`
	IdLemdik            uint
	Nomor               string
	NoLengkapSertifikat string
	CreateAt            string
	UpdateAt            string
}
