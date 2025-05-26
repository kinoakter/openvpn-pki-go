#!/usr/bin/env bash

COMPOSE_FILE="docker-compose-ovpn.yml"


start() {
  docker compose -f "${COMPOSE_FILE}" up
}

destroy () {
  docker compose -f "${COMPOSE_FILE}" down && docker compose -f "${COMPOSE_FILE}" rm -f
}
main() {
  case "$1" in
    start)
      start
      ;;
    destroy)
      destroy
      ;;
    *)
      echo "Usage: $0 {start|destroy}"
      exit 1
      ;;
  esac
}

main "$@"

exit 0