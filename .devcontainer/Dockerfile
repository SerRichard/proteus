FROM debian:trixie as devcontainer

# Set environment variables
ENV DEBIAN_FRONTEND=noninteractive

# Install dependencies and Go
RUN apt-get update && apt-get install -y \
    wget \
    curl \
    git \
    gcc \
    make \
    ca-certificates \
    golang

# Configure Go
ENV GOROOT /usr/lib/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH

RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin

WORKDIR /proteus