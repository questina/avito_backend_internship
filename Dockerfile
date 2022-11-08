# Telling to use Docker's golang ready image
FROM golang:1.14.6-alpine3.12 as builder
COPY go.mod go.sum /go/src/github.com/questina/
WORKDIR /go/src/github.com/questina/
RUN go mod download
COPY . /go/src/github.com/questina/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/questina .

FROM alpine
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /go/src/github.com/questina/build/questina /usr/bin/questina
EXPOSE 8000 8000
ENTRYPOINT ["/usr/bin/questina"]
