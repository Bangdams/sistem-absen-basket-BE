package model

import "time"

type AttendanceLogResponse struct {
	ID        *uint      `json:"id" validate:"required"`
	SessionId uint       `json:"session_id" validate:"required"`
	StudentId uint       `json:"student_id" validate:"required"`
	Title     string     `json:"title" validate:"required"`
	CreatedAt time.Time  `json:"created_at" validate:"required"`
	StartedAt string     `json:"started_at" validate:"required"`
	ExpiresAt string     `json:"expires_at" validate:"required"`
	Nis       string     `json:"nis" validate:"required"`
	FullName  string     `json:"full_name" validate:"required"`
	ScannedAt *time.Time `json:"scanned_at" validate:"required"`
	Status    string     `json:"status" validate:"required"`
}

type StudentAttendanceLogResponse struct {
	ScannedAt *time.Time `json:"scanned_at" validate:"required"`
	Title     string     `json:"title" validate:"required"`
	StartedAt string     `json:"started_at" validate:"required"`
	Status    string     `json:"status" validate:"required"`
}

type UpdateAttendanceLogRequest struct {
	SessionId uint   `json:"session_id" validate:"required"`
	StudentId uint   `json:"student_id" validate:"required"`
	Status    string `json:"status" validate:"required,oneof=hadir sakit alpa"`
}

type AttendanceLogRequest struct {
	QrToken string `json:"qr_token" validate:"required"`
}
