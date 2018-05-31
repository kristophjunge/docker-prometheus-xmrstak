# Docker Prometheus XMR-Stak Exporter

Dockerized Prometheus exporter for XMR-Stak mining statistics written in Go.

Exports status and hashrate.

Example output:

```
xmrstak_up{miner="default"} 1
xmrstak_hashrate{miner="default"} 188.74
```
