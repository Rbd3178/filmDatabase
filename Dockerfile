FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -v ./cmd/apiserver

# Stage 2: Create a minimal Linux distro and run the app
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/apiserver .
COPY --from=builder /app/configs ./configs

EXPOSE 8080

CMD ["./apiserver"]