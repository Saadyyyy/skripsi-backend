package soal_handler

import (
	service "bank_soal/api/soal/soal_service"
	"bank_soal/models"
	"bank_soal/utils/https"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type SoalHandler struct {
	service service.SoalServiceInterface
}

func NewSoalHandler(service service.SoalServiceInterface) *SoalHandler {
	return &SoalHandler{service: service}
}

// handler/soal_handler.go
func (h *SoalHandler) CreateSoal(e echo.Context) error {
	fName := "Soal_handler.CreateSoal"
	ctx := e.Request().Context()

	type reqBody struct {
		SoalID       int64  `json:"soal_id"`
		CategoryID   int64  `json:"category_id"`
		Soal         string `json:"soal"`
		JawabanA     string `json:"jawaban_a"`
		JawabanB     string `json:"jawaban_b"`
		JawabanC     string `json:"jawaban_c"`
		JawabanD     string `json:"jawaban_d"`
		JawabanBenar string `json:"jawaban_benar"`
	}

	req := reqBody{}

	if err := e.Bind(&req); err != nil {
		return https.WriteBadRequestResponse(e, https.ResponseBadRequestError)
	}

	// Validate request body
	if err := validator.New().Struct(&req); err != nil {
		return https.WriteBadRequestResponseWithErrMsg(e, https.ResponseBadRequestError, err)
	}

	resp := models.Soals{
		SoalId:       req.SoalID,
		CategoryId:   req.CategoryID,
		Soal:         req.Soal,
		JawabanA:     req.JawabanA,
		JawabanB:     req.JawabanB,
		JawabanC:     req.JawabanC,
		JawabanD:     req.JawabanD,
		JawabanBenar: req.JawabanBenar,
	}

	result, err := h.service.CreateSoal(ctx, resp)
	if err != nil {
		return https.WriteServerErrorResponse(e, fName, err)
	}
	return https.WriteOkResponse(e, fmt.Sprintf("Berhasil membuat soal dengan id %d", result))
}

func (h *SoalHandler) GetSoal(e echo.Context) error {
	fName := "soal_handler.GetSoal"
	ctx := e.Request().Context()

	var filter models.FilterSoal
	category, _ := strconv.Atoi(e.QueryParam("category_id"))
	filter.Category = int64(category)
	// Bind query parameters manually
	filter.TglMulai = e.QueryParam("tgl_mulai")
	filter.TglSelesai = e.QueryParam("tgl_selesai")
	filter.Keyword = e.QueryParam("soal")

	soal, totalData, err := h.service.GetSoal(ctx, filter)
	if err != nil {
		return https.WriteServerErrorResponse(e, fName, err)
	}

	fmt.Println(soal)
	// Build the response
	response := map[string]interface{}{
		"data":       soal,
		"total_data": totalData,
	}

	return e.JSON(http.StatusOK, response)
}

func (h *SoalHandler) UpdateSoal(e echo.Context) error {
	fName := "soal_handle.UpdateSoal"
	ctx := e.Request().Context()
	soalIDStr := e.QueryParam("soal_id")
	if soalIDStr == "" {
		return https.WriteBadRequestResponseWithErrMsg(e, https.ResponseBadRequestError, fmt.Errorf("missing or invalid soal_id parameter"))
	}

	soalID, err := strconv.ParseInt(soalIDStr, 10, 64)
	if err != nil || soalID <= 0 {
		return https.WriteBadRequestResponseWithErrMsg(e, https.ResponseBadRequestError, fmt.Errorf("invalid soal_id parameter"))
	}

	// Define the request body structure
	type reqBody struct {
		ID           int64  `json:"soal_id"`
		CategoryID   int64  `json:"category_id"`
		Soal         string `json:"soal"`
		JawabanA     string `json:"jawaban_a"`
		JawabanB     string `json:"jawaban_b"`
		JawabanC     string `json:"jawaban_c"`
		JawabanD     string `json:"jawaban_d"`
		JawabanBenar string `json:"jawaban_benar"`
	}
	req := reqBody{}

	// Bind the request body
	if err := e.Bind(&req); err != nil {
		return https.WriteBadRequestResponse(e, https.ResponseBadRequestError)
	}

	// Validate the request body
	if err := validator.New().Struct(&req); err != nil {
		return https.WriteBadRequestResponseWithErrMsg(e, https.ResponseBadRequestError, err)
	}

	// Create the response model
	resp := models.Soals{
		SoalId:       soalID,
		CategoryId:   req.CategoryID,
		Soal:         req.Soal,
		JawabanA:     req.JawabanA,
		JawabanB:     req.JawabanB,
		JawabanC:     req.JawabanC,
		JawabanD:     req.JawabanD,
		JawabanBenar: req.JawabanBenar,
	}

	// Call the service to update the soal
	err = h.service.UpdateSoal(ctx, resp)
	if err != nil {
		return https.WriteServerErrorResponse(e, fName, err)
	}

	return https.WriteOkResponse(e, fmt.Sprintf("Berhasil memperbarui soal dengan ID: %d", resp.SoalId))
}
func (h *SoalHandler) DeletedSoal(e echo.Context) error {
	fName := "soal_handle.DeletedSoal"
	ctx := e.Request().Context()

	soalIDStr := e.QueryParam("soal_id")
	if soalIDStr == "" {
		return https.WriteBadRequestResponseWithErrMsg(e, https.ResponseBadRequestError, fmt.Errorf("missing or invalid soal_id parameter"))
	}

	soalID, err := strconv.ParseInt(soalIDStr, 10, 64)
	if err != nil || soalID <= 0 {
		return https.WriteBadRequestResponseWithErrMsg(e, https.ResponseBadRequestError, fmt.Errorf("invalid soal_id parameter"))
	}

	// Call the service to delete the soal
	err = h.service.DeletedSoal(ctx, soalID)
	if err != nil {
		return https.WriteServerErrorResponse(e, fName, err)
	}

	return https.WriteOkResponse(e, fmt.Sprintf("Berhasil menghapus soal dengan ID: %d", soalID))
}

func (h *SoalHandler) GetSoalById(e echo.Context) error {
	fName := "soal_handle.GetSoalById"
	ctx := e.Request().Context()

	soalIDStr := e.QueryParam("soal_id")
	if soalIDStr == "" {
		return https.WriteBadRequestResponseWithErrMsg(e, https.ResponseBadRequestError, fmt.Errorf("missing or invalid soal_id parameter"))
	}
	soalID, err := strconv.ParseInt(soalIDStr, 10, 64)
	if err != nil || soalID <= 0 {
		return https.WriteBadRequestResponseWithErrMsg(e, https.ResponseBadRequestError, fmt.Errorf("invalid soal_id parameter"))
	}

	resp, err := h.service.GetSoalById(ctx, soalID)
	if err != nil {
		return https.WriteServerErrorResponse(e, fName, err)
	}

	return https.WriteOkResponse(e, resp)
}
