#####################
# Carbon Dockerfile #
#####################

FROM golang

MAINTAINER Adrian "vifino" Pistol

# Make /app a volume, for mounting for example `pwd` to easily run stuff.
VOLUME ["/app"]
WORKDIR /app

# Put the source in that directory.
COPY . /go/src/github.com/carbonsrv/carbon

RUN \
	apt-get update && apt-get install -y --no-install-recommends \
		pkgconf \
		libluajit-5.1-dev \
		libphysfs-dev \
		upx \
	&& rm -rf /var/lib/apt/lists/* \
	&& mkdir -p /go/src/github.com/carbonsrv \
	&& cd /go/src/github.com/carbonsrv/carbon && go get -t -d -v ./... \
	&& go build -v -o /go/bin/carbon \
	&& strip --strip-all /go/bin/carbon \
	&& upx -9 /go/bin/carbon \
	&& apt-get remove -y pkgconf upx

# Run the carbon repl by default!
ENTRYPOINT ["/go/bin/carbon"]
CMD ["-repl"]

# Expose default ports.

EXPOSE 8080
EXPOSE 8443
