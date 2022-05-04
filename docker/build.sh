#!/bin/bash

docker build -t chatbot_api ..
docker tag chatbot_api:latest shrek100500/chatbot_api:smirk
docker push shrek100500/chatbot_api:smirk