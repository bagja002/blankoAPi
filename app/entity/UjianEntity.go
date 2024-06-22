package entity

type SoalUjianLemdik struct {
	IdSoalUjian  uint `gorm:"primaryKey;autoIncrement"`
	IdLemdik     uint
	IdPelatihan  uint
	Soal         string
	JawabanBenar string
	Status       string
	CreateAt     string
	UpdateAt     string
	Jawaban      []Jawaban `gorm:"foreignKey:IdSoalUjian"`
}

type Jawaban struct {
	IdJawaban   uint `gorm:"primaryKey;autoIncrement"`
	IdSoalUjian uint
	NameJawaban string
	Status      string
	CreateAt    string
	UpdateAt    string
}

type UsersSoal struct {
	IdUsersSoal     uint `gorm:"primaryKey;autoIncrement"`
	IdUserPelatihan uint

	IdSoalUjian uint
	CodeAksess  string
	CreateAt    string
	UpdateAt    string
}
