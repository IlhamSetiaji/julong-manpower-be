package utils

import (
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-manpower-be/internal/http/response"
)

var ResponseChannel = make(chan map[string]interface{}, 100)

var Rchans = make(map[string](chan response.RabbitMQResponse))

type RabbitMsg struct {
	QueueName string                  `json:"queueName"`
	Message   request.RabbitMQRequest `json:"message"`
}

// channel to publish rabbit messages
var Pchan = make(chan RabbitMsg, 10)
