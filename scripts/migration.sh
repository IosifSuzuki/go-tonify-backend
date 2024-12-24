#!/usr/bin/env bash

source .env

CMD=$1

HOST=0.0.0.0
function main {
	if [ "$CMD" == "migrate_up" ]; then
		migrate -path db/migration -database "postgresql://${POSTGRESQL_USER}:${POSTGRESQL_PASSWORD}@${HOST}:${POSTGRESQL_PORT}/${POSTGRESQL_NAME}?sslmode=${POSTGRESQL_MODE}" -verbose up
	elif [ "$CMD" == "migrate_down" ]; then
		migrate -path db/migration -database "postgresql://${POSTGRESQL_USER}:${POSTGRESQL_PASSWORD}@${HOST}:${POSTGRESQL_PORT}/${POSTGRESQL_NAME}?sslmode=${POSTGRESQL_MODE}" -verbose down
	else
		migrate -path db/migration -database "postgresql://${POSTGRESQL_USER}:${POSTGRESQL_PASSWORD}@${HOST}:${POSTGRESQL_PORT}/${POSTGRESQL_NAME}?sslmode=${POSTGRESQL_MODE}" -verbose drop
	fi
}

main
