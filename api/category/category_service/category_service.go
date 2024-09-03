package categoryservice

import (
	repository "bank_soal/api/category/category_repository"
	"bank_soal/models"
	"fmt"
	"strings"

	"golang.org/x/net/context"
)

type CategoryService interface {
	CreateCategory(ctx context.Context, ct models.Category) (id int64, err error)
	GetCategoryByID(ctx context.Context, id int64) (ct models.Category, err error)
	GetAllCategory(ctx context.Context, filter models.FilterCategory) (ct []models.Category, totalData int64, err error)
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
		return 0, fmt.Errorf("failed to get category repository %+v", err)
	}

	return id, nil
}

func (s *CategoryServiceImpl) GetCategoryByID(ctx context.Context, id int64) (ct models.Category, err error) {
	result, err := s.categoryRepo.GetCategoryByID(ctx, id)
	if err != nil {
		return models.Category{}, err
	}

	return result, nil
}

func (s *CategoryServiceImpl) GetAllCategory(ctx context.Context, filter models.FilterCategory) (ct []models.Category, totalData int64, err error) {
	params := map[string]interface{}{
		"deleted_at":   nil,
		"custom_query": "",
	}

	if filter.Keyword != "" {
		keywordLower := strings.ToLower(filter.Keyword)
		params["custom_query"] = fmt.Sprintf("%s AND LOWER(category) LIKE '%%%s%%'", params["custom_query"], keywordLower)
	}
	if filter.Page <= 0 {
		filter.Page = 1 // default page
	}

	if filter.Limit <= 0 || filter.Limit > 10 {
		filter.Limit = 10 // default limit
	}

	user, err := s.categoryRepo.GetListCategory(ctx, params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get soal from repository: %+v", err)
	}

	totalData, err = s.categoryRepo.CountUser(ctx, params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count soal from repository: %+v", err)
	}

	resp := make([]models.Category, len(user))
	for i, s := range user {
		resp[i] = models.Category{
			CategoryId: s.CategoryId,
			Category:   s.Category,
			CreatedAt:  s.CreatedAt,
		}
	}
	return resp, totalData, nil
}
