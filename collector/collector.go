package collector

import (
	"context"
	"fmt"
	"os"

	"github.com/containers/podman/v3/pkg/bindings"
	"github.com/containers/podman/v3/pkg/bindings/containers"
	"github.com/containers/podman/v3/pkg/domain/entities"
	"github.com/prometheus/client_golang/prometheus"
)

const namespace = "podman"
const subNamespace = "container"

type containerCollector struct {
	containerCPUUsageCollector           *prometheus.Desc
	containerMemoryUsageCollector        *prometheus.Desc
	containerNetworkInputUsageCollector  *prometheus.Desc
	containerNetworkOutputUsageCollector *prometheus.Desc
	containerBlockInputUsageCollector    *prometheus.Desc
	containerBlockOutputUsageCollector   *prometheus.Desc
}

func newContainerCollector() *containerCollector {
	return &containerCollector{
		containerCPUUsageCollector: prometheus.NewDesc(prometheus.BuildFQName(namespace, subNamespace, "cpu_usage"),
			"CPU Usage of the container",
			[]string{"name"}, nil,
		),
		containerMemoryUsageCollector: prometheus.NewDesc(prometheus.BuildFQName(namespace, subNamespace, "memory_usage"),
			"Memory Usage of the container",
			[]string{"name"}, nil,
		),
		containerNetworkInputUsageCollector: prometheus.NewDesc(prometheus.BuildFQName(namespace, subNamespace, "network_usage_input"),
			"Network Input usage of the container",
			[]string{"name"}, nil,
		),
		containerNetworkOutputUsageCollector: prometheus.NewDesc(prometheus.BuildFQName(namespace, subNamespace, "network_usage_output"),
			"Network output usage of the container",
			[]string{"name"}, nil,
		),
		containerBlockInputUsageCollector: prometheus.NewDesc(prometheus.BuildFQName(namespace, subNamespace, "block_usage_input"),
			"Block Input usage of the container",
			[]string{"name"}, nil,
		),
		containerBlockOutputUsageCollector: prometheus.NewDesc(prometheus.BuildFQName(namespace, subNamespace, "block_usage_output"),
			"Block output usage of the container",
			[]string{"name"}, nil,
		),
	}
}

func (collector *containerCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.containerCPUUsageCollector
	ch <- collector.containerMemoryUsageCollector
	ch <- collector.containerNetworkInputUsageCollector
	ch <- collector.containerNetworkOutputUsageCollector
	ch <- collector.containerBlockInputUsageCollector
	ch <- collector.containerBlockOutputUsageCollector
}

func (collector *containerCollector) Collect(ch chan<- prometheus.Metric) {
	ContainerStats, err := getContainerStats()
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := range ContainerStats {
		for _, stats := range i.Stats {
			fmt.Println(stats)
			ch <- prometheus.MustNewConstMetric(collector.containerCPUUsageCollector, prometheus.GaugeValue, stats.CPU, stats.Name)
			ch <- prometheus.MustNewConstMetric(collector.containerMemoryUsageCollector, prometheus.GaugeValue, float64(stats.MemUsage), stats.Name)
			ch <- prometheus.MustNewConstMetric(collector.containerNetworkInputUsageCollector, prometheus.GaugeValue, float64(stats.NetInput), stats.Name)
			ch <- prometheus.MustNewConstMetric(collector.containerNetworkOutputUsageCollector, prometheus.GaugeValue, float64(stats.NetOutput), stats.Name)
			ch <- prometheus.MustNewConstMetric(collector.containerBlockInputUsageCollector, prometheus.GaugeValue, float64(stats.BlockInput), stats.Name)
			ch <- prometheus.MustNewConstMetric(collector.containerBlockOutputUsageCollector, prometheus.GaugeValue, float64(stats.BlockOutput), stats.Name)
		}
	}
}

func Register() {
	collector := newContainerCollector()
	prometheus.MustRegister(collector)
}

func getContainerList() ([]string, error) {
	var containerNamesList []string

	sock_dir := os.Getenv("XDG_RUNTIME_DIR")
	socket := "unix:" + sock_dir + "/podman/podman.sock"
	connText, err := bindings.NewConnection(context.Background(), socket)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	containerList, err := containers.List(connText, new(containers.ListOptions))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for _, c := range containerList {
		containerNamesList = append(containerNamesList, c.Names[0])
	}

	return containerNamesList, nil
}

func getContainerStats() (chan entities.ContainerStatsReport, error) {
	sock_dir := os.Getenv("XDG_RUNTIME_DIR")
	socket := "unix:" + sock_dir + "/podman/podman.sock"
	connText, err := bindings.NewConnection(context.Background(), socket)
	if err != nil {
		return nil, err
	}

	containerList, err := getContainerList()
	if err != nil {
		return nil, err
	}

	return containers.Stats(connText, containerList, new(containers.StatsOptions).WithStream(false).WithInterval(1))
}
