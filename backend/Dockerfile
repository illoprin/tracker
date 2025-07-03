FROM golang:1.24-alpine

WORKDIR /app

# copy necessary data only

COPY .env go.mod go.sum .

# installing dependencies

RUN go mod download

COPY . .

# build app
RUN go build -o ./dist/app ./cmd/tracker/main.go

EXPOSE 8000