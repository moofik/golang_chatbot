#!/bin/bash

docker build -t chatbot_api ..
docker tag chatbot_api:latest moofik/chatbot_api:latest
docker push moofik/chatbot_api:latest