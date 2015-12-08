#####################
# Carbon Dockerfile #
#####################

FROM golang

MAINTAINER Adrian "vifino" Pistol

# Make directory for sources
RUN mkdir -p /go/src/github.com/carbonsrv

# Make /app a volume, for mounting for example `pwd` to easily run stuff.
VOLUME ["/app"]
WORKDIR /app

# Get library deps
RUN apt-get update && apt-get install -y --no-install-recommends \
		pkgconf \
		libluajit-5.1-dev \
		libphysfs-dev \
	&& rm -rf /var/lib/apt/lists/*

# Copy sources over and build
COPY . /go/src/github.com/carbonsrv/carbon
RUN cd /go/src/github.com/carbonsrv/carbon && go get -t -d -v ./...
RUN cd /go/src/github.com/carbonsrv/carbon && go build -v -o /go/bin/carbon

# Run carbon -h by default!
CMD ["/go/bin/carbon", "-h"]

# Expose default ports.

EXPOSE 8080
EXPOSE 8443
