FROM golang:1.22.4-alpine3.20 AS builder

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o app .

FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /app/app .

CMD ["/app/app"]