package models

type Soals struct {
	SoalId       int64 `gorm:"primaryKey;autoIncrement:true"`
	CategoryId   int64
	Soal         string
	JawabanA     string
	JawabanB     string
	JawabanC     string
	JawabanD     string
	JawabanBenar string
	CreatedAt    string
	UpdatedAt    string
	DeletedAt    string
}

type FilterSoal struct {
	Page       int
	Limit      int
	TglMulai   string
	TglSelesai string
	Keyword    string
	Category   int64
}
