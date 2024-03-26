migrateup :
	migrate -path db/migration -database "postgresql://vexsx:New2021!@167.99.138.192:5432/Bank_db?sslmode=disable" -verbose up

migrateup1 :
	migrate -path db/migration -database "postgresql://vexsx:New2021!@167.99.138.192:5432/Bank_db?sslmode=disable" -verbose up 1

migratedown :
	migrate -path db/migration -database "postgresql://vexsx:New2021!@167.99.138.192:5432/Bank_db?sslmode=disable" -verbose down

migratedown1 :
	migrate -path db/migration -database "postgresql://vexsx:New2021!@167.99.138.192:5432/Bank_db?sslmode=disable" -verbose down 1

test:
	go test -v -cover ./...

server :
	go run main.go

.PHONY: migrateup migrateup  migrateup1 migrateup1 test server