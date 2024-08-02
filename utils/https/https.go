package https

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

const (
	ResponseOk                   = "OK"
	ResponseServerError          = "SERVER_ERROR"
	ResponseBadRequestError      = "BAD_REQUEST"
	ResponseNotFoundError        = "NOT_FOUND"
	ResponseTimedOut             = "TIMED_OUT"
	ResponseUnauthorizedError    = "UNAUTHORIZED"
	ResponseUnauthenticatedError = "UNAUTHENTICATED"
	ResponseWrongPasswordError   = "WRONG_PASSWORD"
	ResponseDataAlreadyExist     = "DATA_ALREADY_EXIST"
	ResponseInvalidDates         = "INVALID_DATES"
)

// BaseResponse represents base http response
type BaseResponse struct {
	Status       string      `json:"status"`
	ErrorMessage string      `json:"error_message,omitempty"`
	ErrorCode    string      `json:"error_code,omitempty"`
	Data         interface{} `json:"data"`
}

// Pagination represents pagination for ListResponse
type Pagination struct {
	TotalPage int64 `json:"total_page"`
	TotalData int64 `json:"total_data"`
}

// ListResponse represents how list is returned as `data` in BaseResponse
type ListResponse struct {
	List       interface{} `json:"list"`
	Pagination Pagination  `json:"pagination"`
}

// func getTimedOutRespBody() string {
// 	timeOutBody := BaseResponse{
// 		Status: ResponseTimedOut,
// 	}
// 	marshal, _ := json.Marshal(timeOutBody)

// 	return string(marshal)
// }

// ===== Response wrapper using echo =====

// WriteServerErrorResponse writes server error response
func WriteServerErrorResponse(ctx echo.Context, functionName string, err error) error {
	return writeServerErrorResponse(ctx, functionName, err, "")
}

// WriteServerErrorResponseWithErrorCode writes server error response with error code
func WriteServerErrorResponseWithErrorCode(ctx echo.Context, functionName, errCode string, err error) error {
	return writeServerErrorResponse(ctx, functionName, err, errCode)
}

// writeServerErrorResponse helper function to write server error response
func writeServerErrorResponse(ctx echo.Context, functionName string, err error, errCode string) error {
	if err == nil {
		return WriteNotOkResponse(ctx, http.StatusInternalServerError, ResponseServerError)
	}

	errMessage := err.Error()
	log.Printf("[ERROR] [%s] %s\n", functionName, errMessage) // use log.Printf instead of log.Println + fmt.Sprintf

	if errCode == "" {
		return WriteNotOkResponseWithErrMsg(ctx, http.StatusInternalServerError, ResponseServerError, errMessage)
	} else {
		return WriteNotOkResponseWithErrorCodeAndMessage(ctx, http.StatusInternalServerError, ResponseServerError, errCode, errMessage)
	}
}

func WriteBadRequestResponse(ctx echo.Context, status string) error {
	return WriteNotOkResponse(ctx, http.StatusBadRequest, status)
}

func WriteBadRequestResponseWithErrMsg(ctx echo.Context, status string, err error) error {
	if err == nil {
		return WriteBadRequestResponse(ctx, status)
	}

	return WriteNotOkResponseWithErrMsg(ctx, http.StatusBadRequest, status, err.Error())
}

func WriteBadRequestResponseWithErrCodeAndMsg(ctx echo.Context, status, errCode string, err error) error {
	if err == nil {
		return WriteBadRequestResponse(ctx, status)
	}

	return WriteNotOkResponseWithErrorCodeAndMessage(ctx, http.StatusBadRequest, status, errCode, err.Error())
}

func WriteNotFoundResponse(ctx echo.Context, status string) error {
	return WriteNotOkResponse(ctx, http.StatusNotFound, status)
}

func WriteUnauthorizedResponse(ctx echo.Context) error {
	return WriteNotOkResponse(ctx, http.StatusUnauthorized, ResponseUnauthorizedError)
}

func WriteUnauthenticatedResponse(ctx echo.Context) error {
	return WriteNotOkResponse(ctx, http.StatusUnauthorized, ResponseUnauthenticatedError)
}

func WriteTimedOutResponse(ctx echo.Context) error {
	return WriteNotOkResponse(ctx, http.StatusGatewayTimeout, ResponseTimedOut)
}

func WriteWrongPasswordResponse(ctx echo.Context) error {
	return WriteNotOkResponse(ctx, http.StatusUnauthorized, ResponseWrongPasswordError)
}

// WriteOkResponse writes 200 response using echo.
func WriteOkResponse(ctx echo.Context, data interface{}) error {
	resp := BaseResponse{
		Status: ResponseOk,
		Data:   data,
	}
	return WriteResponse(ctx, resp, http.StatusOK)
}

// WriteNotOkResponse writes non 200 response.
func WriteNotOkResponse(ctx echo.Context, statusCode int, status string) error {
	resp := BaseResponse{
		Status: status,
	}
	return WriteResponse(ctx, resp, statusCode)
}

func WriteNotOkResponseWithErrorCodeAndMessage(ctx echo.Context, statusCode int, status, errCode, errMsg string) error {
	resp := BaseResponse{
		Status:       status,
		ErrorCode:    errCode,
		ErrorMessage: errMsg,
	}
	return WriteResponse(ctx, resp, statusCode)
}

// WriteNotOkResponseWithErrMsg writes non 200 responses with error message.
func WriteNotOkResponseWithErrMsg(ctx echo.Context, statusCode int, status, errMsg string) error {
	resp := BaseResponse{
		Status:       status,
		ErrorMessage: errMsg,
	}
	return WriteResponse(ctx, resp, statusCode)
}

func WriteResponse(ctx echo.Context, resp BaseResponse, statusCode int) error {
	switch {
	case statusCode >= http.StatusBadRequest:
		return echo.NewHTTPError(statusCode, resp)
	default:
		return ctx.JSON(statusCode, resp)
	}
}

// // WriteFileDownloadExcel writes excel file to response
// // under library github.com/xuri/excelize/v2
// func WriteFileDownloadExcel(ctx echo.Context, fileName string, file *excelize.File) error {
// 	ctx.Response().Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
// 	ctx.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
// 	ctx.Response().Header().Set("Content-Transfer-Encoding", "binary")
// 	ctx.Response().Header().Set("Expires", "0")

// 	return file.Write(ctx.Response())
// }

// // WriteFileDownloadExcel writes excel file to response
// // under library github.com/xuri/excelize/v2
// func WriteFileDownloadPdf(ctx echo.Context, fileName string, pdf *gofpdf.Fpdf) error {
// 	ctx.Response().Header().Set("Content-Type", "application/pdf")
// 	ctx.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
// 	ctx.Response().Header().Set("Content-Transfer-Encoding", "binary")
// 	ctx.Response().Header().Set("Expires", "0")

// 	return pdf.Output(ctx.Response())
// }
