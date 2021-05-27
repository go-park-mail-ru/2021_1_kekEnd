package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/logger"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/search"
)

// Handler структура хендлера поиска
type Handler struct {
	useCase search.UseCase
	Log     *logger.Logger
}

// NewHandler инициализация хендлера поиска
func NewHandler(useCase search.UseCase, Log *logger.Logger) *Handler {
	return &Handler{
		useCase: useCase,
		Log:     Log,
	}
}

// Search поиск
func (h *Handler) Search(ctx *gin.Context) {
	query := ctx.Query("q")
	searchResults, err := h.useCase.Search(query)
	if err != nil {
		h.Log.LogError(ctx, "search", "Search", err)
		ctx.AbortWithStatus(http.StatusInternalServerError) // 500
		return
	}

	ctx.JSON(http.StatusOK, searchResults)
}
