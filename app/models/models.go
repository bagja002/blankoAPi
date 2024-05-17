package models

type Admin struct {
	Name      string `json:"name"`
	ID        string `json:"id"`
	Role      string `json:"role"`
	Satminkal string `json:"satminkal"`
}

// Moderator represents the moderator entity
type Moderator struct {
	Name      string `json:"name"`
	ID        string `json:"id"`
	Role      string `json:"role"`
	Satminkal string `json:"satminkal"`
}

// User represents the user entity

type SuperAdmin struct{ 
	IdSuperAdmin uint  `json:"id_admin"`
	Nama string	`json:"nama"`
	Email string`json:"email"`
	Username string`json:"username"`
	Password string`json:"password"`
	Status string`json:"status"`
}

//Entity Untuk Lemdiklat

type AdminPusat struct{
	IdAdminPusat uint 
	Nama string
	Email string
	Password string
	NoTelpon string
	Nip string	
	Status string
}

type Pelatihan struct {
	IdPelatihan              uint   `gorm:"primary_key;auto_increment" json:"id_pelatihan"`
	IdLemdik                 string `json:"IdLemdik"`
	KodePelatihan            string `json:"KodePelatihan"`
	NamaPelatihan            string `json:"NamaPelatihan"`
	PenyelenggaraPelatihan   string `json:"PenyelenggaraPelatihan"`
	DetailPelatihan          string `json:"DetailPelatihan"`
	JenisPelatihan           string `json:"JenisPelatihan"`
	BidangPelatihan          string `json:"BidangPelatihan"`
	DukunganProgramTerobosan string `json:"DukunganProgramTerobosan"`
	TanggalMulaiPelatihan    string `json:"TanggalMulaiPelatihan"`
	TanggalBerakhirPelatihan string `json:"TanggalBerakhirPelatihan"`
	HargaPelatihan           string `json:"HargaPelatihan"`
	Instruktur               string `json:"instruktur"`
	FotoPelatihan            string
	Status                   string `json:"status"`
	MemoPusat                string `json:"memo_pusat"`
	SilabusPelatihan         string `json:"silabus_pelatihan"`
	LokasiPelatihan          string `json:"lokasi_pelatihan"`
	PelaksanaanPelatihan     string `json:"pelaksanaan_pelatihan"`
	UjiKompotensi            string `json:"uji_kompetensi"`
	KoutaPelatihan           string `json:"kouta_pelatihan"`
	AsalPelatihan            string `json:"asal_pelatihan"`
	AsalSertifikat           string `json:"asal_sertifikat"`
	JenisSertifikat          string `json:"jenis_sertifikat"`
	TtdSertifikat            string `json:"ttd_sertifikat"`
	NoSertifikat             string `json:"no_sertifikat"`
	IdSaranaPrasarana        string `json:"id_sarana_prasarana"`

	
	
	IdKonsumsi               string `json:"id_konsumsi"`
	CreateAt                 string `json:"created_at"`
	UpdateAt                 string `json:"updated_at"`


	//Penambahan Matery 
	NamaMateri string `json:"NamaMateri"`
	Deskripsi string `json:"Deskripsi"`
	JamTeory string `json:"JamTeory"`
	JamPraktek string `json:"JamPraktek"`
}


type Users struct {
	IdUsers             uint 
	Nama                string
	NoTelpon            int
	Email               string
	Password            string
	Kota                string
	Provinsi            string
	Alamat              string
	Nik                 int
	TempatLahir         string
	TanggalLahir        string
	JenisKelamin        string
	Pekerjaan           string
	GolonganDarah       string
	StatusMenikah       string
	Kewarganegaraan     string
	IbuKandung          string
	NegaraTujuanBekerja string
	PendidikanTerakhir  string
	Agama               string
	Foto                string //Photo
	Ktp                 string //KTP
	KK                  string //Kartu Keluarga
	SuratKesehatan      string //SuratKesehatan
	Ijazah  string 
	Status              string
	CreateAt            string
	UpdateAt            string

}

