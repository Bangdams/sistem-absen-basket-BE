package converter

import (
	"absen-qr-backend/internal/entity"
	"absen-qr-backend/internal/model"
	"absen-qr-backend/internal/util"
	"fmt"
	"log"
)

func PrintAttendanceReportResponse(session *entity.Session) *model.PrintAttendanceReportResponse {
	log.Println("log from PrintAttendanceReport to response")

	data := &model.PrintAttendanceReportResponse{
		SessionId: int(session.ID),
		Title:     session.Title,
		Date:      session.CreatedAt.Format("2006-01-02"),
		Time:      fmt.Sprintf("%s - %s", util.ParseToHourMinute(session.StartedAt), util.ParseToHourMinute(session.ExpiresAt)),
		Attendees: []model.AttendanceReportResponse{},
	}

	for _, attendanceLog := range session.AttendanceLog {
		attendanceReport := model.AttendanceReportResponse{
			Nis:      attendanceLog.Student.Nis,
			FullName: attendanceLog.Student.FullName,
			Status:   attendanceLog.Status,
		}

		data.Attendees = append(data.Attendees, attendanceReport)
	}

	return data
}

func PrintAttendanceReportResponses(sessions *[]entity.Session) *[]model.PrintAttendanceReportResponse {
	var printAttendanceReportResponses []model.PrintAttendanceReportResponse

	log.Println("log from PrintAttendanceReport to responses")

	for _, session := range *sessions {
		printAttendanceReportResponses = append(printAttendanceReportResponses, *PrintAttendanceReportResponse(&session))
	}

	return &printAttendanceReportResponses
}
