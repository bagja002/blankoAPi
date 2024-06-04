package entity

type NoSertfikat struct {
	IdNoSertfikat       uint `gorm:"primary_key;auto_increment"`
	IdPelatihan         uint
	NamaLemdik          string
	Nomor               int
	NoLengkapSertifikat string
	CreateAt            string
	UpdateAt            string
}

type NoRegistrasi struct {
	IdNoRegistrasi    uint `gorm:"primary_key;auto_increment"`
	IdPelatihan       uint
	IdUsers           uint
	NamaLemdik        string
	Bidang            string
	Nomor             int
	NoLengkapRegister string
	CreateAt          string
	UpdateAt          string
}
