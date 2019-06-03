package handler

import (
	"context"
	"fmt"
	"io"

	"github.com/micro/go-log"
	blob "github.com/microhq/blob-srv/proto/blob"
	pb "github.com/microhq/image-srv/proto/image"
)

// Image implements image service handler
type Image struct {
	// Blob is blob-srv client
	Blob blob.BlobService
}

// CreateAlbum creates new image album and returns it.
// It calls blob-srv which handles the low level filesystem operations.
// It returns error if the blob-srv fails to allocate dedicated filesystem for the album.
func (h *Image) CreateAlbum(ctx context.Context, req *pb.CreateAlbumReq, resp *pb.CreateAlbumResp) error {
	log.Logf("Received request to create new album")

	// Create bucket with album name
	if _, err := h.Blob.CreateBucket(ctx, &blob.CreateBucketReq{Id: req.Id}); err != nil {
		log.Logf("failed to create album bucket: %s", err)
		return fmt.Errorf("failed to create new album: %s", req.Id)
	}

	return nil
}

// DeleteAlbum deletes image album.
// It calls blob-srv to delete the album from the physical storage.
// It returns error if the album failed to be removed from storage.
func (h *Image) DeleteAlbum(ctx context.Context, req *pb.DeleteAlbumReq, resp *pb.DeleteAlbumResp) error {
	log.Logf("Received request to delete existing album")

	// Delete bucket with album name
	if _, err := h.Blob.DeleteBucket(ctx, &blob.DeleteBucketReq{Id: req.Id}); err != nil {
		log.Logf("failed to delete album bucket: %s", err)
		return fmt.Errorf("failed to delete album: %s", req.Id)
	}

	return nil
}

// Upload uploads image to remote server image album.
// It relies on blob-srv to store the image.
// It returns error if the image failed to be stored by blob-srv.
func (h *Image) Upload(ctx context.Context, stream pb.Image_UploadStream) error {
	log.Logf("Received request to upload image")

	var id string
	var albumId string

	// request stream object to upload the blob via
	blobStream, err := h.Blob.Put(ctx)
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

	return nil
}

// Download uploads image from server image album.
// It relies on blob-srv to retrieve the image from remote filesystem.
// It returns error if the image failed to be retrieved from blob-srv.
func (h *Image) Download(ctx context.Context, req *pb.DownloadImageReq, stream pb.Image_DownloadStream) error {
	log.Logf("Received request to download image")

	// request stream to use for blob download
	blobStream, err := h.Blob.Get(ctx, &blob.GetReq{Id: req.Id, BucketId: req.AlbumId})
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

	return nil
}

// Delete deletes the image from the given album
// It relies on blob-srv to delete the image from the physical storage.
// It returns error if the image failed to be removed.
func (h *Image) Delete(ctx context.Context, req *pb.DeleteImageReq, resp *pb.DeleteImageResp) error {
	log.Logf("Received request to delete image")

	// request stream to use for blob download
	_, err := h.Blob.Delete(ctx, &blob.DeleteReq{Id: req.Id, BucketId: req.AlbumId})
	if err != nil {
		log.Logf("failed to delete image %s/%s: %s", req.Id, req.AlbumId, err)
		return fmt.Errorf("failed to delete image %s: %s", req.Id, err)
	}

	return nil
}

// List lists images in given album.
// It returns error if the images failed to be listed.
func (h *Image) List(ctx context.Context, req *pb.ListImageReq, resp *pb.ListImageResp) error {
	log.Logf("Received request to list images")

	// list all available images in given album
	imgs, err := h.Blob.List(ctx, &blob.ListReq{BucketId: req.AlbumId})
	if err != nil {
		log.Logf("failed to list images in album %s: %s", req.AlbumId, err)
		return err
	}

	resp.Id = imgs.Id

	return nil
}
