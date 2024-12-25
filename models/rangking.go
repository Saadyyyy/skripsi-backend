package models

type Rangking struct {
	RankId     int64 `gorm:"primaryKey;autoIncrement:true"`
	UserId     int64
	CategoryId int64
	Point      int64
	CreatedAt  string
	UpdatedAt  string
	DeletedAt  string
}

type RangkingRes struct {
	UserId   int64  `json:"user_id"`
	Point    int64  `json:"point"`
	Username string `json:"username"`
	Profile  string `json:"profile"`
}
