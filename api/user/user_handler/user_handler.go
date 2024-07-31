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
	"github.com/labstack/echo"
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
	fmt.Println("username", req.UsernameOrEmail)
	fmt.Println("password", req.Password)
	fmt.Println("user", user)
	fmt.Println("token", token)

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
	fmt.Println("Received token:", e.Request().Header.Get("Authorization"))

	var filter models.FilterUser
	filter.Keyword = e.QueryParam("keyword")
	page, err := strconv.Atoi(e.QueryParam("page"))
	if err != nil {
		page = 1
	}

	id, role, err := middleware.ExtractToken(e)
	if err != nil {
		return e.JSON(http.StatusBadRequest, map[string]string{"status": "BAD_REQUEST", "data": "invalid token"})
	}

	fmt.Println(id)

	if role == "user" {
		return https.WriteBadRequestResponse(e, https.ResponseBadRequestError)
	}
	fmt.Println("ID:", id)
	fmt.Println("Role:", role)
	fmt.Println("Keyword:", filter.Keyword)
	fmt.Println("Page:", page)
	fmt.Println("Limit:", filter.Limit)

	filter.Page = page
	limit, err := strconv.Atoi(e.QueryParam("per_page"))
	if err != nil {
		limit = 10
	}
	filter.Limit = limit

	soal, totalPage, totalData, err := h.usecase.GetAllUser(ctx, filter)
	if err != nil {
		return https.WriteServerErrorResponse(e, fName, err)
	}

	response := map[string]interface{}{
		"data":       soal,
		"total_page": totalPage,
		"total_data": totalData,
	}

	return e.JSON(http.StatusOK, response)
}
