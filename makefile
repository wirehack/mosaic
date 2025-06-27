build:
	rm -rf dist/modules/*
	mkdir -p dist/modules/test
	cd modules/test/ui && ng build --output-hashing=none --base-href /ui/test/
	go build -buildmode=plugin -o dist/modules/test/rest.so modules/test/main.go
	env MODULES_PATH=dist/modules/ go run main.go core/