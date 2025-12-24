package converter

import (
	"absen-qr-backend/internal/entity"
	"absen-qr-backend/internal/model"
	"log"
)

func SessionToResponse(session *entity.Session) *model.SessionResponse {
	log.Println("log from session to response")

	return &model.SessionResponse{
		ID:          session.ID,
		CoachId:     session.CoachId,
		FullName:    session.Coach.FullName,
		Title:       session.Title,
		QrCodeImage: session.QrToken,
		StartedAt:   session.StartedAt,
		ExpiresAt:   session.ExpiresAt,
		CreatedAt:   session.CreatedAt,
	}
}

func SessionToResponses(sessions *[]entity.Session) *[]model.SessionResponse {
	var sessionResponses []model.SessionResponse

	log.Println("log from session to responses")

	for _, session := range *sessions {
		sessionResponses = append(sessionResponses, *SessionToResponse(&session))
	}

	return &sessionResponses
}
