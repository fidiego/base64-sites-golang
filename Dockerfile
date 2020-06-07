FROM golang:1.14.4-alpine3.12 AS builder
COPY . /code
WORKDIR /code
RUN go build -i main.go

FROM scratch
COPY --from=builder /code/main /go/bin/main
ENTRYPOINT ["/go/bin/hello"]
