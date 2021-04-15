package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/logger"
)

type Csrf interface {
	CheckCsrf() gin.HandlerFunc
}

type CsrfMiddleware struct {
	Log *logger.Logger
}

func NewCsrfMiddleware(Log *logger.Logger) *CsrfMiddleware {
	return &CsrfMiddleware{
		Log: Log,
	}
}

func (m *CsrfMiddleware) CheckCsrf() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		Token := ctx.GetHeader("X-CSRF-Token")
		Cookie, err := ctx.Cookie("X-CSRF-Cookie")
		if err != nil {
			msg := "No csrf cookie in request" + err.Error()
			m.Log.LogWarning(ctx, "CsrfMiddleware", "CheckCSRF", msg)
			return
		}

		if Token != Cookie {
			msg := "Csrf-Cookie doesn't match Csrf-Token"
			m.Log.LogWarning(ctx, "CsrfMiddleware", "CheckCSRF", msg)
			return
		}

		ctx.Next()
	}
}
