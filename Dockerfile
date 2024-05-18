FROM golang:1.22.3-bookworm

WORKDIR /app

COPY . ./
RUN go mod download

RUN go run ./worker
RUN echo "workers started..." 