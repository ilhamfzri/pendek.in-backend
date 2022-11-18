create_migrate:
	migrate create -ext postgres -dir db/migrations pendek_in_migrations

migrate_up:
	migrate -database "postgres://postgres:example@localhost:5432/pendekin_db?sslmode=disable" -path db/migrations up

migrate_down:
	migrate -database "postgres://postgres:example@localhost:5432/pendekin_db?sslmode=disable" -path db/migrations down

migrate_test_up:
	migrate -database "postgres://postgres:example@localhost:5432/pendekin_db_test?sslmode=disable" -path db/migrations up

migrate_test_down:
	migrate -database "postgres://postgres:example@localhost:5432/pendekin_db_test?sslmode=disable" -path db/migrations down

up_container:
	docker-compose --env-file .env up -d