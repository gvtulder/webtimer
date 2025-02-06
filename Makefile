all: fe
	go build -o dist/webtimer .

fe:
	npx rollup -c

clean:
	rm -rf dist
