
path=cmd/ede/main.go

.PHONY: build
build:
		CGO_ENABLED=0 go build -o cmd $(path)