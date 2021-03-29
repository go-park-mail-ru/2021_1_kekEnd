package actors

import (
	"github.com/gin-gonic/gin"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/actors"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/models"
	_const "github.com/go-park-mail-ru/2021_1_kekEnd/pkg/const"
)

type Handler struct {
	useCase actors.UseCase
}

func NewHandler(useCase actors.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) CreateActor(ctx *gin.Context) {
	actor, _ := ctx.Get(_const.ActorKey)
	actorModel, _ := actor.(models.Actor)


	_ = h.useCase.CreateActor(actorModel)
}

func (h *Handler) GetActor(ctx *gin.Context) {
	id := ctx.Param("actor_id")

	h.useCase.GetActor(id)
}

func (h *Handler) EditActor(ctx *gin.Context) {
	change, _ := ctx.Get(_const.ActorKey)
	changedModel, _ := change.(models.Actor)


	_ = h.useCase.EditActor(changedModel)
}