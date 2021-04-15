package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/logger"
	"net/http"
)

type Csrf interface {
	CheckCSRF() gin.HandlerFunc
}

type CsrfMiddleware struct {
	Log *logger.Logger
}

func NewCsrfMiddleware(Log *logger.Logger) *CsrfMiddleware {
	return &CsrfMiddleware{
		Log: Log,
	}
}

func (m *CsrfMiddleware) CheckCSRF() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		Token := ctx.GetHeader("X-CSRF-Token")
		Cookie, err := ctx.Cookie("X-CSRF-Cookie")
		if err != nil {
			msg := "No csrf cookie in request" + err.Error()
			m.Log.LogWarning(ctx, "CsrfMiddleware", "CheckCSRF", msg)
			ctx.Status(http.StatusBadRequest) // 400
			return
		}

		if Token != Cookie {
			msg := "Csrf-Cookie doesn't match Csrf-Token"
			m.Log.LogWarning(ctx, "CsrfMiddleware", "CheckCSRF", msg)
			ctx.Status(http.StatusForbidden) // 403
			return
		}

		ctx.Next()
	}
}
