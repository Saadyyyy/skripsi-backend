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
		return https.WriteBadRequestResponseWithErrMsg(e, https.ResponseBadRequestError, fmt.Errorf("missing or invalid soal_id parameter"))

	}
	categoryId, err := strconv.ParseInt(categoryIdStr, 10, 64)
	if err != nil {
		return https.WriteBadRequestResponseWithErrMsg(e, https.ResponseBadRequestError, fmt.Errorf("invalid soal_id parameter"))
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

	category, totalPage, totalData, err := h.serviceCategory.GetAllCategory(ctx, filter)
	if err != nil {
		https.WriteServerErrorResponse(e, fName, err)

	}

	respon := map[string]interface{}{
		"data":       category,
		"total_page": totalPage,
		"total_data": totalData,
	}
	return e.JSON(http.StatusOK, respon)

}

// func (h *SoalHandler) GetSoal(e echo.Context) error {
// 	fName := "soal_handler.GetSoal"
// 	ctx := e.Request().Context()

// 	var filter models.FilterSoal
// 	category, _ := strconv.Atoi(e.QueryParam("category_id"))
// 	filter.Category = int64(category)
// 	// Bind query parameters manually
// 	filter.TglMulai = e.QueryParam("tgl_mulai")
// 	filter.TglSelesai = e.QueryParam("tgl_selesai")
// 	filter.Keyword = e.QueryParam("soal")
// 	page, err := strconv.Atoi(e.QueryParam("page"))
// 	if err != nil {
// 		page = 1
// 	}

// 	filter.Page = page
// 	limit, err := strconv.Atoi(e.QueryParam("per_page"))
// 	if err != nil {
// 		limit = 10
// 	}
// 	filter.Limit = limit

// 	soal, totalPage, totalData, err := h.service.GetSoal(ctx, filter)
// 	if err != nil {
// 		return https.WriteServerErrorResponse(e, fName, err)
// 	}

// 	// Build the response
// 	response := map[string]interface{}{
// 		"data":       soal,
// 		"total_page": totalPage,
// 		"total_data": totalData,
// 	}

// 	return e.JSON(http.StatusOK, response)
// }
