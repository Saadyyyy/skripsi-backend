package models

type Rangking struct {
	RangkingId int64 `gorm:"primaryKey;autoIncrement:true"`
	UserId     int64
	CategoryId int64
	SoalId     int64
	Point      int64
	Next       bool
	CreatedAt  string
	UpdatedAt  string
	DeletedAt  string
}
