package fileserver

import "github.com/go-park-mail-ru/2021_1_kekEnd/internal/proto"

//go:generate mockgen -destination=mocks/delivery_mock.go -package=mocks . Delivery
type Delivery interface {
	Upload(stream proto.FileServerService_UploadServer, category string) error
}
