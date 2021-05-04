package grpc

import (
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/logger"
	proto "github.com/go-park-mail-ru/2021_1_kekEnd/internal/proto"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/services/sessions"
	"google.golang.org/grpc/status"
	"time"
)

type Delivery struct {
	UseCase sessions.UseCase
	Log     *logger.Logger
}

func NewDelivery(uc sessions.UseCase, Log *logger.Logger) *Delivery {
	return &Delivery{
		UseCase: uc,
		Log:     Log,
	}
}

func (d *Delivery) Create(session *proto.CreateSession) (*proto.SessionValue, error) {
	expires := time.Duration(session.Expires.AsTime().UnixNano())
	sessionID, err := d.UseCase.Create(session.UserID, expires)
	if err != nil {
		return nil, status.Error(500, "error in CreateSession")
	}
	protoSessionValue := &proto.SessionValue{SessionID: sessionID}
	return protoSessionValue, nil
}

func (d *Delivery) GetUser(sessionID *proto.SessionValue) (*proto.UserValue, error) {
	userID, err := d.UseCase.Check(sessionID.SessionID)
	if err != nil {
		return nil, status.Error(500, "error in GetUser")
	}
	protoUserValue := &proto.UserValue{UserID: userID}
	return protoUserValue, nil
}

func (d *Delivery) Delete(sessionID *proto.SessionValue) error {
	err := d.UseCase.Delete(sessionID.SessionID)
	if err != nil {
		return status.Error(500, "error in Delete")
	}
	return nil
}
