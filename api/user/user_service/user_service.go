package service

import (
	postgresql "bank_soal/api/user/user_repository"
	"bank_soal/middleware"
	"bank_soal/models"
	"bank_soal/utils/healper"
	"context"
	"database/sql"
	"fmt"
	"math"
	"strings"

	"github.com/jmoiron/sqlx"
)

type UserService interface {
	CreateUser(ctx context.Context, user models.Users) (ID int64, err error)
	LoginUser(ctx context.Context, usernameOrEmail, password string) (models.UsersRespon, string, error)
	UpdateUser(ctx context.Context, user models.Users) error
	GetAllUser(ctx context.Context, filter models.FilterUser) (user []models.Users, totalPage int64, totalData int64, err error)
}

type UserServiceImpl struct {
	repo postgresql.UserRepositoryInterface
	db   *sqlx.DB
}

func NewUserService(repo postgresql.UserRepositoryInterface, db *sqlx.DB) UserService {
	return &UserServiceImpl{repo: repo, db: db}
}

// CreateUser service
func (u *UserServiceImpl) CreateUser(ctx context.Context, user models.Users) (ID int64, err error) {
	//validasi empty
	if user.Username == "" || user.Password == "" || user.Email == "" {
		return 0, fmt.Errorf("tidak boleh kosong harus di isi")
	}

	if user.Username != "" {
		var existingUsername string
		queryCekUsername := `select username from users where username = $1 limit 1`
		err = u.db.QueryRowContext(ctx, queryCekUsername, user.Username).Scan(&existingUsername)
		if err == nil {
			return 0, fmt.Errorf("username sudah digunakan")
		} else if err != sql.ErrNoRows {
			return 0, fmt.Errorf("gagal memeriksa username: %v", err)
		}
	}

	if user.Email != "" {
		var existingUsername string
		queryCekUsername := `select email from users where email = $1 limit 1`
		err = u.db.QueryRowContext(ctx, queryCekUsername, user.Username).Scan(&existingUsername)
		if err == nil {
			return 0, fmt.Errorf("email sudah digunakan")
		} else if err != sql.ErrNoRows {
			return 0, fmt.Errorf("gagal memeriksa email: %v", err)
		}
	}

	hashedPassword, err := healper.HashPassword(user.Password)
	if err != nil {
		return 0, fmt.Errorf("gagal hash password")
	}
	user.Password = hashedPassword

	//repository create user
	ID, err = u.repo.CreateUser(ctx, user)
	if err != nil {
		return 0, fmt.Errorf("gagal membuat pengguna: %v", err)
	}

	return ID, nil
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

func (s *UserServiceImpl) GetAllUser(ctx context.Context, filter models.FilterUser) (user []models.Users, totalPage int64, totalData int64, err error) {
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

	if filter.Page <= 0 {
		filter.Page = 1 // default page
	}

	if filter.Limit <= 0 || filter.Limit > 10 {
		filter.Limit = 10 // default limit
	}

	user, err = s.repo.GetAllUser(ctx, params, filter.Page, filter.Limit)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to get soal from repository: %+v", err)
	}

	totalData, err = s.repo.CountUser(ctx, params)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("failed to count soal from repository: %+v", err)
	}

	totalPage = int64(math.Ceil(float64(totalData) / float64(filter.Limit)))
	if totalData == 0 {
		totalPage = 1
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
		}
	}

	return resp, totalPage, totalData, nil
}
