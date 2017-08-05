# drone-git

[![Build Status](http://drone-server.ke-cs.dev.qiniu.io/api/badges/yaoshipu/drone-git/status.svg)](http://drone-server.ke-cs.dev.qiniu.io/yaoshipu/drone-git)

Drone plugin to clone `git` repositories. For the usage information and a
listing of the available options please take a look at [the docs](http://plugins.drone.io/drone-plugins/drone-git).

## Build

Build the binary with the following commands:

```
go build
go test
```

## Docker

Build the docker image with the following commands:

```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo
docker build --rm=true -t plugins/git .
```

Please note incorrectly building the image for the correct x64 linux and with
GCO disabled will result in an error when running the Docker image:

```
docker: Error response from daemon: Container command
'/bin/drone-git' not found or does not exist..
```

## Usage

Clone a commit:

```
docker run --rm \
  -e DRONE_REMOTE_URL=https://github.com/garyburd/redigo.git \
  -e DRONE_WORKSPACE=/go/src/github.com/garyburd/redigo \
  -e DRONE_BUILD_EVENT=push \
  -e DRONE_COMMIT_SHA=d8dbe4d94f15fe89232e0402c6e8a0ddf21af3ab \
  -e DRONE_COMMIT_REF=refs/heads/master \
  plugins/git
```

Clone a pull request:

```
docker run --rm \
  -e DRONE_REMOTE_URL=https://github.com/garyburd/redigo.git \
  -e DRONE_WORKSPACE=/go/src/github.com/garyburd/redigo \
  -e DRONE_BUILD_EVENT=pull_request \
  -e DRONE_COMMIT_SHA=3b4642018d177bf5fecc5907e7f341a2b5c12b8a \
  -e DRONE_COMMIT_REF=refs/pull/74/head \
  plugins/git
```

Clone a tag:

```
docker run --rm \
  -e DRONE_REMOTE_URL=https://github.com/garyburd/redigo.git \
  -e DRONE_WORKSPACE=/go/src/github.com/garyburd/redigo \
  -e DRONE_BUILD_EVENT=tag \
  -e DRONE_COMMIT_SHA=3b4642018d177bf5fecc5907e7f341a2b5c12b8a \
  -e DRONE_COMMIT_REF=refs/tags/74/head \
  plugins/git
```
