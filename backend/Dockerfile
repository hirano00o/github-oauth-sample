FROM golang:latest as builder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
WORKDIR /go/src/
COPY . .
RUN go build -o app main.go

FROM alpine
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/app /app
EXPOSE 8080

ENV DATABASE=mysql
ENV DSN=root:root@tcp(oauthtest:3306)/test

ENTRYPOINT ["/app"]
