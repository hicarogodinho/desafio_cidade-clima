FROM golang:1.22-alpine AS build

WORKDIR /src

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app .

FROM gcr.io/distroless/static-debian11

ENV PORT=8080
EXPOSE 8080

COPY --from=build /app /app

CMD ["/app"]

# WORKDIR /app

# FROM golang:1.21-alpine

# WORKDIR /app
# COPY . .

# RUN go mod tidy
# RUN go build -o app .

# EXPOSE 8080

# CMD ["go", "test", "-v"]