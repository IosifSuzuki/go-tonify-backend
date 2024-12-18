#!/usr/bin/env bash

source .env

WEB_HOOK_PATH="/telegram/bot/update"
HOST="https://956d-2a02-8309-b001-4c00-975-eefa-7de6-aec8.ngrok-free.app"
URL="https://api.telegram.org/bot${TELEGRAM_BOT_TOKEN}/setWebhook?url=${HOST}${WEB_HOOK_PATH}"

echo "perform request: ${URL}"
curl -X POST "${URL}"
