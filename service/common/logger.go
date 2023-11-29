package common

import (
	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

func InitLog() {
	Logger.SetFormatter(&logrus.JSONFormatter{})
	Logger.SetLevel(logrus.TraceLevel)
}
