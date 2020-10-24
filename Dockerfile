FROM golang:1.14 as backend_builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -v \
    -o build/run_server cmd/run_server.go cmd/handler.go

FROM node:14-slim as frontend_builder

WORKDIR /app

COPY frontend/package.json frontend/package-lock.json frontend/
RUN cd frontend && \
    npm install

COPY frontend/ frontend/
RUN cd frontend && \
    npm run build

FROM alpine:3

COPY --from=backend_builder /app/build/run_server /run_server
COPY --from=frontend_builder /app/frontend/out /ui

EXPOSE $PORT

CMD ["/run_server"]
