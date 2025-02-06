all: fe
	CGO_ENABLED=0 go build -ldflags="-s -w" -o dist/webtimer .

multiarch: fe
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/webtimer_linux-amd64 .
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o dist/webtimer_linux-arm64 .
	CGO_ENABLED=0 GOOS=linux GOARM=6 GOARCH=arm go build -ldflags="-s -w" -o dist/webtimer_linux-armv6 .
	CGO_ENABLED=0 GOOS=linux GOARM=7 GOARCH=arm go build -ldflags="-s -w" -o dist/webtimer_linux-armv7 .
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dist/webtimer_darwin-amd64 .
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o dist/webtimer_darwin-arm64 .
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o dist/webtimer.exe .
	CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -ldflags="-s -w" -o dist/webtimer_arm64.exe .

fe:
	npx webpack --config webpack.config.js

clean:
	rm -rf dist
