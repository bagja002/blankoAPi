package entity


type Sertifikasi struct{
	IdSertifikasi uint 
	IdLemdik uint 		//Nanti ambil nama lemdiknya
	NamaSertifikasi string   //nama dari sertifikasinya
	BidangSertifikasi string  // Budidaya, Penangkapan Dll
	IdLsp string
	TempatTuk string
	TanggalUjianKompotensi string
	TanggalAkhir string
	Harga string
	CreateAt string
	UpdaetAt string

}