GOPATH := $(shell go env GOPATH)

bin/httpserver: httpserver/main.go
	mkdir -p ./bin
	cd httpserver && go build -o ../bin/httpserver

build: main.go index.html
	mkdir -p "./public"
	cp $(GOPATH)/misc/wasm/wasm_exec.js "./public/"
	cp index.html "./public"
	GOOS=js GOARCH=wasm go build -o "./public/wasmdemo.wasm" .

serve: build bin/httpserver
	./bin/httpserver -hostport ":6581"
