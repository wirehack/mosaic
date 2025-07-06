build:
	rm -rf dist/modules/*
	mkdir -p dist/modules
	go build -buildmode=plugin -ldflags="-s -w" -o dist/modules/id.so modules/id/main.go
	go build -buildmode=plugin -ldflags="-s -w" -o dist/modules/mail.so modules/mail/main.go
	go build -buildmode=plugin -ldflags="-s -w" -o dist/modules/admin.so modules/admin/main.go
	env MODULES_PATH=dist/modules/ go run main.go core/