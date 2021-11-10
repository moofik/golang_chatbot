FROM golang:1.17-buster
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go mod download
RUN GOOS=linux go build -o /chatbot
EXPOSE 8181
CMD [ "/chatbot"]