# syntax=docker/dockerfile:1

FROM golang:1.20

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

ARG GO_BUILD_ARGS
RUN echo $GO_BUILD_TAGS
RUN CGO_ENABLED=0 GOOS=linux go build $GO_BUILD_ARGS -o /demo ./cmd/demo

EXPOSE 7000

CMD ["/demo"]
