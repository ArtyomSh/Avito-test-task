FROM golang:1.20-alpine AS builder

WORKDIR /usr/local/src

RUN apk --no-cache add bash git make gcc gettext musl-dev

COPY ["app/go.mod", "app/go.sum", "./"]

RUN go mod download

COPY app ./
RUN go build -o ./bin/app cmd/main/main.go

FROM alpine as runner

WORKDIR /

COPY --from=builder /usr/local/src/bin/app /
COPY app/configs/dockerConfig.yml configs/config.yml

CMD ["/app"]
