FROM golang:latest

COPY . /account-agent
WORKDIR /account-agent

RUN go mod download
RUN go build -o bin/ cmd/main.go
EXPOSE 8081
CMD ["/account-agent/bin/main"]