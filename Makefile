html:
	cp -R assets server/html/assets
	GOOS=js GOARCH=wasm go build -o server/html/game.wasm
	go run server/server.go