#!/usr/bin/env bash

source .env

CMD=$1

HOST="${DB_HOST}"

function main {
	if [ "$CMD" == "migrate_up" ]; then
		migrate -path db/migration -database "postgresql://${DB_USER}:${DB_PASSWORD}@${HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_MODE}" -verbose up
	elif [ "$CMD" == "migrate_down" ]; then
		migrate -path db/migration -database "postgresql://${DB_USER}:${DB_PASSWORD}@${HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_MODE}" -verbose down
	else
		migrate -path db/migration -database "postgresql://${DB_USER}:${DB_PASSWORD}@${HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_MODE}" -verbose drop
	fi
}

main