include .env

filename = init # default migration filename

create_migration:
	migrate create -ext=sql -dir=database/migrations -seq $(filename)

migrate_up:
	migrate -path=database/migrations -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}" -verbose up 

migrate_down:
	migrate -path=database/migrations -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}" -verbose down

migrate_force:
	migrate -path=database/migrations -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}" force $(migrate_version) 

# To make `Makefile` to run commands that do not refere to physical files (e.g., `create_migration`) in the file system
.PHONY: create_migration migrate_up migrate_down migrate_force