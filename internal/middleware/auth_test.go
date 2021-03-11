package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/sessions"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/users/usecase"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckAuth(t *testing.T) {
	//testErr := errors.New("error no cookie")
	const userKey = "user"

	t.Run("Check-OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sessionsDelivery := sessions.NewMockDelivery(ctrl)
		userUseCase := usecase.UsersUseCaseMock{}
		mdw := NewAuthMiddleware(&userUseCase, sessionsDelivery)

		Cookie := http.Cookie{
			Name: "session_id",
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

		userUseCase.On("GetUser", username).Return(&userModel, nil)

		handler := mdw.CheckAuth()
		handler(ctx)

		userFromMiddleware, _ := ctx.Get(userKey)
		userModelFromMiddleware := userFromMiddleware.(models.User)

		assert.Equal(t, userModel, userModelFromMiddleware) // 500
		//assert.Equal(t, 1, 1) // 500
	})

	//t.Run("Check-Error-nocookie", func(t *testing.T) {
	//	ctrl := gomock.NewController(t)
	//	defer ctrl.Finish()
	//
	//	mUC := sessions.NewMockUseCase(ctrl)
	//	sessionsDelivery := sessionsDel.NewDelivery(mUC)
	//	userUseCase := usecase.UsersUseCase{}
	//	mdw := NewAuthMiddleware(&userUseCase, sessionsDelivery)
	//
	//	Request :=  new(http.Request)
	//	ctx := new(gin.Context)
	//	ctx.Request = Request
	//
	//
	//	handler := mdw.CheckAuth()
	//	handler(ctx)
	//
	//	assert.Equal(t, http.StatusUnauthorized, ctx.Status) // 401
	//})
}
