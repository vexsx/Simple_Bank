migrateup :
	migrate -path db/migration -database "postgresql://vexsx:New2021!@167.99.138.192:5432/Bank_db?sslmode=disable" -verbose up

migratedown :
	migrate -path db/migration -database "postgresql://vexsx:New2021!@167.99.138.192:5432/Bank_db?sslmode=disable" -verbose down

test:
	go test -v -cover ./...

server :
	go run main.go

.PHONY: migrateup migrateup test server