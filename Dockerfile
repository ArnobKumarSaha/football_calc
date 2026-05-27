FROM ghcr.io/appscode/golang-dev:1.24 AS builder
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o /football-calc ./cmd/server

FROM gcr.io/distroless/static-debian12
COPY --from=builder /football-calc /football-calc
EXPOSE 8080
ENTRYPOINT ["/football-calc"]
