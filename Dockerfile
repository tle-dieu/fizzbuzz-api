FROM golang:1.14-alpine3.11 AS builder
WORKDIR /api
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o fizzbuzz_api ./cmd/

FROM scratch
COPY --from=builder /api/fizzbuzz_api .

CMD ["/fizzbuzz_api"]