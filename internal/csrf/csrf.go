package csrf

import (
	"github.com/gin-gonic/gin"
	_const "github.com/go-park-mail-ru/2021_1_kekEnd/pkg/const"
	uuid "github.com/satori/go.uuid"
)

func CreateCsrfToken(ctx *gin.Context) {
	csrfToken := uuid.NewV4().String()

	ctx.Header("X-CSRF-Token", csrfToken)
	ctx.SetCookie("X-CSRF-Cookie",
		csrfToken,
		int(_const.CsrfExpires),
		"/",
		_const.Host,
		false,
		true,
	)
}
