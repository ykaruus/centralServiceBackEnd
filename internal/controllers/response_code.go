package controllers

import "centralService/internal/http/response"

type HTTPResponse struct {
	Code     int
	ApiError *response.ErrorBody
}

func New(code int) *HTTPResponse {
	return &HTTPResponse{
		Code: code,
	}
}

func Wrap(code int, apierror *response.ErrorBody) *HTTPResponse {
	return &HTTPResponse{
		Code:     code,
		ApiError: apierror,
	}
}

func (h *HTTPResponse) Error() string {
	return h.ApiError.Message
}
