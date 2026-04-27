FROM golang:1.25 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o service-courier ./cmd/service-courier
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o worker ./cmd/worker

FROM gcr.io/distroless/base-debian12 AS service
WORKDIR /
COPY --from=builder /app/service-courier /service-courier
COPY .env /.env
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/service-courier"]

FROM gcr.io/distroless/base-debian12 AS worker

WORKDIR /
COPY --from=builder /app/worker ./worker
COPY .env /.env

USER nonroot:nonroot
ENTRYPOINT ["/worker"]