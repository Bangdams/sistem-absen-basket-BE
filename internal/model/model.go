package model

type WebResponse[T any] struct {
	Data   T              `json:"data"`
	Errors *ErrorResponse `json:"errors"`
}

type WebResponses[T any] struct {
	Data   *[]T           `json:"data"`
	Errors *ErrorResponse `json:"errors"`
}

type WebResponsesPagination[T any] struct {
	Data         *[]T           `json:"data"`
	CurrentPage  int            `json:"current_page"`
	TotalRecords int            `json:"total_records"`
	TotalPages   int            `json:"total_pages"`
	Errors       *ErrorResponse `json:"errors"`
}

type ErrorResponse struct {
	Message string   `json:"message"`
	Details []string `json:"details"`
}
