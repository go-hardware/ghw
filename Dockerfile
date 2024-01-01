FROM golang:1.19-buster as builder
WORKDIR /go/src/github.com/go-hardware/ghw

ENV GOPROXY=direct

# go.mod and go.sum go into their own layers.
COPY go.mod .
COPY go.sum .

# This ensures `go mod download` happens only when go.mod and go.sum change.
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o ghw ./cmd/ghw/main.go

FROM alpine:3.7@sha256:8421d9a84432575381bfabd248f1eb56f3aa21d9d7cd2511583c68c9b7511d10
RUN apk add --no-cache ethtool

WORKDIR /bin

COPY --from=builder /go/src/github.com/go-hardware/ghw/ghwc /bin

CMD ghw
