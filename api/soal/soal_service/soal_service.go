package soal_service

import (
	repository "bank_soal/api/soal/soal_repository"
	userRepo "bank_soal/api/user/user_repository"
	"bank_soal/models"
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

type SoalServiceInterface interface {
	CreateSoal(ctx context.Context, soal models.SoalExsample) (ID int64, err error)
	GetSoal(ctx context.Context, filter models.FilterSoal, rdb *redis.Client) (soal []models.SoalExsample, totalData int64, err error)
	UpdateSoal(ctx context.Context, soal models.SoalExsample) error
	DeletedSoal(ctx context.Context, ID int64) error
	GetSoalById(ctx context.Context, ID int64) (soal models.SoalExsample, err error)
}

type SoalServiceImpl struct {
	repo     repository.SoalRepositoryInterface
	RepoUser userRepo.UserRepositoryInterface
	db       *sqlx.DB
}

func NewSoalService(repo repository.SoalRepositoryInterface, db *sqlx.DB) SoalServiceInterface {
	return &SoalServiceImpl{repo: repo, db: db}
}

func (s *SoalServiceImpl) CreateSoal(ctx context.Context, soal models.SoalExsample) (ID int64, err error) {
	if soal.Soal == "" || soal.JawabanB == "" || soal.JawabanD == "" || soal.JawabanA == "" || soal.JawabanC == "" || soal.JawabanE == "" || soal.JawabanBenar == "" {
		return 0, fmt.Errorf("tidak boleh kosong harus di isi")
	}

	ID, err = s.repo.CreateSoal(ctx, soal)
	if err != nil {
		return 0, fmt.Errorf("gagal memanggil fungsi dari repository %+v", err)
	}

	return ID, nil
}

func (s *SoalServiceImpl) GetSoal(ctx context.Context, filter models.FilterSoal, rdb *redis.Client) ([]models.SoalExsample, int64, error) {
	params := map[string]interface{}{
		"deleted_at":   nil,
		"custom_query": "",
	}

	if filter.Keyword != "" {
		keywordLower := strings.ToLower(filter.Keyword)
		params["custom_query"] = fmt.Sprintf("%s AND LOWER(soal) LIKE '%%%s%%'", params["custom_query"], keywordLower)
	}

	if filter.Category != 0 {
		categoryString := strconv.Itoa(int(filter.Category))
		params["custom_query"] = fmt.Sprintf("%s AND soal_exsamples.category_id='%s'", params["custom_query"], categoryString)
	}

	soal, err := s.repo.GetSoal(ctx, params, rdb)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get soal from repository: %+v", err)
	}

	totalData, err := s.repo.CountSoal(ctx, params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count soal from repository: %+v", err)
	}

	resp := make([]models.SoalExsample, len(soal))
	for i, s := range soal {
		resp[i] = models.SoalExsample{
			SoalId:       s.SoalId,
			CategoryId:   s.CategoryId,
			Soal:         s.Soal,
			JawabanA:     s.JawabanA,
			JawabanB:     s.JawabanB,
			JawabanC:     s.JawabanC,
			JawabanD:     s.JawabanD,
			JawabanE:     s.JawabanE,
			JawabanBenar: s.JawabanBenar,
			CreatedAt:    s.CreatedAt,
		}
	}

	return resp, totalData, nil
}

func (s *SoalServiceImpl) UpdateSoal(ctx context.Context, soal models.SoalExsample) error {

	err := s.repo.UpdateSoal(ctx, soal)
	if err != nil {
		return fmt.Errorf("gagal get update soal dari repository %+v", err)
	}

	return nil
}

func (s *SoalServiceImpl) DeletedSoal(ctx context.Context, ID int64) error {
	err := s.repo.DeleteSoal(ctx, ID)
	if err != nil {
		return fmt.Errorf("gagal get DeleteSoal dari repository %+v ", err)
	}

	return nil
}

func (s *SoalServiceImpl) GetSoalById(ctx context.Context, id int64) (result models.SoalExsample, err error) {
	result, err = s.repo.GetSoalById(ctx, id)
	if err != nil {
		return models.SoalExsample{}, fmt.Errorf("gagal getSoalById dari repository %+v", err)
	}
	return result, err
}
