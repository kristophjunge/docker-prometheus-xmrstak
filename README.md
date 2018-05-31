# Docker Prometheus XMR-Stak Exporter

Dockerized Prometheus exporter for XMR-Stak mining statistics written in Go.

Exports status and hashrate.

Example output:

```
xmrstark_up{miner="default"} 1
xmrstark_hashrate{miner="default"} 188.74
```
