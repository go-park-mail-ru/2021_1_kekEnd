package grpc

import (
	"context"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/proto"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/services/sessions"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

type AuthHandlerServer struct {
	UseCase sessions.UseCase
}

func NewAuthHandlerServer(uc sessions.UseCase) *AuthHandlerServer {
	return &AuthHandlerServer{
		UseCase: uc,
	}
}

func (d *AuthHandlerServer) Create(ctx context.Context, session *proto.CreateSession) (*proto.SessionValue, error) {
	sessionID, err := d.UseCase.Create(session.UserID, time.Duration(session.Expires))
	if err != nil {
		return nil, status.Error(500, "error in CreateSession")
	}
	protoSessionValue := &proto.SessionValue{SessionID: sessionID}
	return protoSessionValue, nil
}

func (d *AuthHandlerServer) GetUser(ctx context.Context, sessionID *proto.SessionValue) (*proto.UserValue, error) {
	userID, err := d.UseCase.Check(sessionID.SessionID)
	if err != nil {
		return nil, status.Error(500, "error in GetUser")
	}
	protoUserValue := &proto.UserValue{UserID: userID}
	return protoUserValue, nil
}

func (d *AuthHandlerServer) Delete(ctx context.Context, sessionID *proto.SessionValue) (*empty.Empty, error) {
	err := d.UseCase.Delete(sessionID.SessionID)
	if err != nil {
		return &emptypb.Empty{}, status.Error(500, "error in Delete")
	}
	return &emptypb.Empty{}, nil
}
