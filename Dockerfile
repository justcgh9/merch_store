FROM golang:1.23 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ cmd/
COPY config/ config/
COPY internal/ internal/

RUN CGO_ENABLED=0 GOOS=linux go build -o bin/merch-store ./cmd/merch-store

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/bin/merch-store ./bin/merch-store
COPY --from=builder /app/config/ ./config

EXPOSE 8080

# установил дефолтные значения, они должны быть перезаписаны в файле docker-compose.yml
ENV CONFIG_PATH=/app/config/local.yml
ENV JWT_SECRET=defaultsecret

CMD ["./bin/merch-store"]
