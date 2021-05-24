package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/logger"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/search"
)

func RegisterHttpEndpoints(router *gin.RouterGroup, searchUC search.UseCase, Log *logger.Logger) {
	handler := NewHandler(searchUC, Log)
	router.GET("/search", handler.Search)
}
