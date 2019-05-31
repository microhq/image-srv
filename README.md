# Image Service

**THIS IS A WIP: THE IMPLEMENTATION OF THIS SERVICE HAS NOT BEEN COMPLETED YET**

This is the Image service

Generated with

```
micro new image-srv --namespace=go.micro --type=srv
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

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
