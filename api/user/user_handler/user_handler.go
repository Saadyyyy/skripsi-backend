package handler

import (
	service "bank_soal/api/user/user_service"
	"bank_soal/middleware"
	"bank_soal/models"
	"bank_soal/utils/https"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	usecase service.UserService
}

func NewUserHandler(usecase service.UserService) *UserHandler {
	return &UserHandler{usecase: usecase}
}

func (h *UserHandler) CreateUser(e echo.Context) error {
	fName := "OrganisasiHttpHandler.InsertMasterPetaOrganisasi"
	ctx := e.Request().Context()

	type reqBody struct {
		UserId   int64  `json:"user_id"`
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	req := reqBody{}

	if err := e.Bind(&req); err != nil {
		return https.WriteBadRequestResponse(e, https.ResponseBadRequestError)
	}

	// validate request body
	if err := validator.New().Struct(&req); err != nil {
		return https.WriteBadRequestResponseWithErrMsg(e, https.ResponseBadRequestError, err)
	}

	resp := models.Users{
		UserId:   req.UserId,
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	}

	result, err := h.usecase.CreateUser(ctx, resp)
	if err != nil {
		return https.WriteServerErrorResponse(e, fName, err)
	}

	return https.WriteOkResponse(e, fmt.Sprintf("Berhasil membuat akun dengan id %d", result))
}
func (h *UserHandler) LoginUser(e echo.Context) error {
	fName := "user_handler.LoginUser"
	ctx := e.Request().Context()

	// Define the request body structure
	type reqBody struct {
		UsernameOrEmail string `json:"username_or_email" validate:"required"`
		Password        string `json:"password" validate:"required"`
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

	// Call the service to login the user
	user, token, err := h.usecase.LoginUser(ctx, req.UsernameOrEmail, req.Password)
	if err != nil {
		return https.WriteServerErrorResponse(e, fName, err)
	}

	// Set the token in a cookie
	middleware.SetTokenCookie(e, token)

	// Return user info and token in response
	response := map[string]interface{}{
		"user":  user,
		"token": token,
	}

	return https.WriteOkResponse(e, response)
}

func (h *UserHandler) UpdateUser(e echo.Context) error {
	fName := "OrganisasiHttpHandler.UpdateUser"
	ctx := e.Request().Context()
	type reqBody struct {
		UserId   int64  `json:"user_id"`
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	req := reqBody{}

	if err := e.Bind(&req); err != nil {
		return https.WriteBadRequestResponse(e, https.ResponseBadRequestError)
	}

	// validate request body
	if err := validator.New().Struct(&req); err != nil {
		return https.WriteBadRequestResponseWithErrMsg(e, https.ResponseBadRequestError, err)
	}

	resp := models.Users{
		UserId:   req.UserId,
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	}

	err := h.usecase.UpdateUser(ctx, resp)
	if err != nil {
		return https.WriteServerErrorResponse(e, fName, err)
	}

	return https.WriteOkResponse(e, fmt.Sprintf("Berhasil membuat akun dengan id %d", resp.UserId))
}

func (h *UserHandler) GetAllUser(e echo.Context) error {
	fName := "soal_handler.GetAllUser"
	ctx := e.Request().Context()

	var filter models.FilterUser
	filter.Keyword = e.QueryParam("keyword")
	page, err := strconv.Atoi(e.QueryParam("page"))
	if err != nil {
		page = 1
	}

	filter.Page = page
	limit, err := strconv.Atoi(e.QueryParam("per_page"))
	if err != nil {
		limit = 10
	}
	filter.Limit = limit

	soal, totalData, err := h.usecase.GetAllUser(ctx, filter)
	if err != nil {
		return https.WriteServerErrorResponse(e, fName, err)
	}

	response := map[string]interface{}{
		"data":       soal,
		"total_data": totalData,
	}

	return e.JSON(http.StatusOK, response)
}

func (h *UserHandler) GetUserByID(e echo.Context) error {
	fName := "soal_handler.GetAllUser"
	ctx := e.Request().Context()

	userIdStr := e.QueryParam("user_id")
	if userIdStr == "" {
		return https.WriteBadRequestResponseWithErrMsg(e, https.ResponseBadRequestError, fmt.Errorf("missing or invalid user_id parameter"))
	}

	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		return https.WriteBadRequestResponseWithErrMsg(e, https.ResponseBadRequestError, fmt.Errorf("invalid user_id parameter"))
	}

	resp, err := h.usecase.GetUserByID(ctx, userId)
	if err != nil {
		return https.WriteServerErrorResponse(e, fName, err)
	}

	type respon struct {
		UserId   int64  `json:"user_id"`
		Username string `json:"username"`
		Role     int64  `json:"role"`
		Profile  string `json:"profile"`
	}

	filter := respon{
		UserId:   resp.UserId,
		Username: resp.Username,
		Role:     resp.Role,
		Profile:  resp.Profile,
	}

	return https.WriteOkResponse(e, filter)
}

func (h *UserHandler) UpdateUserRoleByID(e echo.Context) error {
	fName := "OrganisasiHttpHandler.UpdaUpdateUserRoleByIDteUser"
	ctx := e.Request().Context()
	type reqBody struct {
		UserId int64 `json:"user_id"`
		Role   int64 `json:"role"`
	}

	req := reqBody{}

	if err := e.Bind(&req); err != nil {
		return https.WriteBadRequestResponse(e, https.ResponseBadRequestError)
	}

	// validate request body
	if err := validator.New().Struct(&req); err != nil {
		return https.WriteBadRequestResponseWithErrMsg(e, https.ResponseBadRequestError, err)
	}

	resp := models.Users{
		UserId: req.UserId,
		Role:   req.Role,
	}

	err := h.usecase.UpdateUserRoleByID(ctx, resp)
	if err != nil {
		return https.WriteServerErrorResponse(e, fName, err)
	}

	return https.WriteOkResponse(e, fmt.Sprintf("Berhasil membuat update role dengan id %d", resp.UserId))
}
