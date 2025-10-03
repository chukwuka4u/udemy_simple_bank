postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=mysecret -d postgres:12-alpine
createdb:
	docker exec -it postgres12 createdb --username=postgres --owner=postgres simple_bank
dropdb:
	docker exec -it postgres12 dropdb --username=postgres simple_bank
migrateup:
	migrate -path db/migration -database "postgresql://postgres:mysecret@172.30.16.1:5432/simple_bank?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://postgres:mysecret@172.30.16.1:5432/simple_bank?sslmode=disable" -verbose down
migrateup1:
	migrate -path db/migration -database "postgresql://postgres:mysecret@172.30.16.1:5432/simple_bank?sslmode=disable" -verbose up 1
migratedown1:
	migrate -path db/migration -database "postgresql://postgres:mysecret@172.30.16.1:5432/simple_bank?sslmode=disable" -verbose down 1
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
mock: 
	mockgen -package mockdb -destination db/mock/store.go simplebank/db/sqlc Store
server:
	go run main.go

.PHONY: postgres createdb dropdb migrateup migratedown migrateup1 migratedown1 sqlc test