# Image Service

This is the Image service

Generated with

```
micro new image-srv --namespace=go.micro --type=srv
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)
- [Example](#example)

## Configuration

- FQDN: go.micro.srv.image
- Type: srv
- Alias: image

## Dependencies

Micro services depend on service discovery. The default is multicast DNS, a zeroconf system.

In the event you need a resilient multi-host setup we recommend consul.

```
# install consul
brew install consul

# run consul
consul agent -dev
```

Besides the service discovery, the image service depends on the [blob service](https://github.com/microhq/blob-srv) which provides filesystem storage functionality. You can get is usual way:
```
go get github.com/microhq/blob-srv
```

## Usage

A Makefile is included for convenience

Build the binary

```
make build
```

Run the service
```
./image-srv
```

Build a docker image
```
make docker
```

## Example

Start [blob service](https://github.com/microhq/blob-srv):
```shell
MICRO_REGISTRY=consul blob-srv
2019/05/31 20:25:07 Transport [http] Listening on [::]:50477
2019/05/31 20:25:07 Broker [http] Connected to [::]:50478
2019/05/31 20:25:07 Registry [consul] Registering node: go.micro.srv.blob-994b1c3a-1b64-4092-9761-94b8cfa295d5
2019/05/31 20:25:34 Received request to create bucket: myalbum
2019/05/31 20:25:34 received put blob request
2019/05/31 20:25:34 received 1048576 bytes of blob: myalbum/mypic.jpg
2019/05/31 20:25:34 received 1048576 bytes of blob: myalbum/mypic.jpg
2019/05/31 20:25:34 received 1048576 bytes of blob: myalbum/mypic.jpg
2019/05/31 20:25:34 received 1048576 bytes of blob: myalbum/mypic.jpg
2019/05/31 20:25:34 received 10079 bytes of blob: myalbum/mypic.jpg
2019/05/31 20:25:34 finished receiving data

2019/05/31 20:27:14 received Get request to get blob mypic.jpg from bucket myalbum
2019/05/31 20:27:14 sending 1048576 bytes of myalbum/mypic.jpg
2019/05/31 20:27:14 sending 1048576 bytes of myalbum/mypic.jpg
2019/05/31 20:27:14 sending 1048576 bytes of myalbum/mypic.jpg
2019/05/31 20:27:15 sending 1048576 bytes of myalbum/mypic.jpg
2019/05/31 20:27:15 sending 10079 bytes of myalbum/mypic.jpg
2019/05/31 20:27:15 finished sending data
2019/05/31 20:27:15 closing stream socket
```

Start `image-srv` service:
```shell
$ MICRO_REGISTRY=consul image-srv
2019/05/31 20:25:24 Transport [http] Listening on [::]:50487
2019/05/31 20:25:24 Broker [http] Connected to [::]:50488
2019/05/31 20:25:24 Registry [consul] Registering node: go.micro.srv.image-a7b4f20f-2ca4-434b-997d-1daa9c5a30ab
2019/05/31 20:25:34 Received request to create new album
2019/05/31 20:25:34 Received request to upload image
2019/05/31 20:25:34 received 1048576 bytes of image: mypic.jpg
2019/05/31 20:25:34 received 1048576 bytes of image: mypic.jpg
2019/05/31 20:25:34 received 1048576 bytes of image: mypic.jpg
2019/05/31 20:25:34 received 1048576 bytes of image: mypic.jpg
2019/05/31 20:25:34 received 10079 bytes of image: mypic.jpg
2019/05/31 20:25:34 finished receiving image data
2019/05/31 20:25:34 closing stream socket

2019/05/31 20:27:14 Received request to download image
2019/05/31 20:27:14 received 1048576 bytes of image: mypic.jpg
2019/05/31 20:27:14 received 1048576 bytes of image: mypic.jpg
2019/05/31 20:27:15 received 1048576 bytes of image: mypic.jpg
2019/05/31 20:27:15 received 1048576 bytes of image: mypic.jpg
2019/05/31 20:27:15 received 10079 bytes of image: mypic.jpg
2019/05/31 20:27:15 finished receiving image data
```

Run uploader:
```shell
MICRO_REGISTRY=consul go run upload.go --album_name="myalbum" --img_path="mypic.jpg"
2019/05/31 20:25:34 streaming 1048576 bytes for image: mypic.jpg
2019/05/31 20:25:34 streaming 1048576 bytes for image: mypic.jpg
2019/05/31 20:25:34 streaming 1048576 bytes for image: mypic.jpg
2019/05/31 20:25:34 streaming 1048576 bytes for image: mypic.jpg
2019/05/31 20:25:34 streaming 10079 bytes for image: mypic.jpg
2019/05/31 20:25:34 finished uploading image: mypic.jpg
```

Run downloader:
```shell
MICRO_REGISTRY=consul go run download.go --album_name="myalbum" --img_name="mypic.jpg" --img_path="mypic.jpg"
2019/05/31 20:27:15 received 1048576 bytes for image: mypic.jpg
2019/05/31 20:27:15 received 1048576 bytes for image: mypic.jpg
2019/05/31 20:27:15 received 1048576 bytes for image: mypic.jpg
2019/05/31 20:27:15 received 1048576 bytes for image: mypic.jpg
2019/05/31 20:27:15 received 10079 bytes for image: mypic.jpg
```
