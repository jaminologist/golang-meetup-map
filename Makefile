dev:
	go run ./cmd/main.go -root "."
	cd ./docs && npx http-server

archive:
	go run ./cmd/archive/main.go -root "."

test:
	go test ./...