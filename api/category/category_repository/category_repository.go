package categoryrepository

import (
	"bank_soal/models"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"golang.org/x/net/context"
)

const (
	queryCreateCategory = `
		insert into categories(category,created_at)values($1,$2) returning category_id 
	`

	queryGetCategoryById = `
		select category_id, category,created_at from categories where category_id = $1 and deleted_at is null
	`

	queryGetAllCategory = `
		select category_id , category ,created_at from categories where deleted_at is null
	`

	queryCountCategory = `
		select count(category_id) from categories where deleted_at is null
	`

	queryUpdateCategory = `
		UPDATE categories SET category= $1 , updated_at= $2 where category_id= $3 and deleted_at is null
	`

	queryDeleteCategory = `
		UPDATE categories SET deleted_at =$1 where category_id= $2 and deleted_at is null
	`
)

type CategoryRepository interface {
	CreateCategory(ctx context.Context, ct models.Category) (id int64, err error)
	GetCategoryByID(ctx context.Context, id int64) (ct models.Category, err error)
	GetListCategory(ctx context.Context, searchCriteria map[string]interface{}) (ct []models.Category, err error)
	UpdateCategory(ctx context.Context, ct models.Category) (err error)
	CountUser(ctx context.Context, params map[string]interface{}) (count int64, err error)
	DeletedCategory(ctx context.Context, id int64) error
}

type CategoryRepositoryImpl struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) CategoryRepository {
	return &CategoryRepositoryImpl{db: db}
}

func (r *CategoryRepositoryImpl) CreateCategory(ctx context.Context, ct models.Category) (id int64, err error) {
	created_at := time.Now()
	err = r.db.QueryRowContext(ctx, queryCreateCategory, ct.Category, created_at).Scan(&id)
	if err != nil {
		err = fmt.Errorf("queryInsertSoal err%+v", err)
		return
	}

	return id, nil
}

func (r *CategoryRepositoryImpl) GetCategoryByID(ctx context.Context, id int64) (ct models.Category, err error) {
	err = r.db.QueryRowContext(ctx, queryGetCategoryById, id).Scan(
		&ct.CategoryId, &ct.Category, &ct.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Category{}, err
		}
		return models.Category{}, err
	}

	return ct, nil
}

func (r *CategoryRepositoryImpl) GetListCategory(ctx context.Context, searchCriteria map[string]interface{}) (ct []models.Category, err error) {

	// Ensure there's a space before LIMIT
	sqlQuery := queryGetAllCategory + searchCriteria["custom_query"].(string) + " LIMIT "

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

	var s models.Category
	for rows.Next() {
		err = rows.Scan(&s.CategoryId, &s.Category, &s.CreatedAt)
		if err != nil {
			err = fmt.Errorf("row scan err: %+v", err)
			return nil, err
		}
		ct = append(ct, s)
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("rows iteration err: %+v", err)
		return nil, err
	}

	return ct, nil
}

func (r *CategoryRepositoryImpl) CountUser(ctx context.Context, params map[string]interface{}) (count int64, err error) {
	sqlQuery := queryCountCategory + params["custom_query"].(string)
	err = r.db.QueryRowContext(ctx, sqlQuery).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *CategoryRepositoryImpl) UpdateCategory(ctx context.Context, ct models.Category) (err error) {
	updatedAt := time.Now()
	_, err = r.db.ExecContext(ctx, queryUpdateCategory, ct.Category, updatedAt, ct.CategoryId)
	if err != nil {
		return err
	}

	return nil
}

func (r *CategoryRepositoryImpl) DeletedCategory(ctx context.Context, id int64) error {
	deletedAt := time.Now()

	_, err := r.db.ExecContext(ctx, queryDeleteCategory, deletedAt, id)
	if err != nil {
		return fmt.Errorf("fail get query")
	}
	return nil
}
