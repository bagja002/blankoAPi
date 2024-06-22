package entity

type AdminBalai struct {
	IdAdminBalai uint `gorm:"primaryKey;autoIncrement"`
	NamaBalai    string
	Username     string
	Email        string
	Password     string
	No_telpon    int
	Type         string
	Provinsi     string
	Kota         string
	Status       string
	Is_active    string
}

type User struct {
	IdUser   uint `gorm:"primaryKey;autoIncrement"`
	Nama     string
	Email    string
	Password string
	Status   string
}

// Tipe seperti Ankapin, dan Atkapin atau lainnya CPIB CHCP
type TypeUjian struct {
	IdTypeUjian   uint
	NamaTypeUjian string
	CreateAt      string
	UpdateAt      string
}

// Nama Dari Uian Seperti Ankapin 1,2,3 ATKAPIN 1,2,3
type Ujian struct {
	IdUjian      uint `gorm:"primaryKey;autoIncrement"`
	IdTypeUjian  uint
	IdAdminBalai uint
	NameUjian    string
	Durasi       string
	Angkatan     string
	Keterangan   string
	Status       string
}

type Kompotensi struct {
	IdKompotensi uint `gorm:"primaryKey;autoIncrement"`
	IdUjian      uint

	NamaKomptensi string
	Keterangan    string
	Is_active     string
	Is_update     string
}

type SubKompotensi struct {
	IdSubKompotensi uint `gorm:"primaryKey;autoIncrement"`
	IdUjian         uint
	IdKompotensi    uint
	NamaKompotensi  string
	Keterangan      string
	Is_active       string
	Is_update       string
}

type Pertanyaan struct {
	IdPertanyaan    uint `gorm:"primaryKey;autoIncrement"`
	IdUjian         uint
	IdKompotensi    uint
	IdSubKompotensi uint
	JawabanBenar    string
	NamaSoal        string
	Status          string
}
