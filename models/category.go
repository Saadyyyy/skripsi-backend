package models

type Category struct {
	CategoryId int64 `gorm:"primaryKey;autoIncrement:true"`
	Category   string
	CreatedAt  string
	UpdatedAt  string
	DeletedAt  string
}
