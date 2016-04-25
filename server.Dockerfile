FROM golang:1.6

EXPOSE 8080

COPY . /go/src/quotivational

WORKDIR /go/src/quotivational

# Install quotivational server
RUN go build quotivational/cmd/server

ENTRYPOINT [ "./server" ]
