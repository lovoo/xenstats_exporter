FROM alpine:latest

ENV GOPATH /go
ENV APPPATH $GOPATH/src/github.com/lovoo/xenstat_exporter

RUN apk -U add --update -t build-deps go git

COPY . $APPPATH

RUN cd $APPPATH && go get -d && go build -o /xenstat_exporter \
    && apk del --purge build-deps go git mercurial curl file gcc libgcc libc-dev make automake autoconf libtool libssh2 libcurl expat pcre && rm -rf $GOPATH $APPPATH && rm -rf /var/cache/apk/*

EXPOSE 9290

ENTRYPOINT ["/xenstat_exporter"]
