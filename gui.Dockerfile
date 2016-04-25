# docker run --rm \
#	-v /tmp/.X11-unix:/tmp/.X11-unix \
#	-e DISPLAY=unix$DISPLAY \
#	-v ~/.Xauthority:/tmp/Xauthority \
#	-e XAUTHORITY=/tmp/Xauthority \
#   --name quotivational \
#   endophage/apps:quotivational
FROM golang:1.6

RUN apt-get update && apt-get install -y \
	pkg-config \
	libgtk-3-0 \
	libgtk-3-dev \
	libgtk3.0-cil-dev \
	libcanberra-gtk3-module \
	xauth

COPY . /go/src/github.com/endophage/quotivational

RUN cd /go/src/github.com/endophage/quotivational \
	&& go build ./cmd/quotivational \
	&& adduser --disabled-password --no-create-home --gecos "" runner

WORKDIR /go/src/github.com/endophage/quotivational
USER runner

ENTRYPOINT [ "./quotivational" ]
