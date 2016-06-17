package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"

	"gopkg.in/yaml.v2"
)

var (
	listenAddress = flag.String("web.listen", ":9289", "Address on which to expose metrics and web interface.")
	metricsPath   = flag.String("web.path", "/metrics", "Path under which to expose metrics.")
	configFile    = flag.String("config.file", "config.yml", "Config file Path")
	namespace     = flag.String("xen", "xenstats", "Namespace for the IPMI metrics.")
)

func readConfig() (config Config, err error) {
	config = Config{}

	source, err := ioutil.ReadFile(*configFile)
	if err != nil {
		return config, fmt.Errorf("could not read config: %v", err)
	}

	err = yaml.Unmarshal(source, &config)
	if err != nil {
		return config, fmt.Errorf("could not unmarshal config: %v", err)
	}
	return config, err
}

func main() {
	flag.Parse()

	config, err := readConfig()
	if err != nil {
		log.Printf("%v", err)
		return
	}

	prometheus.MustRegister(NewExporter(config))

	log.Printf("Starting Server: %s", *listenAddress)
	handler := prometheus.Handler()
	if *metricsPath == "" || *metricsPath == "/" {
		http.Handle(*metricsPath, handler)
	} else {
		http.Handle(*metricsPath, handler)
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`<html>
			<head><title>IPMI Exporter</title></head>
			<body>
			<h1>IPMI Exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
		})
	}

	err = http.ListenAndServe(*listenAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
}