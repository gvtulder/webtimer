all: fe
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/webtimer .

fe:
	npx rollup -c

clean:
	rm -rf dist
