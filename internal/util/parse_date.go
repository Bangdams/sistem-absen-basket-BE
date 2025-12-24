package util

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

func ParseStringToDate(data string) (time.Time, error) {
	layout := "2006-01-02"

	// Parsing string to time.Time
	value, err := time.Parse(layout, data)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return time.Time{}, fiber.NewError(fiber.ErrBadRequest.Code, "Bad Request")
	}

	return value, nil
}

func CompareSessionTime(sessionStartedAt, requestStartedAt, sessionExpiresAt, requestExpiresAt string) bool {
	format := "15:04"

	startedSession, _ := time.Parse("15:04:05", sessionStartedAt)
	startedRequest, _ := time.Parse("15:04", requestStartedAt)

	expiresSession, _ := time.Parse("15:04:05", sessionExpiresAt)
	expiresRequest, _ := time.Parse("15:04", requestExpiresAt)

	if startedSession.Format(format) != startedRequest.Format(format) ||
		expiresSession.Format(format) != expiresRequest.Format(format) {
		return true
	}

	return false
}

func ParseTimeToday(data string) (time.Time, error) {
	layout := "15:04"

	// Parsing string to time.Time
	value, err := time.Parse(layout, data)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return time.Time{}, fiber.NewError(fiber.ErrBadRequest.Code, "Bad Request")
	}

	now := time.Now()
	hour, minute := value.Hour(), value.Minute()

	customDate := time.Date(
		now.Year(), now.Month(), now.Day(),
		hour, minute, 0, 0, time.Local)

	return customDate, nil
}
