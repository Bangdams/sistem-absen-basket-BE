package converter

import (
	"absen-qr-backend/internal/entity"
	"absen-qr-backend/internal/model"
	"log"
)

func AttendanceLogToResponse(attendanceLog *entity.AttendanceLog) *model.AttendanceLogResponse {
	log.Println("log from attendanceLog to response")

	return &model.AttendanceLogResponse{
		ID:        &attendanceLog.ID,
		SessionId: attendanceLog.SessionId,
		StudentId: attendanceLog.SessionId,
		Title:     attendanceLog.Session.Title,
		CreatedAt: attendanceLog.Session.CreatedAt,
		StartedAt: attendanceLog.Session.StartedAt,
		ExpiresAt: attendanceLog.Session.ExpiresAt,
		Nis:       attendanceLog.Student.Nis,
		FullName:  attendanceLog.Student.FullName,
		ScannedAt: attendanceLog.ScannedAt,
		Status:    attendanceLog.Status,
	}
}

func AttendanceLogToResponses(attendanceLogs *[]entity.AttendanceLog) *[]model.AttendanceLogResponse {
	var attendanceLogResponses []model.AttendanceLogResponse

	log.Println("log from attendanceLog to responses")

	for _, attendanceLog := range *attendanceLogs {
		attendanceLogResponses = append(attendanceLogResponses, *AttendanceLogToResponse(&attendanceLog))
	}

	return &attendanceLogResponses
}

// for student
func StudentAttendanceLogToResponse(attendanceLog *entity.AttendanceLog) *model.StudentAttendanceLogResponse {
	log.Println("log from attendanceLog to response")

	return &model.StudentAttendanceLogResponse{
		ScannedAt: attendanceLog.ScannedAt,
		Title:     attendanceLog.Session.Title,
		StartedAt: attendanceLog.Session.StartedAt,
		Status:    attendanceLog.Status,
	}
}

func StudentAttendanceLogToResponses(attendanceLogs *[]entity.AttendanceLog) *[]model.StudentAttendanceLogResponse {
	var studentAttendanceLogResponses []model.StudentAttendanceLogResponse

	log.Println("log from attendanceLog to responses")

	for _, attendanceLog := range *attendanceLogs {
		studentAttendanceLogResponses = append(studentAttendanceLogResponses, *StudentAttendanceLogToResponse(&attendanceLog))
	}

	return &studentAttendanceLogResponses
}
