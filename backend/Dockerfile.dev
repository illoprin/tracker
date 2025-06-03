FROM golang:1.24-alpine

WORKDIR /app

# copy necessary data only

COPY .env go.mod go.sum .

# installing dependencies

RUN go mod download

# -- runing app in dev mode

RUN go install github.com/air-verse/air@latest

# copy air config for hot-reload

COPY .air.toml .

EXPOSE 8000

# CMD ["go", "run", "./src/main.go"]