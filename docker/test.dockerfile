FROM steebchen/go-prisma:go_v1.13-prisma_v1.34.10 as builder

ENV GO111MODULE=on

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN go test -v ./...
