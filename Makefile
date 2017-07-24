build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo

publish:
	docker build --rm=true -t index.qiniu.com/spock/git-plugin .
	docker push index.qiniu.com/spock/git-plugin

all: build publish
