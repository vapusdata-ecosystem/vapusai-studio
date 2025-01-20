package pkg

import (
	zerolog "github.com/rs/zerolog"
	dmlogger "github.com/vapusdata-oss/aistudio/core/logger"
)

var DmLogger zerolog.Logger

func InitLogger(debugMode bool) {
	DmLogger = dmlogger.NewZeroLogger(debugMode, true, false, "", "caller", "level", "fields", "time")
}
