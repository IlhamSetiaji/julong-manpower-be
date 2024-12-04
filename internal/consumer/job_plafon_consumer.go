package consumer

import "github.com/sirupsen/logrus"

type IJobPlafonConsumer interface {
	ConsumeCheckJobExistMessage()
}

type JobPlafonConsumer struct {
	Log *logrus.Logger
}

func NewJobPlafonConsumer(log *logrus.Logger) IJobPlafonConsumer {
	return &JobPlafonConsumer{
		Log: log,
	}
}

func (c *JobPlafonConsumer) ConsumeCheckJobExistMessage() {
	
}
