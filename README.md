# Xenstats Exporter

[![GoDoc](https://godoc.org/github.com/lovoo/xenstats_exporter?status.svg)](https://godoc.org/github.com/lovoo/xenstats_exporter)

 Xenstats exporter for prometheus.io, written in go.

## Docker Usage

    Building:
    docker build  --rm --force-rm --no-cache -t lovoo/xenstats_exporter .

    Running:
    docker run  -p 9290:9290 --rm  -v /path/to/config/file/:/configs lovoo/xenstats_exporter:latest  -config.file /configs/config.yml

## Building

    go get -u github.com/lovoo/xenstats_exporter
    go install github.com/lovoo/xenstats_exporter

## Config

  create a yml in form of:

```
  xenhost: "xen1.fqdn.de"
  credentials:
    username: "root"
    password: "password"
```


## Contributing

1. Fork it!
2. Create your feature branch: `git checkout -b my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin my-new-feature`
5. Submit a pull request
