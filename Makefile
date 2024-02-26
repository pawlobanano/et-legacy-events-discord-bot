build:
	go build -o bin/tournament-discord-bot

run: build
	./bin/tournament-discord-bot

vet:
	go vet ./...