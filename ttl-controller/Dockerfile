FROM golang:1.22 as builder

WORKDIR /app
COPY . .

RUN go mod tidy
RUN go build -o controller

FROM gcr.io/distroless/base-debian12

COPY --from=builder /app/controller /controller

ENTRYPOINT ["/controller"]
