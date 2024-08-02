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
		select category_id, category from from categories where category_id = $1 and deleted is null
	`
)

type CategoryRepository interface {
	CreateCategory(ctx context.Context, ct models.Category) (id int64, err error)
	GetCategoryByID(ctx context.Context, id int64) (ct models.Category, err error)
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
		&ct.CategoryId, &ct.Category,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Category{}, err
		}
		return models.Category{}, err
	}

	return ct, nil
}
