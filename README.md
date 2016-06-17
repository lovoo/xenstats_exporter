# Xenstats Exporter

[![GoDoc](https://godoc.org/github.com/lovoo/xenstats_exporter?status.svg)](https://godoc.org/github.com/lovoo/xenstats_exporter)

 Xenstats exporter for prometheus.io, written in go.

## Docker Usage

    docker run --privileged -d --name xenstats_exporter -p 9290:9290 lovoo/xenstats_exporter:latest

## Building

    go get -u github.com/lovoo/xenstats_exporter
    go install github.com/lovoo/xenstats_exporter

## Contributing

1. Fork it!
2. Create your feature branch: `git checkout -b my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin my-new-feature`
5. Submit a pull request
