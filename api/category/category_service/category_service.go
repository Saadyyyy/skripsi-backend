package categoryservice

import (
	repository "bank_soal/api/category/category_repository"
	"bank_soal/models"
	"fmt"

	"golang.org/x/net/context"
)

type CategoryService interface {
	CreateCategory(ctx context.Context, ct models.Category) (id int64, err error)
	// GetCategoryByID(ctx context.Context, id int64) (ct models.Category, err error)
}

type CategoryServiceImpl struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &CategoryServiceImpl{categoryRepo: categoryRepo}
}

func (s *CategoryServiceImpl) CreateCategory(ctx context.Context, ct models.Category) (id int64, err error) {
	id, err = s.categoryRepo.CreateCategory(ctx, ct)
	if err != nil {
		return 0, fmt.Errorf("Failed to get category repository %+v", err)
	}

	return id, nil
}

// func (s *CategoryServiceImpl) GetCategoryByID(ctx context.Context, id int64) (ct models.Category, err error) {

// }
