all:
	@rm main && CGO_ENABLED=1 GOARCH=386 LDFLAGS="-Wl,-O3 -Wl,--sort-common -Wl,--as-needed -Wl,--hash-style=gnu" go build main.go

run:
	CGO_ENABLED=1 GOARCH=386 LDFLAGS="-Wl,-O3 -Wl,--sort-common -Wl,--as-needed -Wl,--hash-style=gnu" go run main.go

#