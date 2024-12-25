package models

type UserSoal struct {
	UserSoalID int64 `gorm:"primaryKey;autoIncrement:true"`
	UserID     int64
	CategoryID int64
	Answer     bool `gorm:"default:false"`
}
