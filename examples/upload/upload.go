package main

import (
	"context"
	"io"
	"os"

	"github.com/micro/cli"
	"github.com/micro/go-log"
	"github.com/micro/go-micro"

	pb "github.com/microhq/image-srv/proto/image"
)

var (
	albumName string
	imgPath   string
	buffSize  int
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
				Name:  "img_path",
				Value: "",
				Usage: "Path to image file",
			},
			cli.IntFlag{
				Name:  "buf_size",
				Value: 1024 * 1024,
				Usage: "Download buffer size",
			},
		),
	)

	// Initialise service
	service.Init(
		micro.Action(func(c *cli.Context) {
			albumName = c.String("album_name")
			imgPath = c.String("img_path")
			buffSize = c.Int("buf_size")
		}),
	)

	// Initialize service client
	client := pb.NewImageService("go.micro.srv.image", service.Client())

	// Create bucket with given bucket ID
	if _, err := client.CreateAlbum(context.Background(), &pb.CreateAlbumReq{Id: albumName}); err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(imgPath)
	if err != nil {
		log.Logf("failed to open file %s", imgPath)
		return
	}
	defer file.Close()

	// request stream object to upload the blob via
	stream, err := client.Upload(context.Background())
	if err != nil {
		log.Logf("failed to upload file %s", imgPath)
		return
	}
	defer stream.Close()

	// flag to control upload
	upload := true

	// sent data in 1M chunks
	buf := make([]byte, buffSize)

	// stream the data across to blob service
	for upload {
		n, err := file.Read(buf)
		if err == io.EOF {
			upload = false
			break
		}

		if err != nil {
			log.Logf("error reading file %s: %s", imgPath, err)
			break
		}

		log.Logf("streaming %d bytes for image: %s", n, imgPath)

		if err := stream.Send(&pb.UploadImageReq{Id: imgPath, AlbumId: albumName, Data: buf[:n]}); err != nil {
			log.Logf("error streaming image %s: %s", imgPath, err)
			break
		}
	}
}
