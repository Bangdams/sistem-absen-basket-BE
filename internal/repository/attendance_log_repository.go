package repository

import (
	"absen-qr-backend/internal/entity"

	"gorm.io/gorm"
)

type AttendanceLogRepository interface {
	ExistsBySessionAndStudent(tx *gorm.DB, sessionId uint, studentId uint) error
	FindAllBySessionId(tx *gorm.DB, attendanceLogs *[]entity.AttendanceLog, sessionId uint) error
	FindAllByStudent(tx *gorm.DB, attendanceLogs *[]entity.AttendanceLog, studentId uint) error
	FindById(tx *gorm.DB, attendanceLog *entity.AttendanceLog) error
	Create(tx *gorm.DB, attendanceLog *entity.AttendanceLog) error
	Update(tx *gorm.DB, attendanceLog *entity.AttendanceLog, updateByUser bool) (int64, error)
	Delete(tx *gorm.DB, attendanceLog *entity.AttendanceLog) error
}

type AttendanceLogRepositoryImpl struct {
	Repository[entity.AttendanceLog]
}

func NewAttendanceLogRepository() AttendanceLogRepository {
	return &AttendanceLogRepositoryImpl{}
}

// ExistsBySessionAndStudent implements AttendanceLogRepository.
func (repository *AttendanceLogRepositoryImpl) ExistsBySessionAndStudent(tx *gorm.DB, sessionId uint, studentId uint) error {
	err := tx.Model(&entity.AttendanceLog{}).Where("session_id = ?", sessionId).Where("student_id = ?", studentId).First(&entity.AttendanceLog{}).Error
	if err != nil {
		return err
	}

	return nil
}

// FindAllBySessionId implements AttendanceLogRepository.
func (repository *AttendanceLogRepositoryImpl) FindAllBySessionId(tx *gorm.DB, attendanceLogs *[]entity.AttendanceLog, sessionId uint) error {
	err := tx.Where("session_id = ?", sessionId).
		Preload("Student", func(db *gorm.DB) *gorm.DB {
			return db.Select("user_id", "nis", "full_name")
		}).
		Preload("Session", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "title", "created_at", "started_at", "expires_at")
		}).
		Find(&attendanceLogs).Error

	if err != nil {
		return err
	}

	if len(*attendanceLogs) == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// FindAll implements AttendanceLogRepository.
func (repository *AttendanceLogRepositoryImpl) FindAllByStudent(tx *gorm.DB, attendanceLogs *[]entity.AttendanceLog, studentId uint) error {
	err := tx.Where("student_id = ?", studentId).
		Preload("Session", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "title", "started_at")
		}).
		Find(attendanceLogs).Error

	if err != nil {
		return err
	}

	if len(*attendanceLogs) == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// FindById implements AttendanceLogRepository.
func (repository *AttendanceLogRepositoryImpl) FindById(tx *gorm.DB, attendanceLog *entity.AttendanceLog) error {
	return tx.Where("id = ?", attendanceLog.ID).
		First(attendanceLog).Error
}

// Update implements AttendanceLogRepository.
// Subtle: this method shadows the method (Repository).Update of AttendanceLogRepositoryImpl.Repository.
func (repository *AttendanceLogRepositoryImpl) Update(tx *gorm.DB, attendanceLog *entity.AttendanceLog, updateByUser bool) (int64, error) {
	query := tx.Model(&entity.AttendanceLog{})

	if updateByUser {
		query = query.Where("status = ?", "Alpa")
	}

	result := query.
		Select("ScannedAt", "Status").
		Where("session_id = ?", attendanceLog.SessionId).
		Where("student_id = ?", attendanceLog.StudentId).
		Updates(map[string]interface{}{
			"ScannedAt": attendanceLog.ScannedAt,
			"Status":    attendanceLog.Status,
		})

	return result.RowsAffected, result.Error
}
