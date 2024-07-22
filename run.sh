#!/usr/bin/env bash

CMD=$1

function main() {
  	if [ "$CMD" == "local" ]; then
		  docker-compose -f ./docker-compose-local.yml build --no-cache
		  docker-compose -f ./docker-compose-local.yml up -d
  	fi
}

main