build:
	go build -buildmode=plugin -o mods/test.so modules/test/main.go
	env MODULES_PATH=mods/ go run main.go core/