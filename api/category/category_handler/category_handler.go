package categoryhandler

import (
	categoryservice "bank_soal/api/category/category_service"
	"bank_soal/models"
	"bank_soal/utils/https"
	"fmt"
	"net/http"
	"strconv"

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

func (h *CategoryHandler) GetCategoryByID(e echo.Context) error {
	fName := "Soal_handler.CreateSoal"
	ctx := e.Request().Context()

	categoryIdStr := e.QueryParam("category_id")
	if categoryIdStr == "" {
		return https.WriteBadRequestResponseWithErrMsg(e, https.ResponseBadRequestError, fmt.Errorf("missing or invalid category_id parameter"))

	}
	categoryId, err := strconv.ParseInt(categoryIdStr, 10, 64)
	if err != nil {
		return https.WriteBadRequestResponseWithErrMsg(e, https.ResponseBadRequestError, fmt.Errorf("invalid category_id parameter"))
	}

	resp, err := h.serviceCategory.GetCategoryByID(ctx, categoryId)

	if err != nil {
		return https.WriteServerErrorResponse(e, fName, err)
	}

	return https.WriteOkResponse(e, resp)

}

func (h *CategoryHandler) GetListCategory(e echo.Context) error {
	fName := "Category_handler.GetListCategory"
	ctx := e.Request().Context()

	filter := models.FilterCategory{}

	filter.Keyword = e.QueryParam("keyword")
	page, err := strconv.Atoi(e.QueryParam("page"))
	if err != nil {
		page = 1
	}
	filter.Page = page

	limit, err := strconv.Atoi(e.QueryParam("per_page"))
	if err != nil {
		limit = 1
	}
	filter.Limit = limit

	category, totalData, err := h.serviceCategory.GetAllCategory(ctx, filter)
	if err != nil {
		https.WriteServerErrorResponse(e, fName, err)

	}

	respon := map[string]interface{}{
		"data":       category,
		"total_data": totalData,
	}
	return e.JSON(http.StatusOK, respon)

}

func (h *CategoryHandler) UpdatedCategory(e echo.Context) error {
	fName := "Category_handler.UpdatedCategory"
	ctx := e.Request().Context()

	categoryIdStr := e.QueryParam("category_id")
	if categoryIdStr == "" {
		return https.WriteBadRequestResponseWithErrMsg(e, https.ResponseBadRequestError, fmt.Errorf("missing or invalid category_id parameter"))
	}

	categoryId, err := strconv.ParseInt(categoryIdStr, 10, 64)
	if err != nil || categoryId <= 0 {
		return https.WriteBadRequestResponseWithErrMsg(e, https.ResponseBadRequestError, fmt.Errorf("failed convert catgeory_id"))

	}

	type reqBody struct {
		Category string `json:"category"`
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
		CategoryId: categoryId,
		Category:   req.Category,
	}

	err = h.serviceCategory.UpdatedCategory(ctx, resp)
	if err != nil {
		return https.WriteServerErrorResponse(e, fName, err)
	}
	return https.WriteOkResponse(e, fmt.Sprintf("Berhasil updated category id %d", categoryId))

}

func (h *CategoryHandler) DeletedCategory(e echo.Context) error {
	fName := "Category_handler.GetListCategory"
	ctx := e.Request().Context()

	categoryIdStr := e.QueryParam("category_id")
	if categoryIdStr == "" {
		return https.WriteBadRequestResponseWithErrMsg(e, https.ResponseBadRequestError, fmt.Errorf("missing or invalid category_id parameter"))
	}

	categoryId, err := strconv.ParseInt(categoryIdStr, 10, 64)
	if err != nil || categoryId <= 0 {
		return https.WriteBadRequestResponseWithErrMsg(e, https.ResponseBadRequestError, fmt.Errorf("failed convert catgeory_id"))

	}

	err = h.serviceCategory.DeletedCategory(ctx, categoryId)
	if err != nil {
		return https.WriteServerErrorResponse(e, fName, err)
	}
	return https.WriteOkResponse(e, fmt.Sprintf("Berhasil deleted category dengan id %d", categoryId))
}
