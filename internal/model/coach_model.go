package model

type CoachResponse struct {
	FullName string `json:"full_name" validate:"required"`
	Role     string `json:"role" validate:"required"`
}

type CoachRequest struct {
	Nip      string `json:"nip" validate:"required"`
	FullName string `json:"full_name" validate:"required"`
}
