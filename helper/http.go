package helper

import (
	"github.com/labstack/echo/v4"
)

type (
	HttpResponse struct {
		Status  int                    `json:"status"`
		Message string                 `json:"message"`
		Data    map[string]interface{} `json:"data"`
	}
	HttpResponseData struct {
		Status  int         `json:"status"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
)

func Response(statuscode int, data interface{}, message string) error {
	response := map[string]interface{}{
		"status":  statuscode,
		"message": message,
		"data":    data,
	}

	return echo.NewHTTPError(statuscode, response)
}
