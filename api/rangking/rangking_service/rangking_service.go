package rangkingservice

import (
	categoryrepository "bank_soal/api/category/category_repository"
	rangkingrepository "bank_soal/api/rangking/rangking_repository"
	"bank_soal/api/soal/soal_repository"
	repository "bank_soal/api/user/user_repository"
	"bank_soal/models"
	"context"
	"fmt"
)

type RangkingService interface {
	CreateRangking(ctx context.Context, rank models.Rangking) (id int64, err error)
	GetPointByUserId(ctx context.Context, id int64) (rank models.Rangking, err error)
	GetUserAndPoint(ctx context.Context) (rank []models.RangkingUser, err error)
}

type RangkingServiceImpl struct {
	repoRank     rangkingrepository.RangkingRepository
	repoSoal     soal_repository.SoalRepositoryInterface
	repoCategory categoryrepository.CategoryRepository
	repoUser     repository.UserRepositoryInterface
}

func NewRangkingService(
	repoRank rangkingrepository.RangkingRepository,
	repoSoal soal_repository.SoalRepositoryInterface,
	repoCategory categoryrepository.CategoryRepository,
	repoUser repository.UserRepositoryInterface) RangkingService {
	return &RangkingServiceImpl{repoRank: repoRank, repoSoal: repoSoal, repoCategory: repoCategory, repoUser: repoUser}
}

func (s *RangkingServiceImpl) CreateRangking(ctx context.Context, rank models.Rangking) (id int64, err error) {
	ctId, err := s.repoCategory.GetCategoryByID(ctx, rank.CategoryId)
	if err != nil {
		return 0, fmt.Errorf("failed to get GetCategoryByID: %w", err)
	}

	soalId, err := s.repoSoal.GetSoalById(ctx, rank.SoalId)
	if err != nil {
		return 0, fmt.Errorf("failed to get GetSoalById: %w", err)
	}

	uId, err := s.repoUser.GetUserByID(ctx, rank.UserId)
	if err != nil {
		return 0, fmt.Errorf("failed to get GetUserByID: %w", err)
	}

	resp := models.Rangking{
		UserId:     uId.UserId,
		CategoryId: ctId.CategoryId,
		SoalId:     soalId.SoalId,
		Point:      100,
		Next:       false,
	}

	id, err = s.repoRank.CreateRangking(ctx, resp)
	if err != nil {
		return 0, fmt.Errorf("failed to create CreateRangking: %w", err)
	}

	return id, nil
}

func (s *RangkingServiceImpl) GetPointByUserId(ctx context.Context, id int64) (rank models.Rangking, err error) {
	rank, err = s.repoRank.GetPointByUserId(ctx, id)
	if err != nil {
		return models.Rangking{}, fmt.Errorf("failed to create CreateRangking: %w", err)
	}
	// hasil := (rank.Point + rank.Point) - rank.Point
	// rank.Point = hasil

	return rank, nil
}

func (s *RangkingServiceImpl) GetUserAndPoint(ctx context.Context) (rank []models.RangkingUser, err error) {
	rank, err = s.repoRank.GetUserAndPoint(ctx)
	if err != nil {
		return nil, fmt.Errorf("Gagal get GetUserAndPoint : ", err)
	}
	return rank, nil
}
