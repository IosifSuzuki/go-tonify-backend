#!/usr/bin/env bash

CMD=$1

function main() {
  	if [ "$CMD" == "local" ]; then
      rm -rf /docs
  	  ./scripts/docs.sh
		  docker-compose -f ./docker-compose-local.yml build --no-cache
		  docker-compose -f ./docker-compose-local.yml up -d
		elif [ "$CMD" == "stop" ]; then
		  docker stop tonify-server
		  docker stop tonify-db
  	fi
}

main