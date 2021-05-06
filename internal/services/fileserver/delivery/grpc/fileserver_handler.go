package grpc

import (
	"fmt"
	"github.com/go-park-mail-ru/2021_1_kekEnd/internal/proto"
	"google.golang.org/grpc/metadata"
	"io"
	"io/ioutil"
	"log"
)

type FileServerHandlerServer struct {
}

func NewFileServerHandlerServer() *FileServerHandlerServer {
	return &FileServerHandlerServer{}
}

func (dl *FileServerHandlerServer) Upload(stream proto.FileServerHandler_UploadServer) error {
	data := make([]byte, 0, 1024)
	meta, _ := metadata.FromIncomingContext(stream.Context())
	fileName := meta.Get("fileName")[0]

	for {
		inData, err := stream.Recv()
		if err == io.EOF {
			out := &proto.UploadStatus{
				Message: "OK",
				Code:    proto.StatusCode_SUCCESS,
			}
			log.Println("Everything sent, size:", len(data))
			err := stream.SendAndClose(out)
			if err != nil {
				log.Println(err)
			}
			break
		}
		if err != nil {
			return err
		}

		data = append(data, inData.Content...)
	}
	err := ioutil.WriteFile(fmt.Sprintf(fileName), data, 0666)
	if err != nil {
		log.Println("Error while saving file: ", err)
		return err
	}
	return nil
}
