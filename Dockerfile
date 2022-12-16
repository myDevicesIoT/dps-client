FROM golang:1.17.13

COPY . /opt/dps-client
WORKDIR /opt/dps-client
RUN go get -d -v ./...
RUN go build -o dps-client ./cmd/dps-client/


