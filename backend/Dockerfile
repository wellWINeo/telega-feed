FROM golang:1.24-alpine AS builder
WORKDIR /app

# fetch dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

# copy sources
COPY . .

# build
RUN CGO_ENABLED=0 GOOS=linux go build  -o ./telegafeed ./cmd/local/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/telegafeed .
RUN ls -la
EXPOSE 8080
ENTRYPOINT ["./telegafeed"]
