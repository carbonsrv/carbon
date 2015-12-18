#####################
# Carbon Dockerfile #
# Arch Linux Base.  #
#####################

FROM masm/archlinux

MAINTAINER Adrian "vifino" Pistol

# Run the carbon repl by default!
ENTRYPOINT ["/usr/bin/carbon"]
CMD ["-repl"]

# Make /app a volume, for mounting for example `pwd` to easily run stuff.
VOLUME ["/app"]
WORKDIR /app

# Put the source in that directory.
COPY . /usr/local/go/src/github.com/carbonsrv/carbon

# Set the Go Path
ENV GOPATH /usr/local/go

RUN \
	pacman -Syy --noconfirm \
		base-devel git \
		luajit \
		physfs \
		upx ucl \
		go \
	&& mkdir -p "$GOPATH" \
	&& mkdir -p $GOPATH/src/github.com/carbonsrv \
	&& cd $GOPATH/src/github.com/carbonsrv/carbon && go get -t -d -v ./... \
	&& go build -v -o /usr/bin/carbon \
	&& strip --strip-all /usr/bin/carbon \
	&& upx --lzma -9 /usr/bin/carbon \
	&& pacman -R --noconfirm \
		upx \
		go \
	&& rm -rf "$GOPATH" \
	&& rm -rf /var/cache/pacman/pkg/ \
	&& rm -rf /usr/share/man/*


# Expose default ports.
EXPOSE 80
EXPOSE 443
