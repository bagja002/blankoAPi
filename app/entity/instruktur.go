package entity



type Instruktur struct{
	IdInstruktur uint `gorm:"primary_key;auto_increment"`
	Nama string
	IdLemdik uint 
	Status string
	PendidikkanTerakhir string

}