package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/sessions"
	usersMock "github.com/go-park-mail-ru/2021_1_kekEnd/internal/users/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckAuth(t *testing.T) {
	testErr := errors.New("error no cookie")
	const userKey = "user"

	t.Run("GetUser-OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sessionsDelivery := sessions.NewMockDelivery(ctrl)
		userUseCase := usersMock.NewMockUseCase(ctrl)
		mdw := NewAuthMiddleware(userUseCase, sessionsDelivery)

		Cookie := http.Cookie{
			Name:  "session_id",
			Value: "1aLetMeIn7e7fa6d",
		}
		username := "tester"
		userModel := models.User{
			Username: username,
		}

		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Request, _ = http.NewRequest("POST", "/", nil)
		ctx.Request.AddCookie(&Cookie)

		sessionsDelivery.
			EXPECT().
			GetUser(Cookie.Value).
			Return(username, nil)

		userUseCase.EXPECT().GetUser(username).Return(&userModel, nil)

		handler := mdw.CheckAuth(true)
		handler(ctx)

		userFromMiddleware, _ := ctx.Get(userKey)
		userModelFromMiddleware := userFromMiddleware.(models.User)

		assert.Equal(t, userModel, userModelFromMiddleware) // 500
	})

	t.Run("GetUser-No-cookie", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sessionsDelivery := sessions.NewMockDelivery(ctrl)
		userUseCase := usersMock.NewMockUseCase(ctrl)
		mdw := NewAuthMiddleware(userUseCase, sessionsDelivery)

		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Request, _ = http.NewRequest("POST", "/", nil)

		handler := mdw.CheckAuth(true)
		handler(ctx)

		assert.Equal(t, http.StatusUnauthorized, ctx.Writer.Status()) // 401
	})

	t.Run("GetUser-No-Username", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sessionsDelivery := sessions.NewMockDelivery(ctrl)
		userUseCase := usersMock.NewMockUseCase(ctrl)
		mdw := NewAuthMiddleware(userUseCase, sessionsDelivery)

		Cookie := http.Cookie{
			Name:  "session_id",
			Value: "1aLetMeIn7e7fa6d",
		}
		username := "tester"
		userModel := models.User{
			Username: username,
		}

		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Request, _ = http.NewRequest("POST", "/", nil)
		ctx.Request.AddCookie(&Cookie)

		sessionsDelivery.
			EXPECT().
			GetUser(Cookie.Value).
			Return(username, nil)

		userUseCase.EXPECT().GetUser(username).Return(&userModel, testErr)

		handler := mdw.CheckAuth(true)
		handler(ctx)

		assert.Equal(t, http.StatusInternalServerError, ctx.Writer.Status()) // 500
	})

	t.Run("GetUser-No-Session", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sessionsDelivery := sessions.NewMockDelivery(ctrl)
		userUseCase := usersMock.NewMockUseCase(ctrl)
		mdw := NewAuthMiddleware(userUseCase, sessionsDelivery)

		Cookie := http.Cookie{
			Name:  "session_id",
			Value: "1aLetMeIn7e7fa6d",
		}

		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Request, _ = http.NewRequest("POST", "/", nil)
		ctx.Request.AddCookie(&Cookie)

		sessionsDelivery.
			EXPECT().
			GetUser(Cookie.Value).
			Return("", testErr)

		handler := mdw.CheckAuth(true)
		handler(ctx)

		assert.Equal(t, http.StatusUnauthorized, ctx.Writer.Status()) // 401
	})
}
