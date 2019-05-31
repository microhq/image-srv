package handler

import (
	"context"
	"fmt"
	"io"

	"github.com/micro/go-log"
	blob "github.com/microhq/blob-srv/proto/blob"
	pb "github.com/microhq/image-srv/proto/image"
)

type Image struct {
	Client blob.BlobService
}

func (h *Image) CreateAlbum(ctx context.Context, req *pb.CreateAlbumReq, resp *pb.CreateAlbumResp) error {
	log.Logf("Received request to create new album")

	// Create bucket with album name
	if _, err := h.Client.CreateBucket(ctx, &blob.CreateBucketReq{Id: req.Id}); err != nil {
		log.Logf("failed to create album bucket: %s", err)
		return fmt.Errorf("failed to create new album: %s", req.Id)
	}

	return nil
}

func (h *Image) DeleteAlbum(ctx context.Context, req *pb.DeleteAlbumReq, resp *pb.DeleteAlbumResp) error {
	log.Logf("Received request to delete existing album")

	// Delete bucket with album name
	if _, err := h.Client.DeleteBucket(ctx, &blob.DeleteBucketReq{Id: req.Id}); err != nil {
		log.Logf("failed to delete album bucket: %s", err)
		return fmt.Errorf("failed to delete album: %s", req.Id)
	}

	return nil
}

func (h *Image) Upload(ctx context.Context, stream pb.Image_UploadStream) error {
	log.Logf("Received request to upload image")

	var id string
	var albumId string
	defer stream.Close()

	// request stream object to upload the blob via
	blobStream, err := h.Client.Put(ctx)
	if err != nil {
		log.Logf("failed to establish stream connection to blob.Service")
		return fmt.Errorf("failed to upload image: %s", err)
	}
	defer blobStream.Close()

	for {
		img, err := stream.Recv()
		if err == io.EOF {
			log.Logf("finished receiving image data")
			break
		}

		if err != nil {
			log.Logf("error storing image: %s", err)
			return err
		}

		albumId = img.AlbumId
		id = img.Id

		log.Logf("Received %d bytes of image: %s", len(img.Data))

		if err := blobStream.Send(&blob.PutReq{Id: id, BucketId: albumId, Data: img.Data}); err != nil {
			log.Logf("error streaming blob %s/%s: %s", albumId, id, err)
			break
		}
	}

	return nil
}

func (h *Image) Download(ctx context.Context, req *pb.DownloadImageReq, stream pb.Image_DownloadStream) error {
	log.Logf("Received request to download image")

	return nil
}

func (h *Image) Delete(context.Context, *pb.DeleteImageReq, *pb.DeleteImageResp) error {
	log.Logf("Received request to delete image")

	return nil
}

func (h *Image) Share(context.Context, *pb.ShareImageReq, *pb.ShareImageResp) error {
	log.Logf("Received request to share image")

	return nil
}

func (h *Image) Unshare(context.Context, *pb.UnshareImageReq, *pb.UnshareImageResp) error {
	log.Logf("Received request to unshare image")

	return nil
}
