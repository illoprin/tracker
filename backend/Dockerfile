FROM golang:1.24-alpine

WORKDIR /app

# copy necessary data only

COPY .env go.mod go.sum .

# installing dependencies

RUN go mod download

EXPOSE 8000

CMD ["go", "run", "./cmd/tracker/main.go"]