FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/main.go

FROM scratch

WORKDIR /app
COPY --from=builder /app/app .

ENV PORT=8080
EXPOSE 8080

ENTRYPOINT ["./app"]
