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

	// request stream object to upload the blob via
	blobStream, err := h.Client.Put(ctx)
	if err != nil {
		log.Logf("failed to start stream with blob.Service")
		return fmt.Errorf("failed to upload image: %s", err)
	}

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

		log.Logf("received %d bytes of image: %s", len(img.Data), id)

		if err := blobStream.Send(&blob.PutReq{Id: id, BucketId: albumId, Data: img.Data}); err != nil {
			log.Logf("error streaming blob %s/%s: %s", albumId, id, err)
			break
		}
	}

	log.Logf("closing stream socket")

	// close stream when done sending
	if err := blobStream.Close(); err != nil {
		log.Logf("error closing stream socket: %s", err)
		return err
	}

	return nil
}

func (h *Image) Download(ctx context.Context, req *pb.DownloadImageReq, stream pb.Image_DownloadStream) error {
	log.Logf("Received request to download image")

	// request stream to use for blob download
	blobStream, err := h.Client.Get(ctx, &blob.GetReq{Id: req.Id, BucketId: req.AlbumId})
	if err != nil {
		log.Logf("failed to get download file: %s", err)
		return fmt.Errorf("failed to download image: %s", err)
	}

	for {
		img, err := blobStream.Recv()
		if err == io.EOF {
			log.Logf("finished receiving image data")
			break
		}

		if err != nil {
			log.Logf("error retrieving image %s: %s", req.Id, err)
			return err
		}

		log.Logf("received %d bytes of image: %s", len(img.Data), req.Id)

		if err := stream.Send(&pb.DownloadImageResp{Data: img.Data}); err != nil {
			log.Logf("error streaming image %s/%s: %s", req.AlbumId, req.Id, err)
			break
		}
	}

	log.Logf("closing stream socket")

	// close stream when done sending
	if err := stream.Close(); err != nil {
		log.Logf("failed to close stream socket: %s", err)
		return err
	}

	return nil
}

func (h *Image) Delete(ctx context.Context, req *pb.DeleteImageReq, resp *pb.DeleteImageResp) error {
	log.Logf("Received request to delete image")

	// request stream to use for blob download
	_, err := h.Client.Delete(ctx, &blob.DeleteReq{Id: req.Id, BucketId: req.AlbumId})
	if err != nil {
		log.Logf("failed to delete image %s/%s: %s", req.Id, req.AlbumId, err)
		return fmt.Errorf("failed to delete image %s: %s", req.Id, err)
	}

	return nil
}
