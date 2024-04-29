FROM ubuntu:latest

COPY . /opt/dps-client

RUN apt update && apt install -y golang
RUN apt install -y git curl perl
RUN apt install -y make gcc tar build-essential
RUN apt install -y upx

# install opkg-build
RUN curl https://git.yoctoproject.org/opkg-utils/snapshot/opkg-utils-0.6.3.tar.gz > opkg-utils-0.6.3.tar.gz
RUN tar -xvf opkg-utils-0.6.3.tar.gz
WORKDIR /opkg-utils-0.6.3
RUN make
RUN cp opkg-build /usr/local/bin


COPY . /opt/dps-client
WORKDIR /opt/dps-client
RUN go get -d -v ./...
# RUN go build -o dps-client ./cmd/dps-client/
