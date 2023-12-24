# Step 1: Modules caching
FROM golang:1.21-alpine as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Tests
FROM golang:1.21-alpine
COPY --from=modules /go/pkg /go/pkg
RUN apk add --no-cache gcc musl-dev 
COPY . /app
WORKDIR /app

CMD ["go", "test", "-count=1", "-v", "./integration-test/..."]