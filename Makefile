LOCAL_MIGRATION_DIR=migrations
PG_DB_NAME=orders_storage
USER=artenso
PWD=artenso_password
LOCAL_MIGRATION_DSN='host=localhost port=5432 dbname=${PG_DB_NAME} user=${USER} password=${PWD} sslmode=disable'
M_NAME=new_column

.PHONY: local-migration-create
local-migration-create:
	goose -dir ${LOCAL_MIGRATION_DIR} -v -s postgres ${LOCAL_MIGRATION_DSN} create ${M_NAME} sql

.PHONY: local-migration-up
local-migration-up:
	goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up

.PHONY: local-migration-down
local-migration-down:
	goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down

.PHONY: local-migration-up-to
local-migration-up-to:
	goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up-to ${V}

.PHONY: local-migration-status
local-migration-status:
	goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} status