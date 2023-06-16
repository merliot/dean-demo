# syntax=docker/dockerfile:1

FROM golang:1.20

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -tags wss -o /demo ./cmd/demo

EXPOSE 7000

CMD ["/demo"]
