mkdir deploy
go build -o ./bin/game ./cmd/game/main.go
go build -o ./bin/session ./cmd/session/main.go
go build -o ./bin/public-api ./cmd/public-api/main.go