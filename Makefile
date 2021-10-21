.PHONY: airdrop

twap:
	mkdir -p build
	go build -o ./build/twap ./cmd/cli.go