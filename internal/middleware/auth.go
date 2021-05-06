package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/sessions"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/users"
	_const "github.com/go-park-mail-ru/2021_1_kekEnd/pkg/const"
	"net/http"
)

func respondWithError(ctx *gin.Context, code int, message interface{}) {
	ctx.AbortWithStatusJSON(code, gin.H{"error": message})
}

type Auth interface {
	CheckAuth(isRequired bool) gin.HandlerFunc
}

type AuthMiddleware struct {
	useCase  users.UseCase
	sessions sessions.Delivery
}

func NewAuthMiddleware(useCase users.UseCase, sessions sessions.Delivery) *AuthMiddleware {
	return &AuthMiddleware{
		useCase:  useCase,
		sessions: sessions,
	}
}

func (m *AuthMiddleware) CheckAuth(isRequired bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionID, err := ctx.Cookie("session_id")
		if err != nil {
			if isRequired {
				fmt.Println("no sessions_id in request", err)
				respondWithError(ctx, http.StatusUnauthorized, "no sessions_id in request") //401
				return
			}
			ctx.Set(_const.AuthStatusKey, false)
			ctx.Next()
			return
		}

		username, err := m.sessions.GetUser(sessionID)
		if err != nil {
			if isRequired {
				fmt.Println("no sessions for this user", err)
				respondWithError(ctx, http.StatusUnauthorized, "no sessions for this user") //401
				return
			}
			ctx.Set(_const.AuthStatusKey, false)
			ctx.Next()
			return
		}

		user, err := m.useCase.GetUser(username)
		if err != nil {
			if isRequired {
				respondWithError(ctx, http.StatusInternalServerError, "no user with this username") //500
				return
			}
			ctx.Set(_const.AuthStatusKey, false)
			ctx.Next()
			return
		}

		ctx.Set(_const.UserKey, *user)
		ctx.Set(_const.AuthStatusKey, true)
		ctx.Next()
	}
}
