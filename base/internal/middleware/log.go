package middleware

import (
	"time"

	"github.com/blackflagsoftware/forge-go/base/config"
	"github.com/sirupsen/logrus"
)

var Default = logrus.New()

func init() {
	Default.SetOutput(config.LogOutput)
	Default.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
	})
}
