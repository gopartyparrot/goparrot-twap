.PHONY: twap

twap:
	mkdir -p build
	go build -o ./build/twap ./cmd/cli.go

twap_linux:
	mkdir -p build
	env GOOS=linux GOARCH=386 go build -o ./build/twap ./cmd/cli.go