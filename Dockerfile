FROM alpine:latest

ENV GOPATH /go
ENV APPPATH $GOPATH/src/github.com/lovoo/xenstat_exporter

COPY . $APPPATH

RUN apk -U add --update -t build-deps go git mercurial

RUN cd $APPPATH && go get -d && go build -o /ipmi_exporter \
    && apk del --purge build-deps git mercurial curl file gcc libgcc libc-dev make automake autoconf libtool && rm -rf $GOPATH

EXPOSE 9290

ENTRYPOINT ["/xenstat_exporter"]
