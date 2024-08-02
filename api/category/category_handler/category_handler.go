package categoryhandler

import (
	categoryservice "bank_soal/api/category/category_service"
	"bank_soal/models"
	"bank_soal/utils/https"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CategoryHandler struct {
	serviceCategory categoryservice.CategoryService
}

func NewCategoryHandler(serviceCategory categoryservice.CategoryService) *CategoryHandler {
	return &CategoryHandler{serviceCategory: serviceCategory}
}

// handler/soal_handler.go
func (h *CategoryHandler) CreateSoal(e echo.Context) error {
	fName := "Soal_handler.CreateSoal"
	ctx := e.Request().Context()

	type reqBody struct {
		CategoryID int64  `json:"category_id"`
		Category   string `json:"category"`
	}

	req := reqBody{}

	if err := e.Bind(&req); err != nil {
		return https.WriteBadRequestResponse(e, https.ResponseBadRequestError)
	}

	// Validate request body
	if err := validator.New().Struct(&req); err != nil {
		return https.WriteBadRequestResponseWithErrMsg(e, https.ResponseBadRequestError, err)
	}

	resp := models.Category{
		CategoryId: req.CategoryID,
		Category:   req.Category,
	}

	result, err := h.serviceCategory.CreateCategory(ctx, resp)
	if err != nil {
		return https.WriteServerErrorResponse(e, fName, err)
	}
	return https.WriteOkResponse(e, fmt.Sprintf("Berhasil membuat soal dengan id %d", result))
}
