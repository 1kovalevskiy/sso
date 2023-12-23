# Step 1: Modules caching
FROM golang:1.21-alpine as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM golang:1.21-alpine as builder
RUN apk add --no-cache gcc musl-dev 
COPY --from=modules /go/pkg /go/pkg
WORKDIR /app
COPY ./cmd /app/cmd
COPY ./config /app/config
COPY ./pkg /app/pkg
COPY go.mod go.sum /app/
RUN CGO_ENABLED=1 go build -a -installsuffix cgo -o /bin/migrator ./cmd/migrator
COPY ./internal /app/internal
RUN CGO_ENABLED=1 go build -a -installsuffix cgo -o /bin/app ./cmd/app

# Step 3: Final
FROM golang:1.21-alpine
COPY --from=builder /bin/app /app
COPY --from=builder /bin/migrator /migrator
COPY ./config /config
COPY ./migrations /migrations
COPY scripts/start.sh /
RUN chmod +x /start.sh
WORKDIR /
CMD ["/app"]