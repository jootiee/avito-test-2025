FROM golang:1.24.4-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go env -w GOPROXY=https://proxy.golang.org

RUN apk add --no-cache git && \
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

COPY . .
RUN go build -o /app ./cmd/app/main.go

FROM alpine:3.18
RUN apk add --no-cache ca-certificates curl
WORKDIR /app
COPY --from=build /app /app/server
COPY --from=build /go/bin/migrate /usr/local/bin/migrate
EXPOSE 8080
ENTRYPOINT ["/app/server"]
