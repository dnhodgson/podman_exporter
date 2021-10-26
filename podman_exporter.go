package main

import (
	//"context"
	//"fmt"
	//"os"

	"net/http"
	"os"

	"github.com/prometheus/common/version"
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"

	"github.com/dnhodgson/podman_exporter/collector"
)

var (
	testMetric = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "test_metric",
		Help: "Test of Metrics",
	})
)

func init() {
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(testMetric)
}

func main() {
	var (
		listenAddress = kingpin.Flag(
			"web.listen-address",
			"Address on which to expose metrics and web interface.",
		).Default("9156").String()
		metricsPath = kingpin.Flag(
			"web.telemetry-path",
			"Path under which to expose metrics.",
		).Default("/metrics").String()
	)

	promlogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.Version(version.Print("podman_exporter"))
	kingpin.CommandLine.UsageWriter(os.Stdout)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	logger := promlog.New(promlogConfig)

	level.Info(logger).Log("msg", "Starting podman_exporter", "version", version.Info())

	collector.Register()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Podman Exporter</title></head>
			<body>
			<h1>Podman Exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})
	http.Handle(*metricsPath, promhttp.Handler())
	server := ":" + *listenAddress
	http.ListenAndServe(server, nil)
}
