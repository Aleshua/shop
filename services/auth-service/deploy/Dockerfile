FROM golang:1.25.4-alpine AS base

WORKDIR /workspace

COPY go.work go.work.sum ./
COPY shared/go.mod shared/go.sum ./shared/
COPY services/auth-service/go.mod services/auth-service/go.sum ./services/auth-service/
COPY services/gateway/go.mod services/gateway/go.sum ./services/gateway/

WORKDIR /workspace/services/auth-service
RUN go mod download

WORKDIR /workspace
COPY . .

FROM base AS dev

RUN go install github.com/air-verse/air@latest

WORKDIR /workspace/services/auth-service

CMD [ "air" ]

FROM base AS builder

WORKDIR /workspace/services/auth-service

RUN CGO_ENABLED=0 GOOS=linux go build -o /main ./cmd/main.go

FROM alpine:3.19 AS prod

WORKDIR /app

COPY --from=builder /main ./main

ENTRYPOINT [ "/app/main" ]
