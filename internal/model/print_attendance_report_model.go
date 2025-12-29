package model

type AttendanceReportResponse struct {
	Nis      string `json:"nis" validate:"required"`
	FullName string `json:"full_name" validate:"required"`
	Status   string `json:"status" validate:"required"`
}

type PrintAttendanceReportResponse struct {
	SessionId int                        `json:"session_id"`
	Title     string                     `json:"title"`
	Date      string                     `json:"date"`
	Time      string                     `json:"time"`
	Attendees []AttendanceReportResponse `json:"attendees"`
}
