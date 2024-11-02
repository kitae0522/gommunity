FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go run github.com/steebchen/prisma-client-go prefetch
RUN go run github.com/steebchen/prisma-client-go generate

RUN rm internal/model/query-engine-debian-openssl-3.0.x_gen.go

RUN go build -o app cmd/main.go

CMD ["./app"]