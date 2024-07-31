package repository

import (
	"bank_soal/models"
	"bank_soal/utils/healper"
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	queryInsertUser = `
		INSERT INTO users (username,password,email,role,created_at,created_by)VALUES ($1, $2, $3, $4, $5,$6)
		RETURNING user_id;
	`
	queryGetUser = `		
		SELECT username , password ,email  
		FROM users 
		WHERE (username = $1 OR email = $2) and password=$3 limit 1
	`

	queryUpdateUser = `
		UPDATE users
		SET 
			username = $1,
			email=$2,
			password=$3,
			updated_at =$5
		where user_id =$4 and deleted_at is null
	`

	queryGetAllUser = `
		select user_id,username,password,email,role,created_at from users where deleted_at is null
	`

	queryCountUser = `
		select count(user_id) from users where deleted_at is null
	`
	queryCekRole = `SELECT role FROM users WHERE id = $1 LIMIT 1`

	queryGetUserById = `
	SELECT id, name, role FROM users WHERE id = $1
	`
)

type UserRepositoryInterface interface {
	CreateUser(ctx context.Context, user models.Users) (ID int64, err error)
	LoginUser(ctx context.Context, usernameOrEmail, password string) (models.UsersRespon, error)
	UpdateUser(ctx context.Context, user models.Users) error
	GetAllUser(ctx context.Context, searchCriteria map[string]interface{}, page int, limit int) (user []models.Users, err error)
	CountUser(ctx context.Context, params map[string]interface{}) (count int64, err error)
	GetUserRole(ctx context.Context, userID int64) (int, error)
	GetUserByID(ctx context.Context, userID int) (models.Users, error)
}
type UserRepositoryInterfaceImpl struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepositoryInterface {
	return &UserRepositoryInterfaceImpl{db: db}
}

// CreateUser repository
func (r *UserRepositoryInterfaceImpl) CreateUser(ctx context.Context, user models.Users) (ID int64, err error) {
	createdAt := time.Now()
	createdBy := healper.GetCreatedByFromCtx(ctx)

	err = r.db.QueryRowContext(ctx, queryInsertUser, user.Username, user.Password, user.Email, user.Role, createdAt, createdBy).Scan(&ID)
	if err != nil {
		err = fmt.Errorf("queryInsertUser err %+v", err)
		return
	}

	return ID, nil
}

func (r *UserRepositoryInterfaceImpl) LoginUser(ctx context.Context, usernameOrEmail, password string) (models.UsersRespon, error) {
	var user models.UsersRespon
	err := r.db.QueryRowContext(ctx, queryGetUser, usernameOrEmail, usernameOrEmail, password).Scan(&user.Username, &user.Email, &user.Password)
	if err != nil {
		return models.UsersRespon{}, fmt.Errorf("error logging in user: %v", err)
	}
	fmt.Println(queryGetUser)
	fmt.Println("usernameOrEmail", usernameOrEmail)
	fmt.Println("password", password)
	return user, nil
}

func (r *UserRepositoryInterfaceImpl) UpdateUser(ctx context.Context, user models.Users) error {
	updatedAd := time.Now()
	_, err := r.db.ExecContext(ctx, queryUpdateUser, user.Username, user.Email, user.Password, user.UserId, updatedAd)
	if err != nil {
		return fmt.Errorf("queryUpdateUser err %+v", err)
	}

	return nil
}

func (r *UserRepositoryInterfaceImpl) GetAllUser(ctx context.Context, searchCriteria map[string]interface{}, page int, limit int) (user []models.Users, err error) {
	if limit > 10 {
		limit = 10
	}
	offset := (page - 1) * limit

	limitString := strconv.Itoa(limit)
	offsetString := strconv.Itoa(offset)

	sqlQuery := queryGetAllUser + searchCriteria["custom_query"].(string) + " LIMIT " + limitString + " OFFSET " + offsetString

	rows, err := r.db.QueryContext(ctx, sqlQuery)
	if err != nil {
		if err != sql.ErrNoRows {
			err = fmt.Errorf("queryGetSoal err: %+v", err)
			return
		}
		err = nil
		return
	}

	defer rows.Close()
	for rows.Next() {
		var u models.Users
		err = rows.Scan(&u.UserId, &u.Username, &u.Password, &u.Email, &u.Role, &u.CreatedAt)
		if err != nil {
			err = fmt.Errorf("rows scan err : %+v", err)
			return nil, err
		}
		user = append(user, u)
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("rows iteration err: %+v", err)
		return nil, err
	}

	return user, nil
}

func (r *UserRepositoryInterfaceImpl) CountUser(ctx context.Context, params map[string]interface{}) (count int64, err error) {
	sqlQuery := queryCountUser + params["custom_query"].(string)
	err = r.db.QueryRowContext(ctx, sqlQuery).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *UserRepositoryInterfaceImpl) GetUserRole(ctx context.Context, userID int64) (int, error) {
	var role int

	err := r.db.QueryRowContext(ctx, queryCekRole, userID).Scan(&role)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("user tidak ditemukan")
		}
		return 0, fmt.Errorf("gagal mendapatkan role pengguna: %v", err)
	}
	return role, nil
}

func (r *UserRepositoryInterfaceImpl) GetUserByID(ctx context.Context, userID int) (models.Users, error) {
	var user models.Users

	err := r.db.QueryRowContext(ctx, queryGetUserById, userID).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("user not found: %w", err)
		}
		return user, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}
