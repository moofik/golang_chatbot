#!/bin/bash

docker compose -f production-traefik.yml up -d --build --force-recreate