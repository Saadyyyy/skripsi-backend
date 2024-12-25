package soal_repository

import (
	"bank_soal/models"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

const (
	queryInsertSoal = `insert into soal_exsamples (category_id,soal,jawaban_a,jawaban_b,jawaban_c,jawaban_d,jawaban_e,jawaban_benar,created_at)
				values($1,$2,$3,$4,$5,$6,$7,$8,$9) returning soal_id`

	queryGetSoal = `
	SELECT 
    soal_exsamples.soal_id, 
    soal_exsamples.category_id, 
    categories.category AS category_name,
    soal_exsamples.soal, 
    soal_exsamples.jawaban_a, 
    soal_exsamples.jawaban_b, 
    soal_exsamples.jawaban_c, 
    soal_exsamples.jawaban_d, 
	soal_exsamples.jawaban_e, 
    soal_exsamples.jawaban_benar, 
    soal_exsamples.created_at
FROM soal_exsamples
JOIN categories ON soal_exsamples.category_id = categories.category_id
WHERE soal_exsamples.deleted_at IS NULL  

		
	`

	queryCountSoal = `
	select count(soal_id) from soal_exsamples JOIN 
    categories 
ON 
    soal_exsamples.category_id = categories.category_id 
WHERE 
    soal_exsamples.deleted_at IS NULL 
		
	`

	queryUpdateSoal = `
		UPDATE soal_exsamples 
		SET 
			category_id = $1,
			soal = $2,
			jawaban_a = $3,
			jawaban_b = $4,
			jawaban_c = $5,
			jawaban_d = $6,
			jawaban_e = $7,
			jawaban_benar = $8,
			updated_at = $9
		WHERE soal_id = $10 AND deleted_at IS NULL
	`

	queryDeleteSoal = `
		UPDATE soal_exsamples
		SET	
			deleted_at =$1
		where soal_id = $2
	`

	queryGetSoalById = `
		select soal_id, category_id,soal,jawaban_a,jawaban_b,jawaban_c,jawaban_d,jawaban_e,jawaban_benar,created_at from soal_exsamples
		 where soal_id =$1 and
		deleted_at is null
	`
)

type SoalRepositoryInterface interface {
	CreateSoal(ctx context.Context, soal models.SoalExsample) (ID int64, err error)
	GetSoal(ctx context.Context, searchCriteria map[string]interface{}, rdb *redis.Client) (resp []models.SoalExsample, err error)
	CountSoal(ctx context.Context, params map[string]interface{}) (count int64, err error)
	UpdateSoal(ctx context.Context, soal models.SoalExsample) error
	DeleteSoal(ctx context.Context, ID int64) error
	GetSoalById(ctx context.Context, id int64) (result models.SoalExsample, err error)
}

type SoalRepositoryImpl struct {
	db *sqlx.DB
}

func NewSoalRepository(db *sqlx.DB) SoalRepositoryInterface {
	return &SoalRepositoryImpl{db: db}
}

func (r *SoalRepositoryImpl) CreateSoal(ctx context.Context, soal models.SoalExsample) (ID int64, err error) {

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		fmt.Printf("Failed to load location: %+v\n", err)
		return
	}

	// Cetak waktu sekarang dalam sistem timezone dan WIB
	fmt.Println("System Time:", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("Time in Asia/Jakarta:", time.Now().In(loc).Format("2006-01-02 15:04:05"))

	// Get current time in Indonesia timezone
	created_at := time.Now()

	err = r.db.QueryRowContext(ctx, queryInsertSoal, soal.CategoryId, soal.Soal, soal.JawabanA, soal.JawabanB, soal.JawabanC, soal.JawabanD, soal.JawabanE, soal.JawabanBenar, created_at).Scan(&ID)
	if err != nil {
		err = fmt.Errorf("queryInsertSoal err%+v", err)
		return
	}

	return ID, nil
}
func (r *SoalRepositoryImpl) GetSoal(ctx context.Context, searchCriteria map[string]interface{}, rdb *redis.Client) ([]models.SoalExsample, error) {
	// Check if cached result exists

	customQuery, ok := searchCriteria["custom_query"].(string)
	if !ok {
		customQuery = "" // Nilai default jika custom_query tidak ada
	}
	sqlQuery := queryGetSoal + customQuery

	cacheKey := fmt.Sprintf("soals:%v", sqlQuery)
	cachedData, err := rdb.Get(ctx, cacheKey).Result()

	if err == nil {
		var cachedSoals []models.SoalExsample
		err := json.Unmarshal([]byte(cachedData), &cachedSoals)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal cached data: %v", err)
		}
		return cachedSoals, nil
	}

	rows, err := r.db.QueryContext(ctx, sqlQuery)

	if err != nil {
		if err != sql.ErrNoRows {
			return nil, fmt.Errorf("queryGetSoal err: %+v", err)
		}
		return nil, nil
	}
	defer rows.Close()

	var soal []models.SoalExsample
	var ct models.Category

	for rows.Next() {
		var s models.SoalExsample
		err = rows.Scan(&s.SoalId, &s.CategoryId, &ct.Category, &s.Soal, &s.JawabanA, &s.JawabanB, &s.JawabanC, &s.JawabanD, &s.JawabanE, &s.JawabanBenar, &s.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("row scan err: %+v", err)
		}
		soal = append(soal, s)
	}

	dataToCache, err := json.Marshal(soal)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data for caching: %v", err)
	}

	err = rdb.Set(ctx, cacheKey, dataToCache, 10*time.Minute).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to set cache: %v", err)
	}

	return soal, nil
}

func (r *SoalRepositoryImpl) UpdateSoal(ctx context.Context, soal models.SoalExsample) error {
	updated_at := time.Now()
	_, err := r.db.ExecContext(ctx, queryUpdateSoal, soal.CategoryId, soal.Soal, soal.JawabanA, soal.JawabanB, soal.JawabanC, soal.JawabanD, soal.JawabanE, soal.JawabanBenar, updated_at, soal.SoalId)
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

func (r *SoalRepositoryImpl) GetSoalById(ctx context.Context, id int64) (soal models.SoalExsample, err error) {

	err = r.db.QueryRowContext(ctx, queryGetSoalById, id).Scan(
		&soal.SoalId,
		&soal.CategoryId,
		&soal.Soal,
		&soal.JawabanA,
		&soal.JawabanB,
		&soal.JawabanC,
		&soal.JawabanD,
		&soal.JawabanE,
		&soal.JawabanBenar,
		&soal.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.SoalExsample{}, fmt.Errorf("soal not found")
		}
		return models.SoalExsample{}, fmt.Errorf("error querying soal: %w", err)
	}

	return soal, nil
}
