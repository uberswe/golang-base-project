FROM golang:1.16.5-alpine AS builder

WORKDIR /app
COPY . .
RUN apk add --update nodejs npm
RUN apk --no-cache add ca-certificates

RUN npm install
RUN npm run build

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "main" -ldflags="-w -s" ./cmd/base/main.go

FROM scratch

WORKDIR /app
COPY --from=builder /app/active.en.toml /app/active.en.toml
COPY --from=builder /app/active.sv.toml /app/active.sv.toml
COPY --from=builder /app/main /usr/bin/
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs/

CMD ["main"]

EXPOSE 80