package service

import (
	"errors"

	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/messaging"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/request"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type INotificationService interface {
	CreatePeriodNotification(createdBy string) error
}

type NotificationService struct {
	Viper         *viper.Viper
	Log           *logrus.Logger
	UserMessage   messaging.IUserMessage
	JulongService IJulongService
}

func NewNotificationService(viper *viper.Viper, log *logrus.Logger, userMessage messaging.IUserMessage, julongService IJulongService) INotificationService {
	return &NotificationService{
		Viper:         viper,
		Log:           log,
		UserMessage:   userMessage,
		JulongService: julongService,
	}
}

func (s *NotificationService) CreatePeriodNotification(createdBy string) error {
	userIDs, err := s.UserMessage.SendGetUserIDsByPermissionNames([]string{"create-mpp"})
	if err != nil {
		s.Log.Error(err)
		return err
	}

	if len(userIDs) == 0 {
		s.Log.Error("No user IDs found with the specified permission names")
		return errors.New("no user IDs found with the specified permission names")
	}

	payload := &request.CreateNotificationRequest{
		Application: "MANPOWER",
		Name:        "Period - Open Period",
		URL:         "/d/location",
		Message:     "Manpower Planning period has officially started. Please fill in and complete your data. Ensure all required information is accurate before the deadline to avoid any issues.",
		UserIDs:     userIDs,
		CreatedBy:   createdBy,
	}

	err = s.JulongService.CreateJulongNotification(payload)
	if err != nil {
		s.Log.Error(err)
		return err
	}

	return nil
}

func NotificationServiceFactory(viper *viper.Viper, log *logrus.Logger) INotificationService {
	userMessage := messaging.UserMessageFactory(log)
	julongService := JulongServiceFactory(viper, log)
	return NewNotificationService(viper, log, userMessage, julongService)
}
