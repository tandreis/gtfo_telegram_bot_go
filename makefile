APP="gtfo_tgbot"
MAIN="cmd/gtfo-bot/main.go"

build:
	go build -o bin/$(APP) $(MAIN)

run:
	go run $(MAIN)

.PHONY: clean tests

clean:
	rm -rf bin/*

tests:
	go test ./...
