package view

import "net/http"

type ApiResponse struct {
	Headers map[string]string `json:"-"`
	Code    int 		`json:"code"`
	Data    interface{} `json:"data"`
	Message string		`json:"message"`
}

func Ok(data interface{}) ApiResponse {
	return ApiResponse{Code: http.StatusOK, Data: data, Message: "success"}
}

func BadRequest(data interface{}) ApiResponse {
	return ApiResponse{Code: http.StatusBadRequest, Data: data, Message: http.StatusText(http.StatusBadRequest)}
}
