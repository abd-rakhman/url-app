postgres:
	docker run --name url_postgres -p 5434:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=bright -d postgres:15-alpine

createdb:
	docker exec -it url_postgres createdb --username=root --owner=root url_app

dropdb:
	docker exec -it url_postgres dropdb url_app

migrateup:
	migrate -path ./db/migration -database "postgresql://root:bright@localhost:5434/url_app?sslmode=disable" -verbose up

migratedown:
	migrate -path ./db/migration -database "postgresql://root:bright@localhost:5434/url_app?sslmode=disable" -verbose down

sqlc:
	sqlc generate