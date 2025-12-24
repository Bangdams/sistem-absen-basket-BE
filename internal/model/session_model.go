package model

import "time"

type SessionResponse struct {
	ID          uint      `json:"id"`
	CoachId     uint      `json:"coach_id"`
	FullName    string    `json:"full_name"`
	Title       string    `json:"title"`
	QrCodeImage string    `json:"qr_code_image"`
	StartedAt   string    `json:"started_at"`
	ExpiresAt   string    `json:"expires_at"`
	CreatedAt   time.Time `json:"created_at"`
}

type SessionRequest struct {
	CoachId   uint   `json:"coach_id" validate:"required"`
	Title     string `json:"title" validate:"required"`
	StartedAt string `json:"started_at" validate:"required"`
	ExpiresAt string `json:"expires_at" validate:"required"`
}

type UpdateSessionRequest struct {
	ID        uint   `json:"id" validate:"required"`
	Title     string `json:"title" validate:"required"`
	StartedAt string `json:"started_at" validate:"required"`
	ExpiresAt string `json:"expires_at" validate:"required"`
}
