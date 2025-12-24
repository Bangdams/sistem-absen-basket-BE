package converter

import (
	"absen-qr-backend/internal/entity"
	"absen-qr-backend/internal/model"
	"log"
)

func UserToResponse(user *entity.User) *model.UserResponse {
	log.Println("log from user to response")

	reponse := &model.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format("2006-01-02"),
	}

	switch user.Role {
	case "coach":
		reponse.CoachRequest = &model.CoachRequest{
			Nip:      user.Coach.Nip,
			FullName: user.Coach.FullName,
		}
	case "student":
		reponse.StudentRequest = &model.StudentRequest{
			Nis:         user.Student.Nis,
			FullName:    user.Student.FullName,
			Address:     user.Student.Address,
			PhoneNumber: user.Student.PhoneNumber,
		}
	}

	return reponse
}

func UserToResponseForUpdate(user *entity.User) *model.UserResponse {
	log.Println("log from UserToResponseForUpdate")

	return &model.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format("2006-01-02"),
	}
}

func LoginUserToResponse(user *entity.User) *model.LoginResponse {
	log.Println("log from login user to response")

	return &model.LoginResponse{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
	}
}

func UserToResponses(users *[]entity.User) *[]model.UserResponse {
	var userResponses []model.UserResponse

	log.Println("log from user to responses")

	for _, user := range *users {
		userResponses = append(userResponses, *UserToResponse(&user))
	}

	return &userResponses
}
