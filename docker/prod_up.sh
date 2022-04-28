#!/bin/bash

docker compose -f production-nginx.yml up -d --build --force-recreate