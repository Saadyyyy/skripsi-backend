package rangkinghandler

import (
	rangkingservice "bank_soal/api/rangking/rangking_service"
	"bank_soal/models"
	"bank_soal/utils/https"
	"fmt"
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

func (h *RangkingHandlerImpl) CreateRangking(e echo.Context) error {
	fName := "rangking_handler.CreateRangking"
	ctx := e.Request().Context()

	type reqBody struct {
		UserId     int64 `json:"user_id"`
		CategoryId int64 `json:"category_id"`
		SoalId     int64 `json:"soal_id"`
		Next       bool  `json:"next"`
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
		UserId:     req.UserId,
		CategoryId: req.CategoryId,
		SoalId:     req.SoalId,
		Next:       req.Next,
		Point:      req.Point,
	}

	result, err := h.service.CreateRangking(ctx, resp)
	if err != nil {
		return https.WriteServerErrorResponse(e, fName, err)
	}
	return https.WriteOkResponse(e, fmt.Sprintf("Berhasil membuat ranking %d", result))
}

func (h *RangkingHandlerImpl) GetPointByUserId(e echo.Context) error {
	fName := "rangking_handler.GetPointByUserId"
	ctx := e.Request().Context()

	userIdStr := e.QueryParam("user_id")
	if userIdStr == "" {
		return https.WriteBadRequestResponseWithErrMsg(e, https.ResponseBadRequestError, fmt.Errorf("missing or invalid category_id parameter"))
	}
	UserId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		return https.WriteBadRequestResponseWithErrMsg(e, https.ResponseBadRequestError, fmt.Errorf("invalid category_id parameter"))
	}

	result, err := h.service.GetPointByUserId(ctx, UserId)
	if err != nil {
		return https.WriteServerErrorResponse(e, fName, err)
	}

	type Resp struct {
		Point int64 `json:"point"`
	}

	resp := Resp{
		Point: result.Point,
	}

	return https.WriteOkResponse(e, resp)
}
