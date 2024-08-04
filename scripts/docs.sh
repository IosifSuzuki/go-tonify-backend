#!/usr/bin/env bash

swag fmt
swag init --parseInternal -g\
  internal/api/route/auth.go\
