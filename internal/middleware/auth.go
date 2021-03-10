package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/sessions"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/users"
	"net/http"
)

const userKey = "user"

func respondWithError(ctx *gin.Context, code int, message interface{}) {
	ctx.AbortWithStatusJSON(code, gin.H{"error": message})
}

type Auth interface {
	CheckAuth() gin.HandlerFunc
}

type AuthMiddleware struct {
	useCase  users.UseCase
	sessions sessions.Delivery
}

func NewAuthMiddleware(useCase users.UseCase, sessions sessions.Delivery) *AuthMiddleware {
	return &AuthMiddleware{
		useCase: useCase,
		sessions: sessions,
	}
}

func (m *AuthMiddleware) CheckAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionID, err := ctx.Cookie("session_id")
		if err != nil {
			respondWithError(ctx, http.StatusUnauthorized, "no sessions_id in request") //401
			return
		}

		username, err := m.sessions.GetUser(sessionID)
		if err != nil {
			respondWithError(ctx, http.StatusUnauthorized, "no sessions for this user") //401
			return
		}

		user, err := m.useCase.GetUser(username)
		if err != nil {
			respondWithError(ctx, http.StatusInternalServerError, "no user with this username") //500
			return
		}

		ctx.Set(userKey, *user)
		ctx.Next()
	}
}
