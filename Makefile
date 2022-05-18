all:
	@rm battleterm & CGO_ENABLED=0 GOARCH=386 go build -o battleterm -v

run:
	CGO_ENABLED=0 GOARCH=386 go run main.go -v

#