package entity

type UsersPelatihan struct {
	IdUserPelatihan    uint `gorm:"primary_key;auto_increment"`
	IdUsers            uint
	Nama               string
	TempatTanggalLahir string
	IdPelatihan        uint
	NoSertifikat       string
	NoRegistrasi       string
	PreTest            int
	PostTest           int
	NilaiTeory         int
	NilaiPraktek       int

	//Nilai Materi
	StatusPembayaran string //Pending dan Void
	MetodoPembayaran string
	WaktuPembayaran  string
	Keterangan       string
	IsActice         string
	FileSertifikat   string
	Institusi        string
	TotalBayar       string
	CreteAt          string
	UpdateAt         string
}
