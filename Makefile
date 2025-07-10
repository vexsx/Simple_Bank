
DB_URL=postgresql://vexsx:New2021!@185.79.96.50:5432/Bank_db?sslmode=disable

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root Bank_db
migrateup :
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1 :
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown :
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1 :
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

test:
	go test -v -cover ./...

server :
	go run main.go

db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

proto:
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
      --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
      --grpc-gateway_out=pb  --grpc-gateway_opt=paths=source_relative \
      --openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
      proto/*.proto
	statik -src=./doc/swagger -dest=./doc

evans:
	evans --host localhost --port 9090 -r repl

redis:
	docker run --name redis -p 6379:6379 -d redis:7-alpine

.PHONY: migrateup migrateup  migrateup1 migrateup1 test server db_docs db_schema proto evans new_migration redis