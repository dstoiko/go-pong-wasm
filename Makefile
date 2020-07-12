.PHONY: wasm
wasm:
	cp main.go main.go.bak
	echo '// +build js,wasm\n' | cat - main.go > temp && mv temp main.go
	GOOS=js GOARCH=wasm go build -o ./html/main.wasm .
	mv main.go.bak main.go

.PHONY: native
native:
	go build -o ./build/pong .