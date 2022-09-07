dev:
	go run ./cmd/main.go -root "."
	cd ./docs && npx http-server

test:
	go test ./...