package entity

type Admin struct {
	IdAdmin  uint `gorm:"primary_key;auto_increment"`
	Nama     string
	NoTelpon int
	Username string
	Password string
	CreateAt string
	UpdateAt string
}

type SuperAdmin struct {
	IdSuperAdmin uint `gorm:"primary_key;auto_increment"`
	Nama         string
	Email        string
	Username     string
	Password     string
	Status       string
}

type Blanko struct {
	IdBlanko         int `gorm:"primary_key;auto_increment"`
	Jumlah           int // Akan berubah berdasarkan table `blanko_keluar`
	NoSeri           string
	TipeBlanko       string // CoC atau CoP
	TanggalPengadaan string
	JumlahPengadaan  int
	CreateAt         string
	UpdateAt         string
}

type BlankoKeluar struct {
	IdBlankoKeluar        int    `gorm:"primary_key;auto_increment"`
	IdBlanko              int    // Relasi ke `blanko`
	TipeBlanko            string // CoC atau CoP
	TanggalKeluar         string
	NamaLemdiklat         string
	NamaPelaksana         string
	TanggalPermohonan     string
	LinkPermohonan        string
	NamaProgram           string
	TanggalPelaksanaan    string
	JumlahPesertaLulus    int
	JumlahBlankoDiajukan  int
	JumlahBlankoDisetujui int
	NoSeriBlanko          string
	Status                string
	IsSd                  bool
	IsCetak               bool
	TipePengambilan       string
	PetugasYangMenerima   string
	PetugasYangMemberi    string
	LinkDataDukung        string
	CreatedAt             string
	UpdatedAt             string
	Keterangan            string
	AsalPendapatan        string
}

type BlankoRusak struct {
	IdBlankoRusak  int `gorm:"primary_key;auto_increment"`
	IdBlankoKeluar int // Relasi ke `blanko`
	NoSeri         string
	Tipe           string // rusak atau hilang
	Keterangan     string
	TanggalRusak   string
	FotoDokumen    string
	CreatedAt      string
	UpdateAt       string
}
