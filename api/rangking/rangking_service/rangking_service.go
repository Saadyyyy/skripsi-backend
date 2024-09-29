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
	UpdateNextUser(ctx context.Context, rank models.Rangking) (id int64, err error)
	CheckingRank(ctx context.Context, userId, soalId, categoryId int64) (check models.CheckRank, err error)
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
	// Check if the category exists
	ctId, err := s.repoCategory.GetCategoryByID(ctx, rank.CategoryId)
	if err != nil {
		return 0, fmt.Errorf("failed to get GetCategoryByID: %w", err)
	}

	// Check if the soal exists
	soalId, err := s.repoSoal.GetSoalById(ctx, rank.SoalId)
	if err != nil {
		return 0, fmt.Errorf("failed to get GetSoalById: %w", err)
	}

	// Check if the user exists
	uId, err := s.repoUser.GetUserByID(ctx, rank.UserId)
	if err != nil {
		return 0, fmt.Errorf("failed to get GetUserByID: %w", err)
	}

	// // Check if the rank with the same userId, soalId, and categoryId already exists
	// checkRank, err := s.repoRank.CheckingRank(ctx, uId.UserId, soalId.SoalId, ctId.CategoryId)
	// if err != nil {
	// 	return 0, fmt.Errorf("failed to check existing rank: %w", err)
	// }
	// if checkRank.RangkingId != 0 { // If a rank exists
	// 	return 0, fmt.Errorf("failed: UserId, SoalId, and CategoryId combination already exists")
	// }

	// Prepare the ranking response
	resp := models.Rangking{
		UserId:     uId.UserId,
		CategoryId: ctId.CategoryId,
		SoalId:     soalId.SoalId,
		Point:      100,
		Next:       false,
	}

	// Create the ranking
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

	return rank, nil
}

func (s *RangkingServiceImpl) GetUserAndPoint(ctx context.Context) (rank []models.RangkingUser, err error) {
	rank, err = s.repoRank.GetUserAndPoint(ctx)
	if err != nil {
		return nil, fmt.Errorf("Gagal get GetUserAndPoint %d :", err)
	}
	return rank, nil
}

func (s *RangkingServiceImpl) UpdateNextUser(ctx context.Context, rank models.Rangking) (id int64, err error) {
	id, err = s.repoRank.UpdateNextUser(ctx, rank)
	if err != nil {
		return 0, fmt.Errorf("gagal get UpdateNextUser")
	}

	return id, nil
}

func (s *RangkingServiceImpl) CheckingRank(ctx context.Context, userId, soalId, categoryId int64) (check models.CheckRank, err error) {
	check, err = s.repoRank.CheckingRank(ctx, userId, soalId, categoryId)
	if err != nil {
		return models.CheckRank{}, fmt.Errorf("Gagal get GetUserAndPoint %d :", err)
	}
	return check, nil

}
