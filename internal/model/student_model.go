package model

type StudentResponse struct {
	FullName    string `json:"full_name" validate:"required"`
	Role        string `json:"role" validate:"required"`
	Address     string `json:"address" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
}

type StudentRequest struct {
	Nis         string `json:"nis" validate:"required"`
	FullName    string `json:"full_name" validate:"required"`
	Address     string `json:"address" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
}
