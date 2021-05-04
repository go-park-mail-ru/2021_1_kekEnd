package usecase

import (
	"context"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/proto"
	"google.golang.org/grpc"
	"time"
)

type AuthClient struct {
	client proto.AuthHandlerClient
}

func NewAuthClient(conn *grpc.ClientConn) *AuthClient {
	client := proto.NewAuthHandlerClient(conn)
	return &AuthClient{client: client}
}

func (ac *AuthClient) Create(userID string, expires time.Duration) (string, error) {
	session := &proto.CreateSession{
		UserID:  userID,
		Expires: int64(expires),
	}
	value, err := ac.client.Create(context.Background(), session)
	if err != nil {
		return "", err
	}
	return value.SessionID, nil
}

func (ac *AuthClient) Check(sessionID string) (string, error) {
	sessionValue := &proto.SessionValue{SessionID: sessionID}
	userID, err := ac.client.GetUser(context.Background(), sessionValue)
	if err != nil {
		return "", err
	}
	return userID.UserID, nil
}

func (ac *AuthClient) Delete(sessionID string) error {
	sessionValue := &proto.SessionValue{SessionID: sessionID}
	_, err := ac.client.Delete(context.Background(), sessionValue)
	return err
}
