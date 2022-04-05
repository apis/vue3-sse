FROM golang:1.17 as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build

FROM scratch

WORKDIR /app

COPY --from=builder /app/frontend-service .

ENTRYPOINT ["/app/frontend-service"]
