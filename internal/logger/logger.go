package logger

import (
	"context"
	"github.com/gin-gonic/gin"
	_const "github.com/go-park-mail-ru/2021_1_kekEnd/pkg/const"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Logger struct {
	*logrus.Logger
}

func NewAccessLogger() *Logger {
	lg := logrus.New()
	return &Logger{lg}
}

func (l *Logger) GetIdFromContext(ctx context.Context) string {
	ridV := ctx.Value(_const.RequestID)
	rid, ok := ridV.(string)
	if !ok {
		l.WithFields(logrus.Fields{
			"id":       "NO_ID",
			"package":  "logger",
			"function": "GetIdFromContext",
		}).Warn("can't get request id from context")
		return ""
	}
	return rid
}

func (l *Logger) StartReq(r http.Request, rid string) {
	l.WithFields(logrus.Fields{
		"id":        rid,
		"usr_addr":  r.RemoteAddr,
		"req_URI":   r.RequestURI,
		"method":    r.Method,
		"usr_agent": r.UserAgent(),
	}).Info("request started")
}

func (l *Logger) EndReq(r http.Request, start time.Time, rid string) {
	l.WithFields(logrus.Fields{
		"id":        rid,
		"usr_addr":  r.RemoteAddr,
		"req_URI":   r.RequestURI,
		"method":    r.Method,
		"usr_agent": r.UserAgent(),
		"μs":        time.Since(start).Microseconds(),
	}).Info("request ended")
}

func (l *Logger) HttpInfo(ctx *gin.Context, msg string, status int) {
	l.WithFields(logrus.Fields{
		"id":     l.GetIdFromContext(ctx),
		"status": status,
	}).Info(msg)
}

func (l *Logger) LogWarning(ctx *gin.Context, pkg string, funcName string, msg string) {
	l.WithFields(logrus.Fields{
		"id":       l.GetIdFromContext(ctx),
		"package":  pkg,
		"function": funcName,
	}).Warn(msg)
}

func (l *Logger) LogError(ctx *gin.Context, pkg string, funcName string, err error) {
	l.WithFields(logrus.Fields{
		"id":       l.GetIdFromContext(ctx),
		"package":  pkg,
		"function": funcName,
	}).Error(err)
}
