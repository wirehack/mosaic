build:
	rm -rf dist/modules/*
	mkdir -p dist/modules
	cd modules/test/ui && ng build --output-hashing=none --base-href /ui/test/
	go build -buildmode=plugin -ldflags="-s -w" -o dist/modules/rest.so modules/test/main.go
	env MODULES_PATH=dist/modules/ go run main.go core/