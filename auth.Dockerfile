FROM golang:1.6.0-alpine

EXPOSE 8081

COPY . /go/src/quotivational

WORKDIR /go/src/quotivational

# Install auth server
RUN go build quotivational/cmd/auth

ENTRYPOINT [ "./auth" ]
CMD [ "gooduser" ]
