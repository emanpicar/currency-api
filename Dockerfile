FROM golang:1.13

ARG PROXY_URI=""

WORKDIR $GOPATH/src/github.com/emanpicar/currency-api

COPY . $GOPATH/src/github.com/emanpicar/currency-api

RUN export http_proxy=$PROXY_URI && \
    export https_proxy=$PROXY_URI && \
    git config --global http.proxy $PROXY_URI && \
    go get -d -v ./...; exit 0

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o currency-api .

FROM scratch

WORKDIR /root/

COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=0 /go/src/github.com/emanpicar/currency-api .

CMD ["./currency-api"]