FROM golang:1.21-alpine AS builder

ARG CMD_PATH

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /server ./cmd/${CMD_PATH}/main.go

FROM alpine:3.18
RUN apk --no-cache add ca-certificates

COPY --from=builder /server /server

ENTRYPOINT ["/server"]
