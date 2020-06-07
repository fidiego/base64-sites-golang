FROM golang:1.14.4-alpine3.12 AS builder

RUN apk add git

COPY . /go/src/base64_site
WORKDIR /go/src/base64_site

RUN go get
RUN go build -i main.go

FROM scratch
COPY --from=builder /go/src/base64_site/main /go/bin/main
ENTRYPOINT ["/go/bin/main"]
