# Prometheus exporter for podman

Exports the following metrics for each running container

* CPU Usage
* Memory Usage
* Netowrk Usage
* Block Usage

Output Example

```
# HELP podman_container_block_usage_input Block Input usage of the container
# TYPE podman_container_block_usage_input gauge
podman_container_block_usage_input{name="test"} 0
podman_container_block_usage_input{name="test2"} 0
# HELP podman_container_block_usage_output Block output usage of the container
# TYPE podman_container_block_usage_output gauge
podman_container_block_usage_output{name="test"} 0
podman_container_block_usage_output{name="test2"} 0
# HELP podman_container_cpu_usage CPU Usage of the container
# TYPE podman_container_cpu_usage gauge
podman_container_cpu_usage{name="test"} 1.9770597121470862e-10
podman_container_cpu_usage{name="test2"} 1.8731004943724213e-10
# HELP podman_container_memory_usage Memory Usage of the container
# TYPE podman_container_memory_usage gauge
podman_container_memory_usage{name="test"} 1.47456e+06
podman_container_memory_usage{name="test2"} 1.396736e+06
# HELP podman_container_network_usage_input Network Input usage of the container
# TYPE podman_container_network_usage_input gauge
podman_container_network_usage_input{name="test"} 0
podman_container_network_usage_input{name="test2"} 0
# HELP podman_container_network_usage_output Network output usage of the container
# TYPE podman_container_network_usage_output gauge
podman_container_network_usage_output{name="test"} 0
podman_container_network_usage_output{name="test2"} 0
```