package model

type WebResponse[T any] struct {
	Data   T              `json:"data"`
	Errors *ErrorResponse `json:"errors"`
}

type WebResponses[T any] struct {
	Data   *[]T           `json:"data"`
	Errors *ErrorResponse `json:"errors"`
}

type ErrorResponse struct {
	Message string   `json:"message"`
	Details []string `json:"details"`
}
