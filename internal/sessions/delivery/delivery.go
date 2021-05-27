package delivery

import (
	"context"
	"time"

	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/proto"
	"google.golang.org/grpc"
)

// AuthClient структура
type AuthClient struct {
	client proto.AuthHandlerClient
}

// NewAuthClient инициализация структуры
func NewAuthClient(conn *grpc.ClientConn) *AuthClient {
	client := proto.NewAuthHandlerClient(conn)
	return &AuthClient{client: client}
}

// Create создать пользователя
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

// GetUser получить пользователя
func (ac *AuthClient) GetUser(sessionID string) (string, error) {
	sessionValue := &proto.SessionValue{SessionID: sessionID}
	userID, err := ac.client.GetUser(context.Background(), sessionValue)
	if err != nil {
		return "", err
	}
	return userID.UserID, nil
}

// Delete удалить пользователя
func (ac *AuthClient) Delete(sessionID string) error {
	sessionValue := &proto.SessionValue{SessionID: sessionID}
	_, err := ac.client.Delete(context.Background(), sessionValue)
	return err
}
