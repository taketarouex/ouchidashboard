FROM golang:1.14 as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -v -o build/run_server cmd/run_server.go

FROM alpine:3

COPY --from=builder /app/build/run_server /run_server

EXPOSE $PORT

CMD ["/run_server"]
