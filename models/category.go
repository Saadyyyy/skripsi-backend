package models

type Category struct {
	CategoryId int64 `gorm:"primaryKey;autoIncrement:true"`
	Category   string
	CreatedAt  string
	UpdatedAt  string
	DeletedAt  string
}

type FilterCategory struct {
	Page       int
	Limit      int
	TglMulai   string
	TglSelesai string
	Keyword    string
}
