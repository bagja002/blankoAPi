package entity

type Sarpras struct {
	IdSarpras   uint `gorm:"primary_key;auto_increment"`
	IdLemdik    uint
	NamaSarpras string
	Harga       int
	Deskripsi   string
	Jenis       string
	CreateAt    string
	UpdateAt    string
	FotoSarpras string

	Sarprasatih []SarprasPelatihan `gorm:"foreignKey:IdSarpras"`
}

// Perlemdik
type SarprasPelatihan struct {
	IdSarprasPelatihan uint `gorm:"primary_key;auto_increment"`
	IdPelatihan        uint
	IdLemdik           uint
	IdSarpras          uint
	CreteAt            string
	UpdateAt           string
}

// Table Prasarana User Pelatihan
type UsersSarprasPelatihan struct {
	IdUsersSarprasPelatihan uint `gorm:"primary_key;auto_increment"`
	IdPelatihan             uint
	IdIdUserPelatihan       uint
	IdSarpras               uint
	CreteAt                 string
	UpdateAt                string
}
