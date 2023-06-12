# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.19 AS build-stage
#FROM harbor.mos.h-net.ru/dockerhub-proxy/library/golang:1.19 AS build-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY internal/controllers ./internal/controllers/
COPY internal/middlewares ./internal/middlewares/
COPY internal/models ./internal/models/
COPY internal/routes ./internal/routes/
COPY internal/utils ./internal/utils/
COPY cmd/otus-sonet/*.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Run the tests in the container
FROM build-stage AS run-test-stage

RUN go test -v ./...

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage
#FROM harbor.mos.h-net.ru/dockerhub-proxy/library/alpine:latest

WORKDIR /app/

COPY --from=build-stage /app/main .

EXPOSE 8000

#USER nonroot:nonroot

ENTRYPOINT ["./main"]