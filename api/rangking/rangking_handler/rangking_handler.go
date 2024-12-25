package rangkinghandler

import (
	rangkingservice "bank_soal/api/rangking/rangking_service"
	"bank_soal/models"
	"bank_soal/utils/https"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type RangkingHandlerImpl struct {
	service rangkingservice.RangkingService
}

func NewRangkingHandler(service rangkingservice.RangkingService) *RangkingHandlerImpl {
	return &RangkingHandlerImpl{service: service}
}

func (h *RangkingHandlerImpl) CreateRank(e echo.Context) error {
	fName := "Rank_handler.CreateRank"
	ctx := e.Request().Context()

	type reqBody struct {
		RankId     int64 `json:"rank_id"`
		UserId     int64 `json:"user_id"`
		CategoryId int64 `json:"category_id"`
		Point      int64 `json:"point"`
	}

	req := reqBody{}
	if err := e.Bind(&req); err != nil {
		return https.WriteBadRequestResponse(e, https.ResponseBadRequestError)
	}

	// Validate request body
	if err := validator.New().Struct(&req); err != nil {
		return https.WriteBadRequestResponseWithErrMsg(e, https.ResponseBadRequestError, err)
	}

	resp := models.Rangking{
		RankId:     req.RankId,
		UserId:     req.UserId,
		CategoryId: req.CategoryId,
		Point:      req.Point,
	}

	result, err := h.service.CreateRank(ctx, resp)
	if err != nil {
		return https.WriteServerErrorResponse(e, fName, err)
	}

	return https.WriteOkResponse(e, fmt.Sprintf("Berhasil membuat soal dengan id %d", result))
}

func (h *RangkingHandlerImpl) GetRank(e echo.Context) error {
	fName := "Rank_handler.GetRank"
	ctx := e.Request().Context()

	var rank models.RangkingRes
	result, err := h.service.GetRank(ctx, rank)
	if err != nil {
		return https.WriteServerErrorResponse(e, fName, err)
	}

	// type resBody struct {
	// 	UserId int64 `json:"user_id"`
	// 	Point  int64 `json:"point"`
	// }

	// response := make([]resBody, len(result))
	// for i, r := range result {
	// 	response[i] = resBody{
	// 		UserId: r.UserId,
	// 		Point:  r.Point,
	// 	}
	// }

	return e.JSON(http.StatusOK, result)

}

func (h *RangkingHandlerImpl) CheckRank(e echo.Context) error {
	fName := "Rank_handler.CheckRank"
	ctx := e.Request().Context()

	userIdStr := e.QueryParam("user_id")
	if userIdStr == "" {
		return https.WriteBadRequestResponseWithErrMsg(e, https.ResponseBadRequestError, fmt.Errorf("missing or invalid user_id parameter"))
	}
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil || userId <= 0 {
		return https.WriteBadRequestResponseWithErrMsg(e, https.ResponseBadRequestError, fmt.Errorf("invalid userId parameter"))
	}

	result, err := h.service.CheckRank(ctx, userId)

	if err != nil {
		return https.WriteServerErrorResponse(e, fName, err)
	}

	return e.JSON(http.StatusOK, result)
}

func (h *RangkingHandlerImpl) UpdateRank(e echo.Context) error {
	fName := "Rank_handler.CheckRank"
	ctx := e.Request().Context()
	userIdStr := e.QueryParam("user_id")
	if userIdStr == "" {
		return https.WriteBadRequestResponseWithErrMsg(e, https.ResponseBadRequestError, fmt.Errorf("missing or invalid user_id parameter"))
	}
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil || userId <= 0 {
		return https.WriteBadRequestResponseWithErrMsg(e, https.ResponseBadRequestError, fmt.Errorf("invalid userId parameter"))
	}

	categoryIdStr := e.QueryParam("category_id")
	if userIdStr == "" {
		return https.WriteBadRequestResponseWithErrMsg(e, https.ResponseBadRequestError, fmt.Errorf("missing or invalid category_id parameter"))
	}
	categoryId, err := strconv.ParseInt(categoryIdStr, 10, 64)
	if err != nil || categoryId <= 0 {
		return https.WriteBadRequestResponseWithErrMsg(e, https.ResponseBadRequestError, fmt.Errorf("invalid category_id parameter"))
	}

	type reqBody struct {
		RankId     int64 `json:"rank_id"`
		UserId     int64 `json:"user_id"`
		CategoryId int64 `json:"category_id"`
		Point      int64 `json:"point"`
	}

	req := reqBody{}
	if err := e.Bind(&req); err != nil {
		return https.WriteBadRequestResponse(e, https.ResponseBadRequestError)
	}

	// Validate request body
	if err := validator.New().Struct(&req); err != nil {
		return https.WriteBadRequestResponseWithErrMsg(e, https.ResponseBadRequestError, err)
	}

	resp := models.Rangking{
		RankId:     req.RankId,
		UserId:     userId,
		CategoryId: categoryId,
		Point:      req.Point,
	}
	_, err = h.service.UpdatedRank(ctx, resp)

	if err != nil {
		return https.WriteServerErrorResponse(e, fName, err)
	}

	return https.WriteOkResponse(e, fmt.Sprintf("Berhasil update rank "))
}
