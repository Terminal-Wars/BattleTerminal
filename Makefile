all:
	@rm battleterm & CGO_ENABLED=1 GOARCH=386 LDFLAGS="-Wl,-O1 -Wl,--sort-common -Wl,--as-needed -Wl,--hash-style=gnu" go build --tags nowayland -o battleterm -v

run:
	CGO_ENABLED=1 GOARCH=386 LDFLAGS="-Wl,-O1 -Wl,--sort-common -Wl,--as-needed -Wl,--hash-style=gnu" go run --tags nowayland main.go -v

#