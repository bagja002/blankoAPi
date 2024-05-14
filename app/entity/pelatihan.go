package entity


type Pelatihan struct{
	IdPelatihan uint `gorm:"primary_key;auto_increment"`
	IdLemdik uint 
	KodePelatihan string
	NamaPelatihan string	//Judul Pelatihan
	PenyelenggaraPelatihan string //Penyengeggara oelatihan
	DetailPelatihan string //Deskripsi Pelatihan 
	FotoPelatihan string
	JenisPelatihan string //Aspirasi, PNBP, Reguler
	BidangPelatihan string //Bidang Pelatihan 
	DukunganProgramTerobosan string //PIT, Non terobosan 
	TanggalMulaiPelatihan string
	TanggalBerakhirPelatihan string
	HargaPelatihan int    //Harga Pelatihan
	Instruktur string 
	Status string	//Aktif Atau Tidak 	
	MemoPusat string	//memo persetujuan ya g dikeluarkan oleh bu kapus melalui persuratan 
	SilabusPelatihan string  //Dsilabus Pelatihan dalam Bentuk File 
	LokasiPelatihan string  //Lokasi Pelatihan 
	PelaksanaanPelatihan string  //Pelaksana Pelatihan 
	UjiKompotensi string  //True Or False
	KoutaPelatihan string
	AsalPelatihan string  //Masyarakat Pelatihan

	//SECTION SERTIFIKAT 
	AsalSertifikat string  //JDPT/BPSDM 
	JenisSertifikat string  //teknis, kepelautan , umum
	TtdSertifikat string  //Pilih Penandatangan
	NoSertifikat string //Nomor Sertifikat Perpelatihan

	//Penambahan Paket Penginapan 
	IdSaranaPrasarana string
	IdKonsumsi string

	CreateAt string 
	UpdateAt string
	SarprasPelatihan []SarprasPelatihan `gorm:"foreignKey:IdPelatihan"`
}