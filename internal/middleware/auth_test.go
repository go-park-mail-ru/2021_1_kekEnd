package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/sessions"
	sessionsDel "github.com/go-park-mail-ru/2021_1_kekEnd/internal/sessions/delivery"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/users/usecase"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestCheckAuth(t *testing.T) {
	//testErr := errors.New("error no cookie")

	t.Run("Check-Error-nouser", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mUC := sessions.NewMockUseCase(ctrl)
		sessionsDelivery := sessionsDel.NewDelivery(mUC)
		userUseCase := usecase.UsersUseCase{}
		mdw := NewAuthMiddleware(&userUseCase, sessionsDelivery)

		Request :=  new(http.Request)
		Cookie := http.Cookie{
			Name: "session_id",
			Value: "1aLetMeIn7e7fa6d",
		}
		Request.AddCookie(&Cookie)
		ctx := gin.Context{
			Request:  Request,
		}

		handler := mdw.CheckAuth()
		handler(&ctx)

		assert.Equal(t, http.StatusInternalServerError, ctx.Status) // 500

	})

	t.Run("Check-Error-nocookie", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mUC := sessions.NewMockUseCase(ctrl)
		sessionsDelivery := sessionsDel.NewDelivery(mUC)
		userUseCase := usecase.UsersUseCase{}
		mdw := NewAuthMiddleware(&userUseCase, sessionsDelivery)

		Request :=  new(http.Request)
		ctx := new(gin.Context)
		ctx.Request = Request


		handler := mdw.CheckAuth()
		handler(ctx)

		assert.Equal(t, http.StatusUnauthorized, ctx.Status) // 401
	})
}
