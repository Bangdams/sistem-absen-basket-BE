FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /absen-api ./cmd/web

FROM alpine:latest

WORKDIR /app

COPY --from=build /absen-api .

COPY cmd/web/.env .

EXPOSE 8080

CMD ["./absen-api"]