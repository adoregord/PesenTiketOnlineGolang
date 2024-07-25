package handler

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

func LogMethod(msg string, method string, waktu time.Time, httpstatus int) {
	log.Info().
		Int("httpStatusCode", httpstatus).
		Str("StatusDescription", http.StatusText(httpstatus)).
		TimeDiff("ProcessTime", time.Now(), waktu).
		Str("httpMethod", method).
		Msg(msg)
}