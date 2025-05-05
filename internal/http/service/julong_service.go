package service

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/request"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IJulongService interface {
	CreateJulongNotification(req *request.CreateNotificationRequest) error
}

type JulongService struct {
	Viper *viper.Viper
	Log   *logrus.Logger
}

func NewJulongService(viper *viper.Viper, log *logrus.Logger) IJulongService {
	return &JulongService{
		Viper: viper,
		Log:   log,
	}
}

func (s *JulongService) CreateJulongNotification(payload *request.CreateNotificationRequest) error {
	url := s.Viper.GetString("notification.url") + "/api/v1/notifications"
	method := "POST"

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		s.Log.Error(err)
		return errors.New("[JulongService.CreateJulongNotification] Error when marshalling payload: " + err.Error())
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
		},
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		s.Log.Error(err)
		return errors.New("[JulongService.CreateJulongNotification] Error when creating request: " + err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", "JulongService/1.0")

	res, err := client.Do(req)
	if err != nil {
		s.Log.Error(err)
		return errors.New("[JulongService.CreateJulongNotification] Error when sending request: " + err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		s.Log.Error(err)
		return errors.New("[JulongService.CreateJulongNotification] Error when sending request: " + string(bodyBytes))
	}

	return nil
}

func JulongServiceFactory(viper *viper.Viper, log *logrus.Logger) IJulongService {
	return NewJulongService(viper, log)
}
