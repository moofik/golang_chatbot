#!/bin/bash

docker build -t chatbot_api ..
docker tag chatbot_api:latest moofik/chatbot_api:smirk
docker push moofik/chatbot_api:smirk