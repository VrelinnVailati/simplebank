.PHONY: postgres
postgres:
	docker run --name back-db -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -d postgres:15.4-alpine3.18

.PHONY: create-db
create-db:
	docker exec -it back-db createdb --username=root --owner=root simple_bank

.PHONY: drop-db
drop-db:
	docker exec -it back-db dropdb simple_bank

.PHONY: migrate-up
migrate-up:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -verbose up

.PHONY: migrate-down
migrate-down:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -verbose down

.PHONY: sqlc
sqlc:
	sqlc generate

.PHONY: initialize-db
initialize-db:
	$(MAKE) postgres
	timeout 5
	$(MAKE) create-db
	timeout 5
	$(MAKE) migrate-up

.PHONY: test
test:
	go test -v -cover ./...