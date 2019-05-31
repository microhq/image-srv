package main

import (
	"github.com/micro/go-log"
	"github.com/micro/go-micro"

	blob "github.com/microhq/blob-srv/proto/blob"
	"github.com/microhq/image-srv/handler"
	pb "github.com/microhq/image-srv/proto/image"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.srv.image"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	pb.RegisterImageHandler(service.Server(), &handler.Image{
		Client: blob.NewBlobService("go.micro.srv.blob", service.Client()),
	})

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
