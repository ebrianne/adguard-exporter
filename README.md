# AdguardHome Prometheus Exporter

![Build/Push (master)](https://github.com/ebrianne/adguard-exporter/workflows/Build/Push%20(master)/badge.svg?branch=master)
[![GoDoc](https://godoc.org/github.com/ebrianne/adguard-exporter?status.png)](https://godoc.org/github.com/ebrianne/adguard-exporter)
[![GoReportCard](https://goreportcard.com/badge/github.com/ebrianne/adguard-exporter)](https://goreportcard.com/report/github.com/ebrianne/adguard-exporter)
![DockerPulls](https://img.shields.io/docker/pulls/ebrianne/adguard-exporter)

This is a Prometheus exporter for [AdguardHome](https://github.com/AdguardTeam/AdguardHome)'s Raspberry PI ad blocker.
It is based on the famous pihole-exporter [available here](https://github.com/eko/pihole-exporter/)

![Grafana dashboard](https://raw.githubusercontent.com/ebrianne/adguard-exporter/master/grafana/dashboard.png)

Grafana dashboard is [available here](https://grafana.com/dashboards/13330) on the Grafana dashboard website and also [here](https://raw.githubusercontent.com/ebrianne/adguard-exporter/master/grafana/dashboard.json) on the GitHub repository.

## Prerequisites

* [Go](https://golang.org/doc/)

## Installation

### Download binary

You can download the latest version of the binary built for your architecture here:

* Architecture **i386** [
[Darwin](https://github.com/ebrianne/adguard-exporter/releases/latest/download/adguard_exporter-darwin-386) /
[FreeBSD](https://github.com/ebrianne/adguard-exporter/releases/latest/download/adguard_exporter-freebsd-386) /
[Linux](https://github.com/ebrianne/adguard-exporter/releases/latest/download/adguard_exporter-linux-386) /
[Windows](https://github.com/ebrianne/adguard-exporter/releases/latest/download/adguard_exporter-windows-386.exe)
]
* Architecture **amd64** [
[Darwin](https://github.com/ebrianne/adguard-exporter/releases/latest/download/adguard_exporter-darwin-amd64) /
[FreeBSD](https://github.com/ebrianne/adguard-exporter/releases/latest/download/adguard_exporter-freebsd-amd64) /
[Linux](https://github.com/ebrianne/adguard-exporter/releases/latest/download/adguard_exporter-linux-amd64) /
[Windows](https://github.com/ebrianne/adguard-exporter/releases/latest/download/adguard_exporter-windows-amd64.exe)
]
* Architecture **arm** [
[Linux](https://github.com/ebrianne/adguard-exporter/releases/latest/download/adguard_exporter-linux-arm)
]
* Architecture **arm64** [
[Linux](https://github.com/ebrianne/adguard-exporter/releases/latest/download/adguard_exporter-linux-arm64)
]

### From sources

Optionally, you can download and build it from the sources. You have to retrieve the project sources by using one of the following way:
```bash
$ go get -u github.com/ebrianne/adguard-exporter
# or
$ git clone https://github.com/ebrianne/adguard-exporter.git
```

Install the needed vendors:

```
$ GO111MODULE=on go mod vendor
```

Then, build the binary (here, an example to run on Raspberry PI ARM architecture):
```bash
$ GOOS=linux GOARCH=arm GOARM=7 go build -o adguard_exporter .
```

## Using Docker

The exporter has been made available as a docker image. You can simply run it by the following command and pass the configuration with environment variables:

```bash
docker run \
-e 'adguard_protocol=http' \
-e 'adguard_hostname=192.168.10.252' \
-e 'adguard_username=admin' \
-e 'adguard_password=mypassword' \
-e 'adguard_port=' \ #optional if adguard is not using port 80 (http)/443 (https)
-e 'interval=10s' \
-e 'log_limit=10000' \
-e 'server_port=9617' \
-p 9617:9617 \
ebrianne/adguard-exporter:latest
```

If you prefer you can use an .env file where the environment variables are defined and using the command:

```bash
docker run --env-file=.env -p 9617:9617 \
ebrianne/adguard-exporter:latest
```

You can also use docker-compose passing the environment file or using secrets locally
### Local with environment file

```yml
version: "3.7"

services:
  adguard_exporter:
    image: ebrianne/adguard-exporter:latest
    restart: always
    ports:
      - "9617:9617"
    env_file:
      - .env
```
### Local with secret file (compose version 3 minimum)

```yml
version: "3.7"

secrets: 
  my-adguard-pass:
    file: ./my-adguard-pass.txt

services:
  adguard_exporter:
    image: ebrianne/adguard-exporter:latest
    restart: always
    secrets:
      - my-adguard-pass
    ports:
      - "9617:9617"
    environment:
      - adguard_protocol=http
      - adguard_hostname=192.168.10.252
      - adguard_username=admin
      - adguard_password=/run/secrets/my-adguard-pass
      - adguard_port= #optional
      - server_port=9617
      - interval=10s
      - log_limit=10000
```
### Swarm mode (docker swarm init)

```bash
echo "mypassword" | docker secret create my-adguard-pass -
```
Here is an example of docker-compose file.

```yml
version: "3.7"

secrets: 
  my-adguard-pass:
    external: true

services:
  adguard_exporter:
    image: ebrianne/adguard-exporter:latest
    restart: always
    secrets:
      - my-adguard-pass
    ports:
      - "9617:9617"
    environment:
      - adguard_protocol=http
      - adguard_hostname=192.168.10.252
      - adguard_username=admin
      - adguard_password=/run/secrets/my-adguard-pass
      - adguard_port= #optional
      - server_port=9617
      - interval=10s
      - log_limit=10000
```

## Usage

In order to run the exporter, type the following command (arguments are optional):

Using a password

```bash
$ ./adguard_exporter -adguard_protocol https -adguard_hostname 192.168.10.252 -adguard_username admin -adguard_password qwerty -log_limit 10000
```

```bash
2020/11/04 17:16:14 ---------------------------------------
2020/11/04 17:16:14 - AdGuard Home exporter configuration -
2020/11/04 17:16:14 ---------------------------------------
2020/11/04 17:16:14 AdguardProtocol : https
2020/11/04 17:16:14 AdguardHostname : 192.168.10.252
2020/11/04 17:16:14 AdguardUsername : admin
2020/11/04 17:16:14 AdGuard Authentication Method : AdguardPassword
2020/11/04 17:16:14 ServerPort : 9617
2020/11/04 17:16:14 Interval : 10s
2020/11/04 17:16:14 LogLimit : 10000
2020/11/04 17:16:14 ---------------------------------------
2020/11/04 17:16:14 New Prometheus metric registered: avg_processing_time
2020/11/04 17:16:14 New Prometheus metric registered: num_dns_queries
2020/11/04 17:16:14 New Prometheus metric registered: num_blocked_filtering
2020/11/04 17:16:14 New Prometheus metric registered: num_replaced_parental
2020/11/04 17:16:14 New Prometheus metric registered: num_replaced_safebrowsing
2020/11/04 17:16:14 New Prometheus metric registered: num_replaced_safesearch
2020/11/04 17:16:14 New Prometheus metric registered: top_queried_domains
2020/11/04 17:16:14 New Prometheus metric registered: top_blocked_domains
2020/11/04 17:16:14 New Prometheus metric registered: top_clients
2020/11/04 17:16:14 New Prometheus metric registered: query_types
2020/11/04 17:16:14 New Prometheus metric registered: running
2020/11/04 17:16:14 New Prometheus metric registered: protection_enabled
2020/11/04 17:16:14 Starting HTTP server
2020/11/04 17:16:30 New tick of statistics: 3824 ads blocked / 36367 total DNS queries
```

Once the exporter is running, you also have to update your `prometheus.yml` configuration to let it scrape the exporter:

```yaml
scrape_configs:
  - job_name: 'adguard'
  static_configs:
  - targets: ['localhost:9617']
```

## Available CLI options
```bash
# Interval of time the exporter will fetch data from Adguard
-interval duration (optional) (default 10s)

# Protocol to use to query Adguard
-adguard_protocol string (optional: "http", "https") (default "http")

# Hostname of the Raspberry PI where Adguard is installed
-adguard_hostname string (optional) (default "127.0.0.1")

# Username to login to Adguard Home
-adguard_username string (optional)

# Password defined on the Adguard interface
-adguard_password string (optional)

# Port to use to communicate with Adguard API
-adguard_port string (optional)

# Limit for the return log data
-log_limit string (optional) (default "1000")

# Port to be used for the exporter
-server_port string (optional) (default "9617")
```

## Available Prometheus metrics

| Metric name                       | Description                                                          |
|:---------------------------------:|----------------------------------------------------------------------|
| adguard_avg_processing_time       | This represent the average DNS query processing time                 |
| adguard_num_blocked_filtering     | This represent the number of blocked DNS queries                     |
| adguard_num_dns_queries           | This represent the number of DNS queries                             |
| adguard_num_replaced_parental     | This represent the number of blocked DNS queries (parental)          |
| adguard_num_replaced_safebrowsing | This represent the number of blocked DNS queries (safe browsing)     |
| adguard_num_replaced_safesearch   | This represent the number of blocked DNS queries (safe search)       |
| adguard_top_blocked_domains       | This represent the top blocked domains                               |
| adguard_top_clients               | This represent the top clients                                       |
| adguard_top_queried_domains       | This represent the top domains that are queried                      |
| adguard_query_types               | This represent the types of DNS queries                              |
| running                           | Is Adguard running?                                                  |
| protection_enabled                | Is the protection enabled?                                           |

## Systemd file 

### Ubuntu

One can enable the program to work at startup by writing a systemd file. You can put this file in /etc/systemd/system/adguard-home.service

```
[Unit]
Description=AdGuard-Exporter
After=syslog.target network-online.target
Requires=AdGuardHome.Service

[Service]
ExecStart=/opt/adguard_exporter/adguard_exporter-linux-arm -adguard_protocol http -adguard_hostname <hostname> -adguard_username <username> -adguard_password <password> -log_limit 5000
Restart=on-failure
RestartSec=10s

[Install]
WantedBy=multi-user.target
```

Then do this command to start the service:
```
$ sudo systemctl start adguard-home.service
```
To enable the service at startup:
```
$ sudo systemctl enable adguard-home.service
```
