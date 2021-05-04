package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/services/sessions"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/users"
	_const "github.com/go-park-mail-ru/2021_1_kekEnd/pkg/const"
	"net/http"
)

func respondWithError(ctx *gin.Context, code int, message interface{}) {
	ctx.AbortWithStatusJSON(code, gin.H{"error": message})
}

type Auth interface {
	RequireAuth() gin.HandlerFunc
	CheckAuth() gin.HandlerFunc
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

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionID, err := ctx.Cookie("session_id")
		if err != nil {
			fmt.Println("no sessions_id in request", err)
			respondWithError(ctx, http.StatusUnauthorized, "no sessions_id in request") //401
			return
		}

		username, err := m.sessions.GetUser(sessionID)
		if err != nil {
			fmt.Println("no sessions for this user", err)
			respondWithError(ctx, http.StatusUnauthorized, "no sessions for this user") //401
			return
		}

		user, err := m.useCase.GetUser(username)
		if err != nil {
			respondWithError(ctx, http.StatusInternalServerError, "no user with this username") //500
			return
		}

		ctx.Set(_const.UserKey, *user)
		ctx.Next()
	}
}

func (m *AuthMiddleware) CheckAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionID, err := ctx.Cookie("session_id")
		if err != nil {
			ctx.Set(_const.AuthStatusKey, false)
			return
		}

		username, err := m.sessions.GetUser(sessionID)
		if err != nil {
			ctx.Set(_const.AuthStatusKey, false)
			return
		}

		user, err := m.useCase.GetUser(username)
		if err != nil {
			ctx.Set(_const.AuthStatusKey, false)
			return
		}

		ctx.Set(_const.UserKey, *user)
		ctx.Set(_const.AuthStatusKey, true)
		ctx.Next()
	}
}
