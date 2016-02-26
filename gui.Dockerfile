FROM golang:1.6

RUN apt-get update && apt-get install -y \
	pkg-config \
	libgtk-3-0 \
	libgtk-3-dev \
	libgtk3.0-cil-dev \
	libcanberra-gtk3-module

RUN go get github.com/tools/godep

COPY . /go/src/github.com/endophage/quotivational

RUN cd /go/src/github.com/endophage/quotivational \
	&& godep restore \
	&& go build ./cmd/quotivational \
	&& adduser --disabled-password --no-create-home --gecos "" runner

WORKDIR /go/src/github.com/endophage/quotivational

ENTRYPOINT [ "/sbin/runuser -u runner './quotivational'" ]
