
DB_URL=postgresql://vexsx:New2021!@167.99.138.192:5432/Bank_db?sslmode=disable

migrateup :
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1 :
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown :
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1 :
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

test:
	go test -v -cover ./...

server :
	go run main.go

db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

.PHONY: migrateup migrateup  migrateup1 migrateup1 test server db_docs db_schema