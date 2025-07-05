build:
	rm -rf dist/modules/*
	mkdir -p dist/modules
	#cd modules/test/ui && ng build --output-hashing=none --base-href /v1/modules/test/ui/
	go build -buildmode=plugin -ldflags="-s -w" -o dist/modules/id.so modules/id/main.go
	go build -buildmode=plugin -ldflags="-s -w" -o dist/modules/mail.so modules/mail/main.go
	go build -buildmode=plugin -ldflags="-s -w" -o dist/modules/admin.so modules/admin/main.go
	env MODULES_PATH=dist/modules/ go run main.go core/