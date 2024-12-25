package rangkingservice

import (
	rangkingrepository "bank_soal/api/rangking/rangking_repository"
	repository "bank_soal/api/user/user_repository"
	"bank_soal/models"
	"context"
	"fmt"
)

type RangkingService interface {
	CreateRank(ctx context.Context, rank models.Rangking) (ID int64, err error)
	UpdatedRank(ctx context.Context, rank models.Rangking) (ID int64, err error)
	CheckRank(ctx context.Context, ID int64) (result []models.Rangking, err error)
	GetRank(ctx context.Context, rank models.RangkingRes) (result []models.RangkingRes, err error)
}

type RangkingServiceImpl struct {
	repo     rangkingrepository.RangkingRepository
	repoUser repository.UserRepositoryInterface
}

func (r *RangkingServiceImpl) CheckRank(ctx context.Context, ID int64) (result []models.Rangking, err error) {
	result, err = r.repo.CheckRank(ctx, ID)
	if err != nil {
		return []models.Rangking{}, fmt.Errorf("gagal memanggil repositry error :", err)
	}

	return result, nil
}

func (r *RangkingServiceImpl) CreateRank(ctx context.Context, rank models.Rangking) (ID int64, err error) {
	check, err2 := r.repo.CheckRank(ctx, rank.UserId)
	if err2 != nil {
		return 0, fmt.Errorf("gagal get function checkrank dari repository", err2)
	}

	for _, v := range check {
		if v.CategoryId == rank.CategoryId && v.UserId == rank.UserId {
			if v.Point < rank.Point {
				ID, err = r.repo.UpdatedRank(ctx, rank)
				if err != nil {
					return 0, fmt.Errorf("gagal menjalankan fungsi UpdateRank dari repository: %w", err)
				}
				return v.RankId, nil
			}
			return v.RankId, nil
		}
	}
	ID, err = r.repo.CreateRank(ctx, rank)

	if err != nil {
		return 0, fmt.Errorf("Gagal mendapatkan mendapatkan fungsi dari reposiotry create rank error :", err)
	}
	return ID, nil
}

func (r *RangkingServiceImpl) GetRank(ctx context.Context, rank models.RangkingRes) (result []models.RangkingRes, err error) {

	result, err = r.repo.GetRank(ctx, rank)
	if err != nil {
		return []models.RangkingRes{}, fmt.Errorf("Gagal mendapatkan mendapatkan fungsi dari reposiotry create rank error :", err)
	}

	nilai := int64(100)

	for i := range result {
		result[i].Point = result[i].Point * nilai
	}

	return result, nil
}

func (r *RangkingServiceImpl) UpdatedRank(ctx context.Context, rank models.Rangking) (ID int64, err error) {

	ID, err = r.repo.UpdatedRank(ctx, rank)
	if err != nil {
		return 0, fmt.Errorf("Gagal mendapatkan update rank dari repository error :", err)
	}

	return ID, nil
}

func NewRangkingService(repo rangkingrepository.RangkingRepository, repoUser repository.UserRepositoryInterface) RangkingService {
	return &RangkingServiceImpl{repo: repo, repoUser: repoUser}
}
