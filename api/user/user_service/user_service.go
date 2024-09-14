package service

import (
	postgresql "bank_soal/api/user/user_repository"
	"bank_soal/middleware"
	"bank_soal/models"
	"bank_soal/utils/healper"
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type UserService interface {
	CreateUser(ctx context.Context, user models.Users) (ID int64, err error)
	LoginUser(ctx context.Context, usernameOrEmail, password string) (models.UsersRespon, string, error)
	UpdateUser(ctx context.Context, user models.Users) error
	GetAllUser(ctx context.Context, filter models.FilterUser) (user []models.Users, totalData int64, err error)
	GetUserByID(ctx context.Context, userId int64) (user models.Users, err error)
	UpdateUserRoleByID(ctx context.Context, user models.Users) error
}

type UserServiceImpl struct {
	repo postgresql.UserRepositoryInterface
	db   *sqlx.DB
}

func NewUserService(repo postgresql.UserRepositoryInterface, db *sqlx.DB) UserService {
	return &UserServiceImpl{repo: repo, db: db}
}

func (u *UserServiceImpl) CreateUser(ctx context.Context, user models.Users) (ID int64, err error) {
	// Validasi kosong
	if user.Username == "" || user.Password == "" || user.Email == "" {
		return 0, fmt.Errorf("tidak boleh kosong harus diisi")
	}

	// Cek username yang sudah ada
	if user.Username != "" {
		var existingUsername string
		queryCekUsername := `SELECT username FROM users WHERE username = $1 LIMIT 1`
		err = u.db.QueryRowContext(ctx, queryCekUsername, user.Username).Scan(&existingUsername)
		if err == nil {
			return 0, fmt.Errorf("username sudah digunakan")
		} else if err != sql.ErrNoRows {
			return 0, fmt.Errorf("gagal memeriksa username: %v", err)
		}
	}

	// Cek email yang sudah ada
	if user.Email != "" {
		var existingEmail string
		queryCekEmail := `SELECT email FROM users WHERE email = $1 LIMIT 1`
		err = u.db.QueryRowContext(ctx, queryCekEmail, user.Email).Scan(&existingEmail)
		if err == nil {
			return 0, fmt.Errorf("email sudah digunakan")
		} else if err != sql.ErrNoRows {
			return 0, fmt.Errorf("gagal memeriksa email: %v", err)
		}
	}

	// Pilih gambar acak
	gambar, err := pilihGambarAcak()
	if err != nil {
		return 0, fmt.Errorf("gagal mendapatkan gambar: %v", err)
	}
	baseURL, exists := os.LookupEnv("BASEURL")
	if !exists {
		return 0, fmt.Errorf("variabel lingkungan BASEURL tidak ditemukan")
	}
	user.Profile = baseURL + gambar

	// Hash password
	hashedPassword, err := healper.HashPassword(user.Password)
	if err != nil {
		return 0, fmt.Errorf("gagal hash password")
	}
	user.Password = hashedPassword

	// Simpan pengguna ke database
	ID, err = u.repo.CreateUser(ctx, user)
	if err != nil {
		return 0, fmt.Errorf("gagal membuat pengguna: %v", err)
	}

	return ID, nil
}

func pilihGambarAcak() (string, error) {

	gambarList := []string{"foto1.png", "foto2.png", "foto3.png", "foto4.png"}
	rand.Seed(time.Now().UnixNano())
	gambarTerpilih := gambarList[rand.Intn(len(gambarList))]

	// Path ke file gambar di folder assets
	path := filepath.Join("assets", gambarTerpilih)

	// Pastikan file gambar dapat diakses
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", fmt.Errorf("file gambar tidak ditemukan: %s", path)
	}

	return path, nil
}

func (u *UserServiceImpl) LoginUser(ctx context.Context, usernameOrEmail, password string) (models.UsersRespon, string, error) {
	user, err := u.repo.LoginUser(ctx, usernameOrEmail, password)
	if err != nil {
		return models.UsersRespon{}, "", fmt.Errorf("gagal get login user repository: %+v", err)
	}

	comparePass := healper.CompareHash(user.Password, password)
	if !comparePass {
		return models.UsersRespon{}, "", fmt.Errorf("gagal compare pass")
	}
	users := models.Users{}

	// Create JWT token
	token, err2 := middleware.CreateToken(users.Username, user.Role)
	if err2 != nil {
		return models.UsersRespon{}, "", fmt.Errorf("gagal create token: %+v", err)
	}
	return user, token, nil
}

func (u *UserServiceImpl) UpdateUser(ctx context.Context, user models.Users) error {
	hashedPassword, err := healper.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("gagal hash password")
	}
	user.Password = hashedPassword

	err2 := u.repo.UpdateUser(ctx, user)
	if err2 != nil {
		return fmt.Errorf("gagal get repo update user %+v", err)
	}

	return nil
}

func (s *UserServiceImpl) GetAllUser(ctx context.Context, filter models.FilterUser) (user []models.Users, totalData int64, err error) {
	params := map[string]interface{}{
		"deleted_at":   nil,
		"custom_query": "",
	}

	if filter.Keyword != "" {
		keywordLower := strings.ToLower(filter.Keyword)
		params["custom_query"] = fmt.Sprintf("%s AND LOWER(username) LIKE '%%%s%%'", params["custom_query"], keywordLower)
	}

	// if filter.Category != 0 {
	// 	categoryString := strconv.Itoa(int(filter.Category))
	// 	params["custom_query"] = fmt.Sprintf("%s AND category_id='%s'", params["custom_query"], categoryString)
	// }

	user, err = s.repo.GetAllUser(ctx, params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get soal from repository: %+v", err)
	}

	totalData, err = s.repo.CountUser(ctx, params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count soal from repository: %+v", err)
	}

	resp := make([]models.Users, len(user))
	for i, s := range user {
		resp[i] = models.Users{
			UserId:    s.UserId,
			Username:  s.Username,
			Password:  s.Password,
			Email:     s.Email,
			Role:      s.Role,
			CreatedAt: s.CreatedAt,
			Profile:   s.Profile,
		}
	}

	return resp, totalData, nil
}

func (s *UserServiceImpl) GetUserByID(ctx context.Context, userId int64) (user models.Users, err error) {
	result, err := s.repo.GetUserByID(ctx, userId)
	if err != nil {
		return models.Users{}, fmt.Errorf("gagal get function dari repository ")
	}
	fmt.Println("result", result)
	return result, nil
}

func (u *UserServiceImpl) UpdateUserRoleByID(ctx context.Context, user models.Users) error {
	err2 := u.repo.UpdateUserRoleByID(ctx, user)
	if err2 != nil {
		return fmt.Errorf("gagal get repo update user %+v", err2)
	}
	return nil
}
