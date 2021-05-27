package grpc

import (
	"context"
	"time"

	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/proto"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/services/sessions"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// AuthHandlerServer структура хендлера авторизации
type AuthHandlerServer struct {
	UseCase sessions.UseCase
}

// NewAuthHandlerServer инициализация структуры хендлера авторизации
func NewAuthHandlerServer(uc sessions.UseCase) *AuthHandlerServer {
	return &AuthHandlerServer{
		UseCase: uc,
	}
}

// Create создание сессии
func (d *AuthHandlerServer) Create(ctx context.Context, session *proto.CreateSession) (*proto.SessionValue, error) {
	sessionID, err := d.UseCase.Create(session.UserID, time.Duration(session.Expires))
	if err != nil {
		return nil, status.Error(500, "error in CreateSession")
	}
	protoSessionValue := &proto.SessionValue{SessionID: sessionID}
	return protoSessionValue, nil
}

// GetUser получение юзера
func (d *AuthHandlerServer) GetUser(ctx context.Context, sessionID *proto.SessionValue) (*proto.UserValue, error) {
	userID, err := d.UseCase.GetUser(sessionID.SessionID)
	if err != nil {
		return nil, status.Error(500, "error in GetUser")
	}
	protoUserValue := &proto.UserValue{UserID: userID}
	return protoUserValue, nil
}

// Delete удалени сессии
func (d *AuthHandlerServer) Delete(ctx context.Context, sessionID *proto.SessionValue) (*empty.Empty, error) {
	err := d.UseCase.Delete(sessionID.SessionID)
	if err != nil {
		return &emptypb.Empty{}, status.Error(500, "error in Delete")
	}
	return &emptypb.Empty{}, nil
}
