all: fe
	go build -ldflags="-s -w" -o dist/webtimer .

fe:
	npx rollup -c

clean:
	rm -rf dist
