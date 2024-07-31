package models

type Users struct {
	UserId    int64 `gorm:"primaryKey;autoIncrement:true"`
	Username  string
	Password  string
	Email     string
	Role      int64
	Token     string
	CreatedAt string
	CreatedBy string
	UpdatedAt string
	UpdatedBy string
	DeletedAt string
	DeletedBy string
}

type UsersRespon struct {
	UserId   string
	Username string
	Password string
	Email    string
	Role     int64
}

type FilterUser struct {
	Page       int
	Limit      int
	TglMulai   string
	TglSelesai string
	Keyword    string
	Category   int64
}

// Buat struct untuk permintaan perubahan kata sandi
type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required"`
	ConfirmPassword string `json:"confirmPassword" binding:"required"` // Tambahkan ConfirmPassword
}
