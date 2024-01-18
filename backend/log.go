package app

import (
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetReportCaller(true)
	log.SetLevel(log.InfoLevel)
}
