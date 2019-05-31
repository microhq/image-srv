package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/micro/cli"
	"github.com/micro/go-log"
	"github.com/micro/go-micro"

	pb "github.com/microhq/image-srv/proto/image"
)

var (
	albumName string
	imgName   string
	imgPath   string
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.srv.image"),
		micro.Version("latest"),
		micro.Flags(
			cli.StringFlag{
				Name:  "album_name",
				Value: "",
				Usage: "Album name",
			},
			cli.StringFlag{
				Name:  "img_name",
				Value: "",
				Usage: "Image name",
			},
			cli.StringFlag{
				Name:  "img_path",
				Value: "",
				Usage: "Path to image file",
			},
		),
	)

	// Initialise service
	service.Init(
		micro.Action(func(c *cli.Context) {
			albumName = c.String("album_name")
			imgName = c.String("img_name")
			imgPath = c.String("img_path")
		}),
	)

	// Initialize service client
	client := pb.NewImageService("go.micro.srv.image", service.Client())

	file, err := os.OpenFile(imgPath, os.O_RDWR|os.O_CREATE, 0640)
	if err != nil {
		log.Logf("failed to open file %s: %s", imgPath, err)
		return
	}

	// request stream to use for blob download
	stream, err := client.Download(context.Background(), &pb.DownloadImageReq{Id: imgName, AlbumId: albumName})
	if err != nil {
		log.Logf("failed to download image: %s", err)
		return
	}

	// flag to control upload
	download := true

	// stream the data across to blob service
	for download {
		img, err := stream.Recv()
		if err == io.EOF {
			log.Logf("finished receiving data")
			break
		}

		if err != nil {
			fmt.Errorf("error downloading %s from album %s: %s", imgName, albumName, err)
			break
		}

		log.Logf("received %d bytes for image: %s", len(img.Data), imgName)

		if _, err := file.Write(img.Data); err != nil {
			fmt.Errorf("error writing data to file %s: %s", file.Name(), err)
			break
		}
	}

	if err := file.Close(); err != nil {
		log.Fatalf("Error closing file %s: %s", file.Name(), err)
	}
}
