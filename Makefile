build:
	go build -o bin/et-legacy-events-discord-bot

run: build
	./bin/et-legacy-events-discord-bot

vet:
	go vet ./...