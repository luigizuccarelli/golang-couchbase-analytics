FROM registry.access.redhat.com/ubi8/ubi-init:latest

LABEL maintainer="lzuccarelli@tfd.ie"

RUN dnf remove -y subscription-manager
# gcc for cgo
RUN dnf install -y git gcc make && rm -rf /var/lib/apt/lists/*

ENV GOLANG_VERSION 1.14.2
ENV GOLANG_DOWNLOAD_URL https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz
ENV GOLANG_DOWNLOAD_SHA256 6272d6e940ecb71ea5636ddb5fab3933e087c1356173c61f4a803895e947ebb3

RUN curl -fsSL "$GOLANG_DOWNLOAD_URL" -o golang.tar.gz \
	&& echo "$GOLANG_DOWNLOAD_SHA256  golang.tar.gz" | sha256sum -c - \
	&& tar -C /usr/local -xzf golang.tar.gz \
	&& rm golang.tar.gz

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
COPY build/microservice uid_entrypoint.sh /go/ 

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 0755 "$GOPATH"
WORKDIR $GOPATH

USER 1001

ENTRYPOINT [ "./uid_entrypoint.sh" ]

# This will change depending on each microservice entry point
CMD ["./microservice"]
