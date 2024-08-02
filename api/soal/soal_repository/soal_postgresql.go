package soal_repository

import (
	"bank_soal/models"
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	queryInsertSoal = `insert into soals (category_id,soal,jawaban_a,jawaban_b,jawaban_c,jawaban_d,jawaban_benar,created_at)
				values($1,$2,$3,$4,$5,$6,$7,$8) returning soal_id`

	queryGetSoal = `
	SELECT 
    soals.soal_id, 
    soals.category_id, 
    categories.category AS category_name,
    soals.soal, 
    soals.jawaban_a, 
    soals.jawaban_b, 
    soals.jawaban_c, 
    soals.jawaban_d, 
    soals.jawaban_benar, 
    soals.created_at
FROM soals
JOIN categories ON soals.category_id = categories.category_id
WHERE soals.deleted_at IS NULL

		
	`

	queryCountSoal = `
	select count(soal_id) from soals JOIN 
    categories 
ON 
    soals.category_id = categories.category_id 
WHERE 
    soals.deleted_at IS NULL
		
	`

	queryUpdateSoal = `
		UPDATE soals 
		SET 
			category_id = $1,
			soal = $2,
			jawaban_a = $3,
			jawaban_b = $4,
			jawaban_c = $5,
			jawaban_d = $6,
			jawaban_benar = $7,
			updated_at = $8
		WHERE soal_id = $9 AND deleted_at IS NULL
	`

	queryDeleteSoal = `
		UPDATE soals
		SET	
			deleted_at =$1
		where soal_id = $2
	`

	queryGetSoalById = `
		select soal_id, category_id,soal,jawaban_a,jawaban_b,jawaban_c,jawaban_d,jawaban_benar,created_at from soals
		 where soal_id =$1 and
		deleted_at is null
	`
)

type SoalRepositoryInterface interface {
	CreateSoal(ctx context.Context, soal models.Soals) (ID int64, err error)
	GetSoal(ctx context.Context, searchCriteria map[string]interface{}, page int, limit int) (resp []models.Soals, err error)
	CountSoal(ctx context.Context, params map[string]interface{}) (count int64, err error)
	UpdateSoal(ctx context.Context, soal models.Soals) error
	DeleteSoal(ctx context.Context, ID int64) error
	GetSoalById(ctx context.Context, id int64) (result models.Soals, err error)
}

type SoalRepositoryImpl struct {
	db *sqlx.DB
}

func NewSoalRepository(db *sqlx.DB) SoalRepositoryInterface {
	return &SoalRepositoryImpl{db: db}
}

func (r *SoalRepositoryImpl) CreateSoal(ctx context.Context, soal models.Soals) (ID int64, err error) {
	created_at := time.Now()
	err = r.db.QueryRowContext(ctx, queryInsertSoal, soal.CategoryId, soal.Soal, soal.JawabanA, soal.JawabanB, soal.JawabanC, soal.JawabanD, soal.JawabanBenar, created_at).Scan(&ID)
	if err != nil {
		err = fmt.Errorf("queryInsertSoal err%+v", err)
		return
	}

	return ID, nil
}
func (r *SoalRepositoryImpl) GetSoal(ctx context.Context, searchCriteria map[string]interface{}, page int, limit int) (soal []models.Soals, err error) {
	if limit > 10 {
		limit = 10
	}
	offset := (page - 1) * limit

	limitString := strconv.Itoa(limit)
	offsetString := strconv.Itoa(offset)

	// Ensure there's a space before LIMIT
	sqlQuery := queryGetSoal + searchCriteria["custom_query"].(string) + " LIMIT " + limitString + " OFFSET " + offsetString

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
	var ct models.Category

	for rows.Next() {
		var s models.Soals
		err = rows.Scan(&s.SoalId, &s.CategoryId, &ct.Category, &s.Soal, &s.JawabanA, &s.JawabanB, &s.JawabanC, &s.JawabanD, &s.JawabanBenar, &s.CreatedAt)
		if err != nil {
			err = fmt.Errorf("row scan err: %+v", err)
			return nil, err
		}
		soal = append(soal, s)
	}

	if err = rows.Err(); err != nil {
		err = fmt.Errorf("rows iteration err: %+v", err)
		return nil, err
	}

	return soal, nil
}

func (r *SoalRepositoryImpl) UpdateSoal(ctx context.Context, soal models.Soals) error {
	updated_at := time.Now()
	_, err := r.db.ExecContext(ctx, queryUpdateSoal, soal.CategoryId, soal.Soal, soal.JawabanA, soal.JawabanB, soal.JawabanC, soal.JawabanD, soal.JawabanBenar, updated_at, soal.SoalId)
	if err != nil {
		err = fmt.Errorf("UpdateSoal err%+v", err)
		return err
	}

	return nil
}

func (r *SoalRepositoryImpl) DeleteSoal(ctx context.Context, ID int64) error {
	deleted_at := time.Now()
	_, err := r.db.ExecContext(ctx, queryDeleteSoal, deleted_at, ID)
	if err != nil {
		return fmt.Errorf("queryDeleteSoal gagal %+v", err)
	}
	return nil
}

func (r *SoalRepositoryImpl) CountSoal(ctx context.Context, params map[string]interface{}) (count int64, err error) {
	sqlQuery := queryCountSoal + params["custom_query"].(string)
	err = r.db.QueryRowContext(ctx, sqlQuery).Scan(&count)

	if err != nil {
		err = fmt.Errorf("queryCountTotalKlaim error: %+v", err)
		return
	}

	return count, nil
}

func (r *SoalRepositoryImpl) GetSoalById(ctx context.Context, id int64) (soal models.Soals, err error) {

	err = r.db.QueryRowContext(ctx, queryGetSoalById, id).Scan(
		&soal.SoalId,
		&soal.CategoryId,
		&soal.Soal,
		&soal.JawabanA,
		&soal.JawabanB,
		&soal.JawabanC,
		&soal.JawabanD,
		&soal.JawabanBenar,
		&soal.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Soals{}, fmt.Errorf("soal not found")
		}
		return models.Soals{}, fmt.Errorf("error querying soal: %w", err)
	}

	return soal, nil
}
