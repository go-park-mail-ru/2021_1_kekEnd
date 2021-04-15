package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/logger"
	"math/rand"
	"time"
)

type AccessLog interface {
	AccessLogMiddleware(log *logger.Logger) gin.HandlerFunc
}

func AccessLogMiddleware(log *logger.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rand.Seed(time.Now().UnixNano())
		rid := fmt.Sprintf("%016x", rand.Int())[:5]

		log.StartReq(*ctx.Request, rid)
		start := time.Now()

		ctx.Next()
		log.EndReq(*ctx.Request, start, rid)
	}
}
